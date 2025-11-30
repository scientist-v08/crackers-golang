package middleware

import (
	"net/http"
	"os"
	"strings"

	"slices"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/scientist-v08/crackers/initializers"
	"github.com/scientist-v08/crackers/model"
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

            // Check if a user with the given ID exists
            var user model.User
            result := initializers.DB.First(&user, uint(claims["sub"].(float64)))
            // If it doesn't exist return an error
            if result.Error != nil {
                c.JSON(http.StatusUnauthorized, gin.H{
                    "error": "Unable to find any existing user",
                })
                c.Abort()
			    return
            }

            // Check if user has any of the required roles
            hasValidRole := false
            for _, userRole := range user.Roles {
                if slices.Contains(allowedRoles, userRole) {
                    hasValidRole = true
                }
                if hasValidRole {
					c.Set("roles", user.Roles)
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