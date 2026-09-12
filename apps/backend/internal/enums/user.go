package enums

type UserRole string

const (
	ADMIN      UserRole = "admin"
	USER       UserRole = "user"
	STUDENT    UserRole = "student"
	INSTRUCTOR UserRole = "instructor"
)
