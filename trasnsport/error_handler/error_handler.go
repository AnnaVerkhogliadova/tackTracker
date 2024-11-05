package error_handler

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	internal_error "taskTracker/errors"
)

var unrecoverableInternalErrorStatus = status.New(codes.Internal, "unrecoverable internal error")

func handleTaskNotFound(err error) *status.Status {
	var errTaskNotFound internal_error.ErrTaskNotFound

	if errors.As(err, &errTaskNotFound) {
		return status.New(codes.NotFound, errTaskNotFound.Error())
	}

	return nil
}

func handleTitleTaskAlreadyExist(err error) *status.Status {
	var errTitleTaskAlreadyExist internal_error.ErrTitleTaskAlreadyExist

	if errors.As(err, &errTitleTaskAlreadyExist) {
		return status.New(codes.InvalidArgument, errTitleTaskAlreadyExist.Error())
	}

	return nil
}

func HandleError(err error) error {
	if handledErrStatus := handleTaskNotFound(err); handledErrStatus != nil {
		return handledErrStatus.Err()
	}

	if handledErrStatus := handleTitleTaskAlreadyExist(err); handledErrStatus != nil {
		return handledErrStatus.Err()
	}

	return err
}
