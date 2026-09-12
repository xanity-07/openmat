package user

import (
	"github.com/xanity-07/openmat/internal/enums"
	"github.com/xanity-07/openmat/internal/model"
)

type User struct {
	model.Base
	Email    string         `json:"email" validate:"required,email" db:"email"`
	Password string         `json:"-" validate:"required,min=8,max=32" db:"password_hash"`
	Role     enums.UserRole `json:"role" validate:"required,oneof=user admin student instructor" db:"role"`
}
