package errors

import "fmt"

type ErrTitleTaskAlreadyExist struct {
	title string
}

func NewErrTitleTaskAlreadyExist(title string) ErrTitleTaskAlreadyExist {
	return ErrTitleTaskAlreadyExist{
		title: title,
	}
}

func (err ErrTitleTaskAlreadyExist) Error() string {
	return fmt.Sprintf("Title %s already exist", err.title)
}

func (err ErrTitleTaskAlreadyExist) GetApplicationUuid() string {
	return err.title
}
