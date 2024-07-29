package pkg

import "time"

type UserLoginAction string

// const (
// 	Login  UserLoginAction = "login"
// 	Logout UserLoginAction = "logout"
// )

type LoginActivity struct {
	UserID    string    `json:"userId" bson:"userId"`
	Action    string    `json:"action" bson:"action"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}

const CacheKeyActiveUsers = "active_users"
