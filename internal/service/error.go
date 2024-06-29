package service

import "errors"

var (
	ErrTaskAlreadyCompleted = errors.New("task already completed")
)