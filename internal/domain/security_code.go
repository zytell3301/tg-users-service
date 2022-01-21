package domain

import "time"

type SecurityCode struct {
	Phone        string
	SecurityCode string
	CreatedAt    time.Time
}
