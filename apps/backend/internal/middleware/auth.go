// internal/middleware/auth.go
package middleware

import (
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xanity-07/openmat/internal/errs"
	"github.com/xanity-07/openmat/internal/lib/utils"
	"github.com/xanity-07/openmat/internal/repository"
)

func RequireAuth(jwtSecret string, sessionRepo *repository.SessionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			errs.WriteHTTPError(c, errs.NewUnauthorizedError("missing or invalid authorization header"))
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := utils.ValidateJWT(jwtSecret, tokenString)
		if err != nil {
			errs.WriteHTTPError(c, errs.NewUnauthorizedError("invalid or expired token"))
			c.Abort()
			return
		}

		sess, err := sessionRepo.Get(c, claims.SessionID)
		if err != nil {
			errs.WriteHTTPError(c, errs.NewInternalServerError())
			c.Abort()
			return
		}
		if sess == nil {
			errs.WriteHTTPError(c, errs.NewUnauthorizedError("session expired or logged out"))
			c.Abort()
			return
		}

		c.Set(UserIDKey.String(), claims.UserID)
		c.Set(UserRoleKey.String(), claims.Role)
		c.Set(SessionIDKey.String(), claims.SessionID)

		c.Next()
	}
}

func RequireRole(key CTXKey, roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(key.String())
		if !exists {
			errs.WriteHTTPError(c, errs.NewForbiddenError("no role found"))
			c.Abort()
			return
		}
		roleStr, _ := role.(string)
		if slices.Contains(roles, roleStr) {
			c.Next()
			return
		}
		errs.WriteHTTPError(c, errs.NewForbiddenError("insufficient permissions"))
		c.Abort()
	}
}

// func RequireAcademyMembership(memberRepo *repository.AcademyMemberRepository) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		academyID := c.Param("academyId")
// 		if academyID == "" {
// 			errs.WriteHTTPError(c, errs.NewBadRequestError("missing academy id", false, nil))
// 			c.Abort()
// 			return
// 		}

// 		userID := GetUserID(c) // already set by RequireAuth, which must run before this

// 		member, err := memberRepo.GetMembership(c, academyID, userID)
// 		if err != nil {
// 			errs.WriteHTTPError(c, errs.NewInternalServerError())
// 			c.Abort()
// 			return
// 		}
// 		if member == nil || member.Status != "active" {
// 			errs.WriteHTTPError(c, errs.NewForbiddenError("not a member of this academy"))
// 			c.Abort()
// 			return
// 		}

// 		c.Set(AcademyIDKey.String(), academyID)
// 		c.Set(AcademyRoleKey.String(), member.Role)

// 		c.Next()
// 	}
// }
