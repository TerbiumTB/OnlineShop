package model

import "errors"

type Status string

var OrderStatusError = errors.New("invalid order status")

const (
	CREATED  Status = "CREATED"
	SUCCESS         = "SUCCESS"
	CANCELED        = "CANCELED"
)

func ParseStatus(s string) (Status, error) {
	status := Status(s)
	switch status {
	case CREATED, SUCCESS, CANCELED:
		return status, nil
	default:
		return "", OrderStatusError
	}
}
