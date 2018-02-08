package main

import (
	"bytes"
	"fmt"
	"heshi/errors"
	"image"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
	"util"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	filetype "gopkg.in/h2non/filetype.v1"
)

func validateUploadedSingleFile(fileHeader *multipart.FileHeader, product string, fileType string, fileMaxSize int64) (string, errors.HSMessage, error) {
	if fileHeader == nil {
		return "", errors.HSMessage{}, nil
	}
	if fileHeader.Size > fileMaxSize {
		return "", errors.HSMessage{Code: 20020, Message: "File size too big"}, nil
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", errors.HSMessage{}, err
	}
	defer file.Close()

	var Buf bytes.Buffer
	io.Copy(&Buf, file)

	exts := strings.SplitN(fileHeader.Filename, ".", 2)
	if len(exts) != 2 {
		return "", errors.HSMessage{Code: 20020, Message: "Uploaded file has no extension"}, nil
	}

	var filename string
	if fileType == "video" {
		if !filetype.IsVideo([]byte(Buf.String())) {
			return "", errors.HSMessage{Code: 20020, Message: "Uploaded file is not video"}, nil
		}
		if exts[1] == "mp4" || exts[1] == "mov" || exts[1] == "ogv" || exts[1] == "webm" {
			filename = fmt.Sprintf("beyoudiamond-video-%d.%s", time.Now().UnixNano(), exts[1])
		} else {
			return "", errors.HSMessage{Code: 20020, Message: "Uploaded file extension is not supported"}, nil
		}
	} else {
		if !filetype.IsImage([]byte(Buf.String())) {
			return "", errors.HSMessage{Code: 20020, Message: "Uploaded file is not image"}, nil
		}
		filename = fmt.Sprintf("beyoudiamond-image-%d.%s", time.Now().UnixNano(), exts[1])
	}
	// Upload the file to specific dst.
	if err := os.MkdirAll(filepath.Join("."+fileType, product), 0644); err != nil {
		return "", errors.HSMessage{}, err
	}

	return filename, errors.HSMessage{}, nil
}

func saveUploadedSingleFile(c *gin.Context, product string, fileType string, filename string) error {
	if filename == "" {
		return nil
	}
	// Upload the file to specific dst.
	fileHeader, err := c.FormFile("video")
	if err != nil {
		return err
	}
	dst := filepath.Join("."+fileType, product, filename)
	err = c.SaveUploadedFile(fileHeader, dst)
	if err == nil {
		util.Println(fmt.Sprintf("'%s' uploaded!", dst))
	}
	return err
}

func validateUploadedMultipleFile(c *gin.Context, product string, fileType string, fileMaxSize int64) ([]string, errors.HSMessage, error) {
	form, err := c.MultipartForm()
	if err != nil {
		fmt.Println(errors.GetMessage(err))
	}
	files := form.File["images"]
	var fileNames []string
	for _, file := range files {
		time.Sleep(1 * time.Microsecond)
		filename, vemsg, err := validateUploadedSingleFile(file, product, fileType, fileMaxSize)
		if err != nil {
			return nil, errors.HSMessage{}, err
		}
		if vemsg != (errors.HSMessage{}) {
			return nil, vemsg, nil
		}
		fileNames = append(fileNames, filename)
	}
	return fileNames, errors.HSMessage{}, nil
}

func saveUploadedMultipleFile(c *gin.Context, product string, fileType string, fileNames []string) error {
	form, _ := c.MultipartForm()
	files := form.File["images"]
	for k, file := range files {
		dst := filepath.Join("."+fileType, product, fileNames[k])
		if err := saveImage(file, dst); err != nil {
			return err
		}
	}
	return nil
}

func saveImage(fileHeader *multipart.FileHeader, dst string) error {
	src, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	imgConfig, _, err := image.DecodeConfig(src)
	if err != nil {
		return err
	}

	if imgConfig.Height > 640 || imgConfig.Width > 320 {
		//limit image with to 320 pixel, resize image
		img, err := imaging.Decode(src)
		if err != nil {
			return err
		}
		img = imaging.Resize(img, 320, 0, imaging.Lanczos)
		return imaging.Save(img, dst)
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
