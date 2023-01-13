package middleware

import (
	"fmt"
	"net/http"
	"time"

	"package/main/internal/config"
	"package/main/internal/db"
	"package/main/internal/models"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

var conf *config.Config
var DB *gorm.DB

func init() {
	conf = config.Cfg
	DB = db.DB
}

func RequireAuth(c *gin.Context) {
	// --- get the cookie off req
	tokenString, err := c.Cookie("Authorization")

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	// --- Decode/validate it
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(conf.JWT_SECRET), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		// ---check the exp
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		// ---Find the user with token sub
		var user models.User
		DB.First(&user, claims["sub"])

		if user.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		// --- attach to req
		c.Set("user", user)

		// --- Continue

		c.Next()

	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

}
