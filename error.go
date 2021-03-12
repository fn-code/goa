package goa

import "errors"

var (
	errModeInvalid   = errors.New("error mode is invalid")
	errMaskInvalid   = errors.New("error mask is invalid")
	errEmptyByte     = errors.New("error receive empty data")
	errHeaderInvalid = errors.New("header invalid")
)

const (
	// statusOK is status ok
	statusOK = byte(0x80)
	// statusNotOK is status is not ok
	statusNotOK = byte(0x81)
)
