package events

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"service-worker-sqs-s3-postgres/core/domain/exceptions"
	cases "service-worker-sqs-s3-postgres/core/usecases/events"
	env "service-worker-sqs-s3-postgres/dataproviders/utils"
)

// EventController encapsulates all the data necessary for the implementation of the EventsService.
type EventController struct {
	eventUseCases cases.IEventCaseUses
}

// NewEventController instantiate a new event controller.
func NewEventController(es cases.IEventCaseUses) *EventController {
	return &EventController{
		eventUseCases: es,
	}
}

// GetID return a event by ID [eventsService.GetID].
func (ec *EventController) GetID(c echo.Context) error {
	ID, err := env.GetParam(c, "id")
	if err != nil {
		return exceptions.NewError(http.StatusBadRequest, err)
	}
	events, err := ec.eventUseCases.GetID(ID)
	if err != nil {
		return exceptions.HandleServiceError(err)
	}
	return c.JSON(http.StatusOK, events.ToDomainEvents())
}
