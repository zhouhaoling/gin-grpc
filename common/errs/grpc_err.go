package errs

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"test.com/common"
)

func GrpcError(err *BError) error {
	return status.Error(codes.Code(err.Code), err.Msg)
}

func ParseGrpcError(err error) (common.ResCode, string) {
	fromError, _ := status.FromError(err)
	return common.ResCode(fromError.Code()), fromError.Message()
}

func ToBError(err error) *BError {
	fromError, _ := status.FromError(err)
	return NewError(common.ResCode(fromError.Code()), fromError.Message())
}
