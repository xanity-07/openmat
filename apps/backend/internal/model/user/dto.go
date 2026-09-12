package user

import (
	"github.com/go-playground/validator/v10"
	"github.com/xanity-07/openmat/internal/enums"
)

type CreateUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (p *CreateUserPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

// ------------------------------------------------------------------------

type UpdateUserPayload struct {
	Email string `json:"email" validate:"required,email"`
	ID    string `json:"id" validate:"required,min=11,max=11"`
}

func (p *UpdateUserPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

// ------------------------------------------------------------------------

type DeleteUserPayload struct {
	ID    string `json:"id" validate:"required,min=11,max=11"`
	Email string `json:"email" validate:"required,email"`
}

func (p *DeleteUserPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

// ------------------------------------------------------------------------

type GetUserByIDPayload struct {
	ID string `json:"id" validate:"required,min=11,max=11"`
}

func (p *GetUserByIDPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

// ------------------------------------------------------------------------

type GetUsersQuery struct {
	Page   *int            `json:"page" validate:"omitempty,min=1"`
	Limit  *int            `json:"limit" validate:"omitempty,min=1"`
	Search *string         `json:"search" validate:"omitempty,min=1"`
	Sort   *string         `json:"sort" validate:"omitempty,min=1"`
	Order  *string         `json:"order" validate:"omitempty,oneof=desc asc"`
	Role   *enums.UserRole `json:"role" validate:"omitempty,oneof=user admin student instructor"`
}

func (q *GetUsersQuery) Validate() error {
	validate := validator.New()
	if err := validate.Struct(q); err != nil {
		return err
	}

	if q.Page == nil {
		defaultPage := 1
		q.Page = &defaultPage
	}

	if q.Limit == nil {
		defaultLimit := 10
		q.Limit = &defaultLimit
	}

	if q.Sort == nil {
		defaultSort := "desc"
		q.Sort = &defaultSort
	}

	return nil
}
