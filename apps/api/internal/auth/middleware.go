package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/HisakeyT/backpacker-platform/internal/model"
	"github.com/gin-gonic/gin"
)

type UserVersionLookup interface {
	FindByID(id uint) (*model.User, error)
}

func AuthMiddleware(jwtManager *JWTManager, users UserVersionLookup) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		identity, err := jwtManager.ValidateIdentity(parts[1])
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		user, err := users.FindByID(identity.UserID)
		if err != nil {
			if errors.Is(err, model.ErrUserNotFound) {
				c.AbortWithStatus(http.StatusUnauthorized)
			} else {
				c.AbortWithStatus(http.StatusInternalServerError)
			}
			return
		}
		if user == nil || user.AuthVersion != identity.AuthVersion {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("userID", identity.UserID)

		c.Next()
	}
}
