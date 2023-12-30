package controllers

import (
	"net/http"
	userApp "vix-btpns/app/user"
	"vix-btpns/helpers"
	"vix-btpns/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type userController struct {
	db *gorm.DB
}

func NewUserController(db *gorm.DB) *userController {
	return &userController{db}
}

func (h *userController) Update(c *gin.Context) {
	var input userApp.UpdateUserProfileInput

	// Get input user
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helpers.FormatValidationError(err)
		errorMessages := gin.H{"errors": errors}

		response := helpers.APIResponse("failed to update profile", http.StatusUnprocessableEntity, "error", errorMessages)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// Get logedin user
	currentUser := c.MustGet("currentUser").(models.User)

	// Update profile
	if input.Username != "" {
		currentUser.Username = input.Username
	}

	err = h.db.Save(&currentUser).Error
	if err != nil {
		response := helpers.APIResponse("failed to update profile", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := userApp.FormatUser(currentUser)
	response := helpers.APIResponse("update profile success", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

func (h *userController) Delete(c *gin.Context) {
	// Get logedin user
	currentUser := c.MustGet("currentUser").(models.User)
	
	err := h.db.Delete(&models.User{}, currentUser.ID).Error
	if err != nil {
		response := helpers.APIResponse("failed to delete user", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helpers.APIResponse("delete user success", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}
