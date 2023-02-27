package controllers

import (
	"time"
)

const (
	defaultRequeuePeriod    = 60 * time.Second
	defaultErrRequeuePeriod = 5 * time.Second
)
