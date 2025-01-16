// pkg/models/User.go - defines the data structure for the user
package models

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}