package user

import (
	"time"
)

// USER MODEL

type User struct {
	ID    	  uint32    `json:"id"`
	Firstname string    `json:"firstname" form:"alpha_zero,max=255"`
	Lastname  string    `json:"lastname" form:"alpha_zero,max=255"`
	Email     string    `json:"email" form:"required,email,max=255"`
	Password  string    `json:"password" form:"required,max=255"`
	CreatedAt time.Time `json:"created_at"`
}