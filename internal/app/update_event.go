package app

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/ptypes"

	calendarpb "github.com/farir1408/simple-calendar/pkg/api/calendar"
)

// UpdateEvent ...
func (i *Implementation) UpdateEvent(ctx context.Context, req *calendarpb.UpdateEventRequest) (*calendarpb.UpdateEventResponse, error) {
	createAt, err := ptypes.Timestamp(req.GetStartAt())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid start time: %s", req.GetStartAt())
	}
	duration, err := ptypes.Duration(req.GetDuration())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid duration: %s", req.GetDuration())
	}

	err = i.logic.UpdateEvent(
		ctx,
		req.GetTitle(),
		req.GetDescription(),
		int64(req.GetUserId()),
		createAt,
		duration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "can't update event: %s", err)
	}

	return &calendarpb.UpdateEventResponse{}, nil
}
