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

type Photo struct {
	db *gorm.DB
}

func NewPhoto(db *gorm.DB) *Photo {
	return &Photo{db}
}

func (p *Photo) Add(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(models.User)
	userID := currentUser.ID

	var numberOfPhoto int64
	p.db.Model(&models.UserPhoto{}).Where("user_id = ?", userID).Count(&numberOfPhoto)
	if numberOfPhoto > 0 {
		data := gin.H{
			"is_uploaded": false,
		}
		response := helpers.APIResponse("photo already exist", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var input photoApp.UploadUserPhotoInput
	err := c.ShouldBind(&input)
	if err != nil {
		errors := helpers.FormatValidationError(err)
		errorMessages := gin.H{"errors": errors}

		response := helpers.APIResponse("failed to upload user photo", http.StatusUnprocessableEntity, "error", errorMessages)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		errors := helpers.FormatValidationError(err)
		errorMessages := gin.H{"errors": errors}

		response := helpers.APIResponse("failed to upload user photo", http.StatusUnprocessableEntity, "error", errorMessages)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	splitedFileName := strings.Split(file.Filename, ".")
	fileFormat := splitedFileName[len(splitedFileName)-1]
	path := fmt.Sprint("images/user/", userID, "_", time.Now().Format("010206150405"), ".", fileFormat)

	err = c.SaveUploadedFile(file, "public/"+path)
	if err != nil {
		data := gin.H{
			"is_uploaded": false,
		}
		response := helpers.APIResponse("upload ke direktori gagal", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	newUserPhoto := models.UserPhoto{}
	newUserPhoto.UserID = int(currentUser.ID)
	newUserPhoto.Title = input.Title
	newUserPhoto.Caption = input.Caption
	newUserPhoto.PhotoUrl = path

	err = p.db.Create(&newUserPhoto).Error
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

func (p *Photo) Fetch(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(models.User)
	userID := currentUser.ID

	var userPhoto models.UserPhoto
	err := p.db.Where("user_id = ?", userID).Find(&userPhoto).Error
	if err != nil {
		response := helpers.APIResponse("Get photo failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if userPhoto.PhotoUrl == "" {
		response := helpers.APIResponse("User photo still empty", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := photoApp.FormatUserPhoto(userPhoto)
	response := helpers.APIResponse("", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

func (p *Photo) Modify(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(models.User)
	userID := currentUser.ID

	var userPhoto models.UserPhoto
	err := p.db.Where("user_id = ?", userID).Find(&userPhoto).Error
	if err != nil {
		response := helpers.APIResponse("Failed to update user photo", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if userPhoto.ID == 0 {
		response := helpers.APIResponse("Photo not created yet", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var input photoApp.UpdateUserPhotoInput
	err = c.ShouldBind(&input)
	if err != nil {
		response := helpers.APIResponse("Failed to update user photo", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	file, err := c.FormFile("file")
	if err == nil {
		splitedFileName := strings.Split(file.Filename, ".")
		fileFormat := splitedFileName[len(splitedFileName)-1]
		path := fmt.Sprint("images/user/", userID, "_", time.Now().Format("010206150405"), ".", fileFormat)

		err = c.SaveUploadedFile(file, "public/"+path)
		if err != nil {
			response := helpers.APIResponse("Failed to update user photo", http.StatusBadRequest, "error", nil)
			c.JSON(http.StatusBadRequest, response)
			return
		}

		userPhoto.PhotoUrl = path
	}

	if input.Title != "" {
		userPhoto.Title = input.Title
	}
	if input.Caption != "" {
		userPhoto.Caption = input.Caption
	}

	err = p.db.Save(&userPhoto).Error
	if err != nil {
		response := helpers.APIResponse("Failed to update user photo", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := photoApp.FormatUserPhoto(userPhoto)
	response := helpers.APIResponse("Update user photo success", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

func (p *Photo) Remove(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(models.User)
	userID := currentUser.ID

	var userPhoto models.UserPhoto
	err := p.db.Where("user_id = ?", userID).Find(&userPhoto).Error
	if err != nil {
		data := gin.H{
			"is_deleted": false,
		}
		response := helpers.APIResponse("Failed to delete user photo", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	userPhoto.PhotoUrl = ""

	err = p.db.Save(&userPhoto).Error
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
