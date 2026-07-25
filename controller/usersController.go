package controller

import (
	"net/http"
	"os"
	"regexp"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/scientist-v08/crackers/constants"
	"github.com/scientist-v08/crackers/initializers"
	"github.com/scientist-v08/crackers/model"
	"golang.org/x/crypto/bcrypt"
)

func Contains(slice []string, item string) bool {
    return slices.Contains(slice, item)
}

func isPasswordValid(password string) bool {
	var (
		hasMinLen   = len(password) >= 8
		hasUpper, _ = regexp.MatchString(`[A-Z]`, password)
		hasLower, _ = regexp.MatchString(`[a-z]`, password)
		hasSpecial, _ = regexp.MatchString(`[\W_]`, password)
	)
	return hasMinLen && hasUpper && hasLower && hasSpecial
}

func SignUp(c *gin.Context) {
	// Get the email or password from the request body
	var user struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	if c.Bind(&user) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read new body",
		})
		return
	}

	// Check if User already exists
	var existingUser model.User
	if err := initializers.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "User already exists",
		})
		return
	}

	// Validate password using regex
	// At least 1 lowercase, 1 uppercase, 1 special char, and 8+ chars
	if !isPasswordValid(user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password must contain at least 1 uppercase letter, 1 lowercase letter, 1 special character, and be at least 8 characters long",
		})
		return
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// Create the User
	newUser := model.User{Email: user.Email, Password: string(hash), Roles: []string{"ROLE_USER"}}
	result := initializers.DB.Create(&newUser)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create a new user",
		})
		return
	}

	// Respond
	c.JSON(200, gin.H{
		"Creating a new user": "Successful",
	})
}

func AdminSignUp(c *gin.Context) {
	// Get the email or password from the request body
	var user struct {
		Email string `json:"email"`
		Password string `json:"password"`
		IsAdmin bool `json:"isAdmin"`
	}
	if c.Bind(&user) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read new body",
		})
		return
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	if user.IsAdmin {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Only admins can use this API"})
		return
	}

	// Create the User
	newUser := model.User{Email: user.Email, Password: string(hash), Roles: []string{"ROLE_ADMIN","ROLE_USER"}}
	result := initializers.DB.Create(&newUser)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create a new user",
		})
		return
	}

	// Respond
	c.JSON(200, gin.H{
		"Creating a new user": "Successful",
	})
}

func getRoutesByRole(role bool) ([]model.Routes, error) {
	roleOfUser := ""
	if role {
		roleOfUser = constants.RoleAdmin
	} else {
		roleOfUser = constants.RoleUser
	}
	var routes []model.Routes
	err := initializers.DB.Where("role = ?", roleOfUser).Find(&routes).Error
	return routes, err
}

func Login(c *gin.Context) {
	// Get the email and password from the request body
	var user struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	if c.Bind(&user) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body",
		})
		return
	}

	// Look up requested user
	var existingUser model.User
	initializers.DB.First(&existingUser, "email = ?", user.Email)

	if existingUser.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email ID",
		})
		return
	}

	// Compare password from request body to saved password in DB
	bcryptErr := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))

	if bcryptErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid password",
		})
		return
	}

	// Generate a JWT token
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": existingUser.Roles,
		"exp": time.Now().Add(time.Hour * 12).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	hmacSampleSecret := os.Getenv("JWT_SECRET")
	tokenString, jwtTokenErr := token.SignedString([]byte(hmacSampleSecret))

	if jwtTokenErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create JWT token",
		})
		return
	}

	// Check for the role and fetch the routes
	isAnAdmin := false
	userRoles := existingUser.Roles
	if Contains(userRoles, "ROLE_ADMIN") {
        isAnAdmin = true
    }
	routes, errRoutes := getRoutesByRole(isAnAdmin)

	if errRoutes != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errRoutes,
		})
	}

	// Send it back
	c.JSON(200, gin.H{
		"access_token": tokenString,
		"routes": routes,
	})
}