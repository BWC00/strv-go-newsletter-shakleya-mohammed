package user

import (
	"time"
)

type User struct {
	ID    	  uint32    `json:"user_id"`
	Firstname string    `json:"firstname" form:"alpha_space,max=255"`
	Lastname  string    `json:"lastname" form:"alpha_space,max=255"`
	Email     string    `json:"email" form:"required,email,max=255"`
	Password  string    `json:"password" form:"required,max=255"`
	CreatedAt time.Time `json:"created_at"`
}