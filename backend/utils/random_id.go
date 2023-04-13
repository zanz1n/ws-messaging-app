package utils

import "github.com/google/uuid"

func RandomId() string {
	return uuid.New().String()
}
