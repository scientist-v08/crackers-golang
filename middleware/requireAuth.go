package middleware

import (
	"net/http"
	"os"
	"strings"

	"slices"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// RequireAuth verifies JWT tokens in the Authorization header
func RequireAnyRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Check if header is in correct format (Bearer token)
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format. Use 'Bearer {token}'",
			})
			c.Abort()
			return
		}

		// Get the token part
		tokenString := parts[1]

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, gin.Error{
					Err:  jwt.ErrSignatureInvalid,
				}
			}

			// Return the secret key for validation
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		// Handle token parsing errors
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			c.Abort()
			return
		}

		// Check if token is valid
		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token is not valid",
			})
			c.Abort()
			return
		}

		// Extract claims and set them in context
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

			raw, ok := claims["sub"].([]interface{})
			if !ok {
				// handle missing or wrong type
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Invalid JWT",
				})
				c.Abort()
				return // or whatever
			}
            // Check if a user with the given ID exists
			roles := make([]string, 0, len(raw))
			for _, v := range raw {
				if s, ok := v.(string); ok {
					roles = append(roles, s)
				}
			}

            // Check if user has any of the required roles
            hasValidRole := false
            for _, userRole := range roles {
                if slices.Contains(allowedRoles, userRole) {
                    hasValidRole = true
                }
                if hasValidRole {
					c.Set("roles", roles)
                    break
                }
            }

            // If the user does not have the role return an error
            if !hasValidRole {
                c.JSON(http.StatusUnauthorized, gin.H{
                    "error": "You do not have access to this API",
                })
                c.Abort()
			    return
            }

            // Set user ID in context for use in subsequent handlers
			c.Set("userID", claims["sub"])

            } else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token claims",
			})
			c.Abort()
			return
		}

		// Continue to the next handler if everything is valid
		c.Next()
	}
}