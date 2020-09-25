package app

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	calendarpb "github.com/farir1408/simple-calendar/pkg/api/calendar"
)

// DeleteEvent ...
func (i *Implementation) DeleteEvent(ctx context.Context, req *calendarpb.DeleteEventRequest) (*calendarpb.DeleteEventResponse, error) {
	err := i.logic.DeleteEvent(ctx, req.GetEventId())
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &calendarpb.DeleteEventResponse{}, nil
}
