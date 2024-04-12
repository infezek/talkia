package controller_users

import (
	"io"

	"github.com/gofiber/fiber/v2"
	"github.com/infezek/app-chat/pkg/usecase/usecase_user"
)

func uploadImage(c *fiber.Ctx) (usecase_user.File, error) {
	fileAvatar, nameAvatar, err := fileHeaderToBytes(c)
	if err != nil {
		return usecase_user.File{}, err
	}
	return usecase_user.File{
		File: fileAvatar,
		Name: nameAvatar,
	}, nil
}

func fileHeaderToBytes(c *fiber.Ctx) ([]byte, string, error) {
	formFile, err := c.FormFile("avatar")
	if err != nil {
		return nil, "", err
	}
	file, err := formFile.Open()
	if err != nil {
		return nil, "", err
	}
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, "", err
	}
	return bytes, formFile.Filename, nil
}
