package photo

import "vix-btpns/models"

type UserPhotoFormatter struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	Caption  string `json:"caption"`
	PhotoUrl string `json:"photo_url"`
}

func FormatUserPhoto(photo models.UserPhoto) UserPhotoFormatter {
	formatter := UserPhotoFormatter{}
	formatter.ID = photo.ID
	formatter.Title = photo.Title
	formatter.Caption = photo.Caption
	formatter.PhotoUrl = photo.PhotoUrl

	return formatter
}
