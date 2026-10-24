package middleware

import (
	"encoding/json"
	"strings"

	"slices"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/scientist-v08/crackers/initializers"
)

// RequireAuth verifies JWT tokens in the Authorization header
func RequireAnyRole(allowedRoles ...string) fiber.Handler {
	return func(c fiber.Ctx) error {
		// Get the Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header is required",
			})
		}

		// Check if header is in correct format (Bearer token)
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format. Use 'Bearer {token}'",
			})
		}

		// Get the token part
		tokenString := parts[1]

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			// Return the secret key for validation
			return []byte(initializers.Secret), nil
		})

		// Handle token parsing errors
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		// Check if token is valid
		if !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token is not valid",
			})
		}

		// Extract roles from "sub" claim
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

			raw, ok := claims["sub"].(string)
			if !ok {
				// handle missing or wrong type
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Invalid JWT",
				})
			}

            // Check if a user with the given ID exists								
			var roles []string
			if err := json.Unmarshal([]byte(raw), &roles); err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Invalid JWT - sub claim is not a valid JSON array of roles",
				})
			}

            // Check if user has any of the required roles
            hasValidRole := false
            for _, userRole := range roles {
                if slices.Contains(allowedRoles, userRole) {
                    hasValidRole = true
                }
                if hasValidRole {
                    break
                }
            }

            // If the user does not have the role return an error
            if !hasValidRole {
                return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                    "error": "You do not have access to this API",
                })
            }

        } else {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token claims",
			})
		}

		// Continue to the next handler if everything is valid
		return c.Next()
	}
}