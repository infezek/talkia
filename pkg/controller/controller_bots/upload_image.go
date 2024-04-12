package controller_bots

import (
	"fmt"
	"io"

	"github.com/gofiber/fiber/v2"
	"github.com/infezek/app-chat/pkg/usecase/usecase_bot"
)

func uploadImage(c *fiber.Ctx) ([]usecase_bot.File, error) {
	fileAvatar, nameAvatar, _ := fileHeaderToBytes(c, TypeImageAvatar)
	fileBackground, BackgroundName, _ := fileHeaderToBytes(c, TypeImageBackground)
	if len(fileAvatar) == 0 && len(fileBackground) == 0 {
		return []usecase_bot.File{}, fmt.Errorf("no file to upload")
	}
	var files []usecase_bot.File
	if fileAvatar != nil {
		files = append(files, usecase_bot.File{
			File: fileAvatar,
			Name: nameAvatar,
			Type: usecase_bot.TypeImageAvatar,
		})
	}
	if fileBackground != nil {
		files = append(files, usecase_bot.File{
			File: fileBackground,
			Name: BackgroundName,
			Type: usecase_bot.TypeImageBackground,
		})
	}
	return files, nil
}

func fileHeaderToBytes(c *fiber.Ctx, typeFile TypeImage) ([]byte, string, error) {
	formFile, err := c.FormFile(typeFile.String())
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
