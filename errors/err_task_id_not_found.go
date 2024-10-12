package errors

import "fmt"

type ErrTaskNotFound struct {
	taskId uint64
}

func NewErrTaskNotFound(taskId uint64) ErrTaskNotFound {
	return ErrTaskNotFound{
		taskId: taskId,
	}
}

func (err ErrTaskNotFound) Error() string {
	return fmt.Sprintf("Task id %d not found", err.taskId)
}

func (err ErrTaskNotFound) GetApplicationUuid() uint64 {
	return err.taskId
}
