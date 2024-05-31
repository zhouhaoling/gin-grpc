package errs

import (
	codes "google.golang.org/grpc/codes"
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
