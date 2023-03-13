package gapi

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func fieldViolation(field string, err error) *errdetails.BadRequest_FieldViolation {
	return &errdetails.BadRequest_FieldViolation{
		Field: field,
		Description: err.Error(),
	}
}

func invalidArgError(violations []*errdetails.BadRequest_FieldViolation) error {
	statusInvalid := status.New(codes.InvalidArgument, "invalid arg")
	BadRequest := &errdetails.BadRequest{FieldViolations: violations}

	statusInvalidaDetail, err := statusInvalid.WithDetails(BadRequest)
	if err != nil {
		return statusInvalid.Err()
	}

	return statusInvalidaDetail.Err()
}