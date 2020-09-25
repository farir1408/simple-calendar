package app

import (
	"context"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/ptypes"

	calendarpb "github.com/farir1408/simple-calendar/pkg/api/calendar"
)

// NewEvent ...
func (i *Implementation) NewEvent(ctx context.Context, req *calendarpb.NewEventRequest) (*calendarpb.NewEventResponse, error) {
	createAt, err := ptypes.Timestamp(req.GetStartAt())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid start time: %s", req.GetStartAt())
	}
	duration, err := ptypes.Duration(req.GetDuration())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid duration: %s", req.GetDuration())
	}

	eventID, err := i.logic.CreateEvent(
		ctx,
		req.GetTitle(),
		req.GetDescription(),
		int64(req.GetUserId()),
		createAt,
		duration,
	)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "can't create event: %s", err)
	}

	resp := &calendarpb.NewEventResponse{
		Id: eventID,
	}
	return resp, nil
}
