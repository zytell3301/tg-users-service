package domain

import "time"

type SecurityCode struct {
	Phone        string
	SecurityCode string
	Action       string
	CreatedAt    time.Time
}
