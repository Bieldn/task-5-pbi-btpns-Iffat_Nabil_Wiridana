package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	photoApp "vix-btpns/app/photo"
	"vix-btpns/helpers"
	"vix-btpns/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type photoController struct {
	db *gorm.DB
}

func NewPhotoController(db *gorm.DB) *photoController {
	return &photoController{db}
}

func (h *photoController) Create(c *gin.Context) {
	// Get user from Middleware
	currentUser := c.MustGet("currentUser").(models.User)
	userID := currentUser.ID

	// Check if user photo already exist
	var numberOfPhoto int64
	h.db.Model(&models.UserPhoto{}).Where("user_id = ?", userID).Count(&numberOfPhoto)
	if numberOfPhoto > 0 {
		data := gin.H{
			"is_uploaded": false,
		}
		response := helpers.APIResponse("photo already exist", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Get input form
	var input photoApp.UploadUserPhotoInput
	err := c.ShouldBind(&input)
	if err != nil {
		errors := helpers.FormatValidationError(err)
		errorMessages := gin.H{"errors": errors}

		response := helpers.APIResponse("failed to upload user photo", http.StatusUnprocessableEntity, "error", errorMessages)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Get input file
	file, err := c.FormFile("file")
	if err != nil {
		errors := helpers.FormatValidationError(err)
		errorMessages := gin.H{"errors": errors}

		response := helpers.APIResponse("failed to upload user photo", http.StatusUnprocessableEntity, "error", errorMessages)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// Create file name
	splitedFileName := strings.Split(file.Filename, ".")
	fileFormat := splitedFileName[len(splitedFileName)-1]
	path := fmt.Sprint("images/user/", userID, "_", time.Now().Format("010206150405"), ".", fileFormat)

	// Save image to directory
	err = c.SaveUploadedFile(file, "public/"+path)
	if err != nil {
		data := gin.H{
			"is_uploaded": false,
		}
		response := helpers.APIResponse("upload ke direktori gagal", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Save image to database
	newUserPhoto := models.UserPhoto{}
	newUserPhoto.UserID = int(currentUser.ID)
	newUserPhoto.Title = input.Title
	newUserPhoto.Caption = input.Caption
	newUserPhoto.PhotoUrl = path

	err = h.db.Create(&newUserPhoto).Error
	if err != nil {
		data := gin.H{
			"is_uploaded": false,
		}
		response := helpers.APIResponse("failed to upload user photo", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{
		"is_uploaded": true,
	}
	response := helpers.APIResponse("upload user photo success", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

func (h *photoController) Get(c *gin.Context) {
	// Get user from Middleware
	currentUser := c.MustGet("currentUser").(models.User)
	userID := currentUser.ID

	// Check if user photo already exist
	var userPhoto models.UserPhoto
	err := h.db.Where("user_id = ?", userID).Find(&userPhoto).Error
	if err != nil {
		response := helpers.APIResponse("Get photo failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Check if Photo empty
	if userPhoto.PhotoUrl == "" {
		response := helpers.APIResponse("User photo still empty", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := photoApp.FormatUserPhoto(userPhoto)
	response := helpers.APIResponse("", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

func (h *photoController) Update(c *gin.Context) {
	// Get user from Middleware
	currentUser := c.MustGet("currentUser").(models.User)
	userID := currentUser.ID

	// Check if user photo already exist
	var userPhoto models.UserPhoto
	err := h.db.Where("user_id = ?", userID).Find(&userPhoto).Error
	if err != nil {
		response := helpers.APIResponse("Failed to update user photo", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Check if user photo exist
	if userPhoto.ID == 0 {
		response := helpers.APIResponse("Photo not created yet", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Get input form
	var input photoApp.UpdateUserPhotoInput
	err = c.ShouldBind(&input)
	if err != nil {
		response := helpers.APIResponse("Failed to update user photo", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Get input file
	file, err := c.FormFile("file")
	if err == nil {
		// Create file name
		splitedFileName := strings.Split(file.Filename, ".")
		fileFormat := splitedFileName[len(splitedFileName)-1]
		path := fmt.Sprint("images/user/", userID, "_", time.Now().Format("010206150405"), ".", fileFormat)

		// Save photo to directory
		err = c.SaveUploadedFile(file, "public/"+path)
		if err != nil {
			response := helpers.APIResponse("Failed to update user photo", http.StatusBadRequest, "error", nil)
			c.JSON(http.StatusBadRequest, response)
			return
		}

		userPhoto.PhotoUrl = path
	}

	// Save data to database
	if input.Title != "" {
		userPhoto.Title = input.Title
	}
	if input.Caption != "" {
		userPhoto.Caption = input.Caption
	}

	err = h.db.Save(&userPhoto).Error
	if err != nil {
		response := helpers.APIResponse("Failed to update user photo", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := photoApp.FormatUserPhoto(userPhoto)
	response := helpers.APIResponse("Update user photo success", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

func (h *photoController) Delete(c *gin.Context) {
	// Get user from Middleware
	currentUser := c.MustGet("currentUser").(models.User)
	userID := currentUser.ID

	// Check if user photo already exist
	var userPhoto models.UserPhoto
	err := h.db.Where("user_id = ?", userID).Find(&userPhoto).Error
	if err != nil {
		data := gin.H{
			"is_deleted": false,
		}
		response := helpers.APIResponse("Failed to delete user photo", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	userPhoto.PhotoUrl = ""

	err = h.db.Save(&userPhoto).Error
	if err != nil {
		data := gin.H{
			"is_deleted": false,
		}
		response := helpers.APIResponse("Failed to delete user photo", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{
		"is_deleted": true,
	}
	response := helpers.APIResponse("Delete user photo success", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}	