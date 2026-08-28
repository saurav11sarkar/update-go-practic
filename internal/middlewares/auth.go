package middlewares

import (
	"net/http"
	"strings"

	"go-prictic/internal/auth"
	"go-prictic/internal/httpresponse"

	"github.com/labstack/echo/v5"
)

type Middleware struct {
	jwtService auth.JwtService
}

func NewAuthMiddleware(jwtService auth.JwtService) *Middleware {
	return &Middleware{jwtService: jwtService}
}

func (m *Middleware) AuthMiddleware(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
			parts := strings.Fields(authHeader)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				return c.JSON(http.StatusUnauthorized, httpresponse.Error{
					Code: http.StatusUnauthorized, Message: "Missing or invalid authorization header",
				})
			}

			claims, err := m.jwtService.ValidateToken(parts[1])
			if err != nil {
				return c.JSON(http.StatusUnauthorized, httpresponse.Error{
					Code: http.StatusUnauthorized, Message: "Invalid or expired access token",
				})
			}

			if len(allowedRoles) > 0 && !roleAllowed(claims.Role, allowedRoles) {
				return c.JSON(http.StatusForbidden, httpresponse.Error{
					Code: http.StatusForbidden, Message: "You do not have permission to access this resource",
				})
			}

			c.Set("id", claims.ID)
			c.Set("email", claims.Email)
			c.Set("role", claims.Role)

			return next(c)
		}
	}
}

func roleAllowed(userRole string, allowedRoles []string) bool {
	for _, role := range allowedRoles {
		if strings.EqualFold(userRole, role) {
			return true
		}
	}
	return false
}
