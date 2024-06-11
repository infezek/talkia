package usecase_user

type UploadImageDtoInput struct {
	UserID string
	Avatar File
}

type File struct {
	File []byte
	Name string
}

type UploadImageOutput struct {
	AvatarURL string `json:"avatar"`
}
