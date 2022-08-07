package appErrors

import "errors"

var (
	ErrDeclareQueue     = errors.New("failed to declare a queue")
	ErrQueryParamsEmpty = errors.New("query parameter query parameter is empty. Please, repeat request")
	ErrParsingParams    = errors.New("problems in parsing query parameter")
	ErrInvalidPathKey   = errors.New("invalid pathKey parameter(only non-negative integers)")
	ErrReadAllBody      = errors.New("can`t read body. Please, try again")
	ErrSendToQueue      = errors.New("can`t send body to rabbitmq queue")
	ErrFormatFile       = errors.New("the provided file format is not allowed. Please upload a JPEG or PNG image")
	ErrGetBodyFromQueue = errors.New("we can`t get any body from queue (rabbit_mq)")
)
