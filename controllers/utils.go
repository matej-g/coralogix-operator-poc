package controllers

import (
	"time"
)

const (
	defaultRequeuePeriod    = 30 * time.Second
	defaultErrRequeuePeriod = 5 * time.Second
)
