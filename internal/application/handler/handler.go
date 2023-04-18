// Package handler contains the HTTP handlers for the port-related operations.
package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/canbo-x/port-service/internal/application/service"
	errs "github.com/canbo-x/port-service/internal/error"
)

// PortHandler is the HTTP handler for port-related operations.
type PortHandler struct {
	portService *service.PortService
}

// NewPortHandler creates a new PortHandler instance with the given port service.
func NewPortHandler(portService *service.PortService) *PortHandler {
	return &PortHandler{
		portService: portService,
	}
}

// GetPort handles the HTTP GET request to retrieve a port by its ID.
// It returns an appropriate error response if the ID is invalid, the port is not found,
// or there is an internal server error. Otherwise, it returns the port data as JSON.
func (h *PortHandler) GetPort(c echo.Context) error {
	id := c.Param("id")
	port, err := h.portService.GetPort(c.Request().Context(), id)
	if err == errs.ErrInvalidPortID {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if err == errs.ErrPortNotFound {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, port)
}
