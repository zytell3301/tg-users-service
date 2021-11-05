package domain

import (
	"time"
)

type User struct {
	Id            string
	Name          string
	Lastname      string
	Bio           string
	Username      string
	Phone         string
	Online_status bool
	Created_at    time.Time
}
