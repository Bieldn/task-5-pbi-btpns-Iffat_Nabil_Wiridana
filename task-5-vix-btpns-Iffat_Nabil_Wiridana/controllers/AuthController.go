package controllers

import (
	"encoding/base64"
	"net/http"
	"strconv"
	"time"

	userApp "vix-btpns/app/user"
	"vix-btpns/helpers"
	"vix-btpns/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Auth struct {
	db *gorm.DB
}

func NewAuth(db *gorm.DB) *Auth {
	return &Auth{db}
}

func (a *Auth) SignUp(c *gin.Context) {
	var input userApp.RegisterUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helpers.FormatValidationError(err)
		errorMessages := gin.H{"errors": errors}

		response := helpers.APIResponse("failed to update profile", http.StatusUnprocessableEntity, "error", errorMessages)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	user := models.User{}
	user.Username = input.Username
	user.Email = input.Email

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		errorMessages := gin.H{"errors": "pasword hash error"}
		response := helpers.APIResponse("failed to create account", http.StatusInternalServerError, "error", errorMessages)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	user.Password = string(passwordHash)

	err = a.db.Create(&user).Error
	if err != nil {
		response := helpers.APIResponse("failed to create account", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	byteUserID := []byte(strconv.FormatUint(uint64(user.ID), 10))
	encryptedUserID, err := helpers.Encrypt(byteUserID)
	if err != nil {
		response := helpers.APIResponse("failed to create account", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	encodedUserID := base64.StdEncoding.EncodeToString(encryptedUserID)
	claims := helpers.NewClaims(map[string]interface{}{
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(1 * time.Hour).Unix(),
		"sub": encodedUserID,
	})
	token, err := helpers.EncodeJWT(helpers.GetEnv("JWT_SECRET_KEY"), claims)
	if err != nil {
		response := helpers.APIResponse("failed to create account", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := userApp.FormatTokenRespons(user, token)
	response := helpers.APIResponse("your account has been craeted", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

func (a *Auth) SignIn(c *gin.Context) {
	var input userApp.LoginInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helpers.FormatValidationError(err)
		errorMessages := gin.H{"errors": errors}

		response := helpers.APIResponse("login failed", http.StatusUnprocessableEntity, "error", errorMessages)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	var user models.User

	email := input.Email
	password := input.Password

	err = a.db.Where("email = ?", email).Find(&user).Error
	if err != nil {
		response := helpers.APIResponse("login failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if user.ID == 0 {
		response := helpers.APIResponse("login failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		response := helpers.APIResponse("login failed", http.StatusUnauthorized, "error", nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	byteUserID := []byte(strconv.FormatUint(uint64(user.ID), 10))
	encryptedUserID, err := helpers.Encrypt(byteUserID)
	if err != nil {
		response := helpers.APIResponse("login failed", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	encodedUserID := base64.StdEncoding.EncodeToString(encryptedUserID)
	claims := helpers.NewClaims(map[string]interface{}{
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(1 * time.Hour).Unix(),
		"sub": encodedUserID,
	})
	token, err := helpers.EncodeJWT(helpers.GetEnv("JWT_SECRET_KEY"), claims)
	if err != nil {
		response := helpers.APIResponse("generate token failed", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := userApp.FormatTokenRespons(user, token)
	response := helpers.APIResponse("login success", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}
