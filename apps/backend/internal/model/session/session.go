package session

import "github.com/xanity-07/openmat/internal/enums"

type Session struct {
	UserID string         `json:"user_id"`
	Role   enums.UserRole `json:"role"`
}
