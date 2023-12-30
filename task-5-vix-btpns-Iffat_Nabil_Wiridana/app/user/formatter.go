package user

import "vix-btpns/models"

type TokenResponsFormatter struct {
	ID       uint   `json:"id"`
	Username string `json:"name"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

func FormatTokenRespons(user models.User, token string) TokenResponsFormatter {
	formatter := TokenResponsFormatter{}
	formatter.ID = user.ID
	formatter.Username = user.Username
	formatter.Email = user.Email
	formatter.Token = token

	return formatter
}

type UserFormatter struct {
	ID       uint   `json:"id"`
	Username string `json:"name"`
	Email    string `json:"email"`
}

func FormatUser(user models.User) UserFormatter {
	formatter := UserFormatter{}
	formatter.ID = user.ID
	formatter.Username = user.Username
	formatter.Email = user.Email

	return formatter
}