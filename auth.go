package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) requireAuth(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
		return
	}
	// the header looks like : Bearer <token>
	parts := strings.Split(header, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization"})
		return
	}
	tokenString := parts[1]
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		//security : make sure the token was signed the way WE sign, not some
		// attacker-chosen method. Then hand back our secret to verify with/
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(cfg.jwtSecret), nil
	})
	if err != nil || !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		return
	}
	c.Next() // token is genuine and unexpired → let the request through
}

func (cfg *apiConfig) login(c *gin.Context) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Temp: hardcoded check, real apps lok up the user in the DB and
	// compare a bcrypt-hashed passowrd.
	if creds.Username != "admin" || creds.Password != "password123" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	// BUILD the token's claims - the info written on the wristband
	claims := jwt.MapClaims{
		"username": creds.Username,
		"exp":      time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //HS256= sign with our secret
	signed, err := token.SignedString([]byte(cfg.jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": signed})
}
