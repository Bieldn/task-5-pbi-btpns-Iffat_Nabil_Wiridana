package photo

type UploadUserPhotoInput struct {
	Title   string `form:"title" binding:"required"`
	Caption string `form:"caption" binding:"required"`
}

type UpdateUserPhotoInput struct {
	Title   string `form:"title"`
	Caption string `form:"caption"`
}

