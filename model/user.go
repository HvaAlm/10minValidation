package model

import "time"

type User struct {
	UserId       string
	Password     string
	LastVerified time.Time
	Token        string
}


