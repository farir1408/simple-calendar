package app

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/ptypes"

	calendarpb "github.com/farir1408/simple-calendar/pkg/api/calendar"
)

// GetEvent ...
func (i *Implementation) GetEvent(ctx context.Context, req *calendarpb.GetEventRequest) (*calendarpb.GetEventResponse, error) {
	event, err := i.logic.GetEventByID(ctx, req.GetEventId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "can't get event by id: %d, err: %s", req.GetEventId(), err)
	}
	startAt, err := ptypes.TimestampProto(event.Start)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "invalid start time: %d, err: %s", req.GetEventId(), err)
	}

	return &calendarpb.GetEventResponse{
		Item: &calendarpb.Event{
			Id:          req.GetEventId(),
			Title:       event.Title,
			Description: event.Description,
			UserId:      uint64(event.UserID),
			StartAt:     startAt,
			Duration:    ptypes.DurationProto(event.Duration),
		},
	}, nil
}
