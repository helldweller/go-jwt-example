package controllers

import (
	"net/http"
	"net/mail"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"package/main/internal/config"
	"package/main/internal/db"
	"package/main/internal/models"
)

var conf *config.Config
var DB *gorm.DB

func init() {
	conf = config.Cfg
	DB = db.DB
}

// AuthSignup    godoc
// @Summary      signup example
// @Schemes
// @Description  do signup
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body  models.User  true  "Email/Password JSON"
// @Success      200
// @Failure      400
// @Failure      404
// @Failure      500
// @Router /auth/signup [post]
func Signup(c *gin.Context) {

	// get the email/pass off req body
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// validate email
	_, err := mail.ParseAddress(body.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email is incorrect",
		})
		return
	}

	// validate password

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// Create the user
	user := models.User{Email: body.Email, Password: string(hash)}
	result := DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{})
}

// AuthLogin     godoc
// @Summary      login example
// @Schemes
// @Description  do login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body  models.User  true  "Email/Password JSON"
// @Success      200
// @Failure      400
// @Failure      404
// @Failure      500
// @Router /auth/login [post]
func Login(c *gin.Context) {
	// get the email and pass off req body
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return

	}

	// Look up requested user
	var user models.User
	DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// compare sent in pass with saved user pass hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email/password",
		})
		return
	}
	// Generate a jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(conf.JWT_SECRET))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to sign token",
		})
	}

	// send it back
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{})
}

// AuthValidate  godoc
// @Summary      validate example
// @Schemes
// @Description  do validate
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /auth/validate [get]
func Validate(c *gin.Context) {
	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}
