package filereader

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/buger/jsonparser"

	"github.com/canbo-x/port-service/internal/domain/model"
)

// JSONFileReader holds the filename and buffer size for reading the JSON file
type JSONFileReader struct {
	Filename   string
	BufferSize int
}

// ReadPorts reads ports from the JSON file and sends them to output channels
func (fr *JSONFileReader) ReadPorts(ctx context.Context, skipBroken bool) (<-chan *model.Port, <-chan error) {
	// Create the output channels
	portsCh := make(chan *model.Port, 1)
	errCh := make(chan error, 1)

	// Launch a goroutine to process the file
	go func() {
		defer close(portsCh)
		defer close(errCh)

		// Open the file
		file, err := os.Open(fr.Filename)
		if err != nil {
			errCh <- fmt.Errorf("os.Open: failed with: %w", err)
			return
		}

		// Ensure that the file is closed before returning
		defer func() {
			if err = file.Close(); err != nil {
				errCh <- fmt.Errorf("file.Close: failed with: %w", err)
			}
		}()

		// Create a new scanner to read the file line by line
		scanner := bufio.NewScanner(file)
		jsonBuffer := make([]byte, 0, fr.BufferSize)

		// Read the file line by line
		for scanner.Scan() {
			line := scanner.Bytes()

			// Skip empty lines
			if len(bytes.TrimSpace(line)) == 0 {
				continue
			}

			// Append the line to the JSON buffer
			jsonBuffer = append(jsonBuffer, line...)

			// Check if the JSON buffer has a complete JSON object
			if bytes.HasSuffix(bytes.TrimSpace(jsonBuffer), []byte("}")) {
				// Process the JSON object
				err = jsonparser.ObjectEach(jsonBuffer, func(key, value []byte, dataType jsonparser.ValueType, _ int) error {
					handleJSONValue(key, value, dataType, portsCh, errCh, skipBroken)
					return nil
				})

				// Clear the buffer for the next JSON object
				jsonBuffer = jsonBuffer[:0]
			}

			// Check if the context is done
			select {
			case <-ctx.Done():
				return
			default:
			}

			// Check scanner error
			if err := scanner.Err(); err != nil {
				errCh <- fmt.Errorf("scanner.Err: failed with: %w", err)
				return
			}
		}
	}()

	return portsCh, errCh
}

// handleJSONValue processes the JSON value and sends it to the output channel
func handleJSONValue(key, value []byte, dataType jsonparser.ValueType,
	portsCh chan<- *model.Port, errCh chan<- error, skipErrors bool,
) {
	// If the JSON value is not an object, skip it
	if dataType != jsonparser.Object {
		return
	}

	// Process the port JSON and send it to the output channel
	port, err := processPort(key, value, skipErrors)
	if err != nil {
		errCh <- err
		return
	}

	if port != nil {
		portsCh <- port
	}
}

// processPort Unmarshals the port JSON and returns a Port instance
func processPort(key, value []byte, skipErrors bool) (*model.Port, error) {
	port := new(model.Port)

	// Unmarshal the JSON value into the Port struct
	if err := json.Unmarshal(value, port); err != nil {
		// If skipping errors is enabled, return nil without an error
		if skipErrors {
			return nil, nil
		}

		// Otherwise, return an error with details
		return nil, fmt.Errorf("json.Unmarshal: failed with: %w data: (key: %s, value: %s)", err, key, value)
	}

	// Assign the key as the ID of the port
	port.ID = string(key)

	return port, nil
}
