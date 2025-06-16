package model

import "errors"

type Status string

var PaymentStatusError = errors.New("invalid order status")

const (
	SUCCESS Status = "SUCCESS"
	FAIL           = "FAIL"
)

func ParseStatus(s string) (Status, error) {
	status := Status(s)
	switch status {
	case SUCCESS, FAIL:
		return status, nil
	default:
		return "", PaymentStatusError
	}
}
