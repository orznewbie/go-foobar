package user_service

import (
	"context"

	"google.golang.org/genproto/googleapis/longrunning"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (u *UserServiceImpl) ListOperations(ctx context.Context, in *longrunning.ListOperationsRequest) (*longrunning.ListOperationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListOperations not implemented")
}

func (u *UserServiceImpl) GetOperation(ctx context.Context, in *longrunning.GetOperationRequest) (*longrunning.Operation, error) {
	return &longrunning.Operation{
		Name:     "xxxxxxxxxxx",
		Metadata: nil,
		Done:     false,
		Result:   nil,
	}, nil
}

func (u *UserServiceImpl) DeleteOperation(ctx context.Context, in *longrunning.DeleteOperationRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteOperation not implemented")
}

func (u *UserServiceImpl) CancelOperation(ctx context.Context, in *longrunning.CancelOperationRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CancelOperation not implemented")
}

func (u *UserServiceImpl) WaitOperation(ctx context.Context, in *longrunning.WaitOperationRequest) (*longrunning.Operation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WaitOperation not implemented")
}
