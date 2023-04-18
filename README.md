# Port Service
This is a simple microservice that reads port data from a file (ports.json) and either creates new records or updates existing ones in an in-memory database. The service is designed to handle files of unknown size and can be run inside a Docker container, with a Dockerfile included for easy deployment.

The project aims to demonstrate best practices in Go development and showcase a strong understanding of the language and its idiomatic usage.


## Prerequisites
To run this project:

- Go >= 1.18
- Docker
- Docker Compose (optional)
- Node.js (for `cspell`, the spell checker) (Please make sure to install Node.js to run the spell checker)

## Installation
1. Clone the repository
```bash
git clone https://github.com/canbo-x/port-service.git
cd port-service
```

2. Install dependencies
```bash
make deps
```

## Usage
### Running the Service

To build and run the service, execute the following command:
```bash
make build
./bin/port-service
```

Or, you can use Docker:
```bash
make docker-build
make docker-run
```

Or, you can use Docker Compose:
```bash
make compose-up
```

## Running Tests
To run tests, execute the following command:
```bash
make test
```

## Linting
To run the linter, execute the following command:
```bash
make lint
```

## Formatting
To format the code, execute the following command:
```bash
make format
```

## Spell Checking
To check for spelling errors in the code, execute the following command:
```bash
make spell
```

## Cleaning Up
To clean up the build artifacts, execute the following command:
```bash
make clean
```

## API
The service exposes a single HTTP endpoint:

- GET /ports/{id} - Retrieves a port record by its ID

Example response:
```json
{"id":"GBLON","name":"London","city":"London","province":"London, City of","country":"United Kingdom","alias":[],"regions":[],"coordinates":[-0.1277583,51.5073509],"timezone":"Europe/London","unlocs":["GBLON"],"code":"41352"}
```

## Signals Handling
The service can handle the following signals:

- SIGTERM: Initiates a graceful shutdown
- SIGKILL: Initiates a graceful shutdown

## Database
This service uses an in-memory database to store the port records. The in-memory database is implemented using a Go map with proper synchronization mechanisms to ensure thread-safety.

## File Reading
The service reads the ports.json file upon starting up. It processes the file line by line, allowing it to handle large files without consuming too much memory. When a port record is read, the service either creates a new record in the database or updates the existing one.

## Docker Security
This project's Dockerfile follows suggested practices for securing Docker containers:

- Using an official, minimal base image
- Running the service as a non-root user
- Avoiding compilation in the Docker container
---
## Further Improvements
## Configuration
- Consider using a configuration file (like YAML) or environment variables to configure the application, such as the path to the input file, buffer size, port, and other settings.

## Implement Logger
 - Implement a proper logging system that can be configured to log different levels of messages (info, warning, error, etc.).
 - Logging library, such as Zap or Logrus can be used.

 ## Monitoring, Metrics, and Tracing
- Add monitoring to the service to ensure that it is running properly and to detect any issues.
- Add metrics to the service to track the performance and usage of the service.
- Add tracing to the service to track the flow of requests and detect any issues.

## Error handling:
- Make sure to handle all possible errors and edge cases throughout the code. For instance file errors handling is good a starting point.

 ## Add Validations
 - Add validation for JSON object to ensure that the data is in the correct format and contains the required fields as well as complains with business requirements.
 - Add input validations.

## Tests
- Add more tests to cover more edge cases and scenarios.
- Use tags like `go:build e2e` and add more granulites to make file to perform tests in different environments and scenarios.
- Add code coverage to the tests.

## CI/CD
- Add CI/CD pipeline to automate the build, test, and deployment process.

## Github Actions (Automated actions for CI/CD)
- Add a linter to the pipeline to ensure that the code is formatted and linted properly.
- Add a spell checker to the pipeline to ensure that there are no spelling errors in the code.
- Add a code coverage tool to the pipeline to ensure that the code is fully tested.
- Add a security scanner to the pipeline to ensure that the code is secure.
- Add a formatter to the pipeline to ensure that the code is formatted properly.
- Add a Docker image builder to the pipeline to ensure that the Docker image is built properly.
---
## Personal Thoughts and Notes
## Context in `Memory DB`. 
Integrating context into the in-memory DB operations can be advantageous in specific scenarios, such as when canceling or timing out operations. However, in this case, it may not offer substantial benefits since in-memory operations are generally fast and not expected to be time-consuming. Nevertheless, adding context for demonstrating real-world scenarios and potential future improvements could be a good idea.

Two different approaches could be chosen:
1. Acquire the lock **before** checking the context cancellation
```
func (db *MemoryDB) Upsert(ctx context.Context, port *model.Port) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		db.ports[port.ID] = port
	}

	return nil
}
```
This approach ensures that the Upsert operation is atomic and not interrupted by context cancellation once the lock is acquired. However, it can lead to a situation where the Upsert operation might wait for the lock even if the context is already canceled.

Example: In a scenario requiring high data consistency, it is more critical to ensure the Upsert operation is atomic and uninterruptible once the lock is acquired, even if the context has been canceled. This guarantees that the data remains in a consistent state, despite potentially taking longer to respond to context cancellation.

2. Acquire the lock **after** checking the context cancellation
```
func (m *MemoryDB) Upsert(ctx context.Context, port *model.Port) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		m.mu.Lock()
		defer m.mu.Unlock()

		m.ports[port.ID] = port
		
		return nil
	}
}
```
In this approach, the context cancellation is checked before acquiring the lock. This allows the operation to return immediately without waiting for the lock if the context is already canceled. However, it also means that the Upsert operation can be interrupted by context cancellation after acquiring the lock, which might not be desirable in some cases.

Consider an interactive application where users initiate multiple actions that require Upsert operations. In such a case, responsiveness to user input is more important than ensuring absolute data consistency. By checking for context cancellation before acquiring the lock, the operation can quickly respond to context cancellation, resulting in a more responsive application and a better user experience.

**Conclusion**

Both implementations are suitable, depending on the requirements and preferences. If prioritizing atomicity and ensuring that the Upsert operation is not interrupted once the lock has been acquired is preferred, the first implementation would be better. If prioritizing responsiveness to context cancellation and avoiding waiting for the lock when the context is already canceled is preferred, the second implementation would be better.

## Validation of Port ID
To validate data in a Go application, two common approaches are available: using a separate validation function or a validation method on the data struct.

A separate validation function is useful for keeping the validation logic decoupled from the data struct, which promotes reusability. The validation function should be named clearly, defined in a separate file within the same package as the data struct, and kept simple.

Alternatively, a validation method can be used when you want the validation logic to be associated with the data struct, providing better encapsulation. To use this approach, the validation method should be named clearly and defined in the same file as the data struct. Keeping the method simple is also recommended.

Ultimately, the approach used depends on the specific requirements of the application and the desired level of coupling between the validation logic and the data struct.

1. Separate validation function:
```
func ValidatePortID(id string) error {
	...
}
```

This approach defines a standalone validation function that accepts a port ID as a string and returns an error if the ID is invalid. This function can be placed in a separate validator.go file within the model package.

2. Validation method on the struct:
```
func (p *Port) ValidateID() error {
	...
}
```

This approach defines a method on the Port struct that validates the port ID directly. The method can be placed within the port.go file, where the Port struct is defined.

**Conclusion**

In this demonstration, a separate validation function is chosen for the following reasons:

- Decoupling: By keeping the validation logic separate from the data struct, it promotes a cleaner separation of concerns. This makes it easier to maintain, modify, and test the validation logic without affecting the data struct itself.

- Reusability: A separate validation function can be easily reused in different parts of the application or even in other projects, without the need to include the entire data struct.

- Flexibility: Using a separate validation function allows for the possibility of validating different aspects of the data with different functions, making it easier to compose complex validation logic from simpler functions.

By choosing a separate validation function, the application's design remains modular and maintainable, while also promoting the reusability and flexibility of the validation logic.

## processPort implementation heap vs stack allocation
1. Implementation using `port := new(model.Port)`
```
func processPort(key, value []byte, skipErrors bool) (*model.Port, error) {
	port := new(model.Port)

	if err := json.Unmarshal(value, port); err != nil {
		if skipErrors {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to unmarshal port JSON: %v", err)
	}

	port.ID = string(key)

	return port, nil
}
```

2. Implementation using `var port model.Port`
```
func processPort(key, value []byte, skipErrors bool) (*model.Port, error) {
	var port model.Port

	if err := json.Unmarshal(value, &port); err != nil {
		if skipErrors {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to unmarshal port JSON: %v", err)
	}

	port.ID = string(key)

	return &port, nil
}
```

Both implementations of processPort are correct and accomplish the same task of unmarshalling a JSON into a model.

However, there are some differences in how they achieve this goal.

The first implementation uses heap allocation with port := new(model.Port) to allocate memory for a new instance of model.Port, while the second implementation uses stack allocation with var port model.Port to declare a new model.Port struct on the stack.

Heap allocation requires extra overhead for managing memory allocation and deallocation, while stack allocation is generally faster because it allocates memory directly on the function's stack frame. However, because we're returning a pointer to the model.Port struct from the function, heap allocation with new is a better choice in this case.

Using heap allocation with new ensures that the returned pointer to the model.Port struct is always valid and can be used safely outside of the function. This makes the code more robust and less error-prone, which is a better choice.

Overall, the differences in efficiency between these two implementations are likely to be small for small to medium-sized model.Port structs. The important factor is to choose the implementation that is more readable, maintainable, and follows Go's best practices, which is the first implementation using port := new(model.Port).