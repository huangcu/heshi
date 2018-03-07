package main

import (
	"bytes"
	"fmt"
	"heshi/errors"
	"image"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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

	ext := filepath.Ext(fileHeader.Filename)
	if ext == "" {
		return "", errors.HSMessage{Code: 20020, Message: "Uploaded file has no extension"}, nil
	}

	var filename string
	if fileType == "video" {
		if !filetype.IsVideo([]byte(Buf.String())) {
			return "", errors.HSMessage{Code: 20020, Message: "Uploaded file is not video"}, nil
		}
		if ext == "mp4" || ext == "mov" || ext == "ogv" || ext == "webm" {
			filename = fmt.Sprintf("beyoudiamond-video-%d.%s", time.Now().UnixNano(), ext)
		} else {
			return "", errors.HSMessage{Code: 20020, Message: "Uploaded file extension is not supported"}, nil
		}
	} else {
		if !filetype.IsImage([]byte(Buf.String())) {
			return "", errors.HSMessage{Code: 20020, Message: "Uploaded file is not image"}, nil
		}
		filename = fmt.Sprintf("beyoudiamond-image-%d.%s", time.Now().UnixNano(), ext)
	}
	// Upload the file to specific dst.
	if err := os.MkdirAll(filepath.Join("."+fileType, product), 0755); err != nil {
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

func handleImage(filepath, dst string) error {
	src, err := os.Open(filepath)
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

func bulkUpload(c *gin.Context) {
	fileHeader, _ := c.FormFile("upload")
	if fileHeader == nil {
		c.JSON(http.StatusOK, "NO FILE UPLOADED")
		return
	}

	if !strings.HasSuffix(fileHeader.Filename, ".zip") {
		c.JSON(http.StatusOK, "File name must end with .zip")
		return
	}
	//limit to 100MB
	if fileHeader.Size > 100*1024*1024 {
		c.JSON(http.StatusOK, "File size mustn't exceed 100MB")
		return
	}
	filename := filepath.Join(os.TempDir(), fileHeader.Filename)
	if err := c.SaveUploadedFile(fileHeader, filename); err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	filemsgMap, err := handleUploadedZip(filename, c.PostForm("product"), c.PostForm("fileType"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.GetMessage(err))
		return
	}
	c.JSON(http.StatusOK, filemsgMap)
	return
}

func handleUploadedZip(file, product, fileType string) (map[string]string, error) {
	tempDir := os.TempDir() + strconv.Itoa(time.Now().Nanosecond())
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return nil, err
	}
	if err := util.Unzip(file, tempDir); err != nil {
		return nil, err
	}
	if fileType == "video" {
		return handleUploadedZipVideo(tempDir), nil
	}
	return handleUploadedZipImages(tempDir, product), nil
}

//only jewelry has video
func handleUploadedZipVideo(tempDir string) map[string]string {
	fileHandleMsgMap := make(map[string]string)
	filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		var msg string
		j := jewelry{}
		if info.IsDir() {
			msg = fmt.Sprintf("SKIP folder %s", info.Name())
			util.Println(msg)
			return nil
		}
		var filename string
		bs, err := ioutil.ReadFile(path)
		if err != nil {
			msg = errors.GetMessage(err)
			return nil
		}
		if !filetype.IsVideo(bs) {
			msg = fmt.Sprintf("%s  is not video", info.Name())
			util.Println(msg)
			return nil
		}
		ext := filepath.Ext(info.Name())
		if ext == "" {
			msg = fmt.Sprintf("SKIP file %s as the file has no extension", info.Name())
			util.Println(msg)
			return nil
		}

		if ext == "mp4" || ext == "mov" || ext == "ogv" || ext == "webm" {
			filename = fmt.Sprintf("beyoudiamond-video-%s", info.Name())
		} else {
			msg = fmt.Sprintf("Uploaded file %s extension is not supported", info.Name())
			util.Println(msg)
			return nil
		}

		j.StockID = strings.Split(info.Name(), "-")[0]
		if err := j.isJewelryExistByStockID(); err != nil {
			msg = errors.GetMessage(err)
			return nil
		}
		if err := util.RunWithStdOutput("mv", path, filepath.Join(".video", "jewelry", filename)); err != nil {
			msg = errors.GetMessage(err)
			return nil
		}
		j.VideoLink = filename
		q := j.composeUpdateQuery()
		if _, err := dbExec(q); err != nil {
			msg = errors.GetMessage(err)
			return nil
		}

		if msg == "" {
			fileHandleMsgMap[info.Name()] = msg
		} else {
			fileHandleMsgMap[info.Name()] = "uploaded"
		}
		return nil
	})
	return fileHandleMsgMap
}

func handleUploadedZipImages(tempDir, product string) map[string]string {
	switch product {
	case "diamond":
		return handleUploadedZipImagesDiamond(tempDir)
	case "jewelry":
		return handleUploadedZipImagesJewelry(tempDir)
	case "gem":
		return handleUploadedZipImagesGem(tempDir)
	}
	return nil
}

//images foldname-stockid
func handleUploadedZipImagesDiamond(tempDir string) map[string]string {
	fileHandleMsgMap := make(map[string]string)
	filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		var msg string
		d := diamond{}
		if !info.IsDir() {
			msg = fmt.Sprintf("SKIP file %s", info.Name())
			util.Println(msg)
			return nil
		}
		var filenames []string
		files, _ := ioutil.ReadDir(info.Name())
		for _, file := range files {
			if file.IsDir() {
				util.Printf("SKIP folder %s/%s", info.Name(), file)
			} else {
				bs, err := ioutil.ReadFile(path)
				if err != nil {
					msg = errors.GetMessage(err)
					return nil
				}
				if !filetype.IsImage(bs) {
					msg = fmt.Sprintf("%s  is not image", info.Name())
					util.Println(msg)
					return nil
				}

				ext := filepath.Ext(info.Name())
				if ext == "" {
					msg = fmt.Sprintf("SKIP file %s as the file has no extension", info.Name())
					util.Println(msg)
				} else {
					filename := fmt.Sprintf("beyoudiamond-image-%s", info.Name())
					if err := handleImage(path, filepath.Join(".image", "diamond", filename)); err != nil {
						msg = errors.GetMessage(err)
					} else {
						filenames = append(filenames, filename)
					}
				}
			}
		}
		d.Images = filenames
		d.DiamondID = strings.Split(info.Name(), "-")[0]
		if err := d.isDiamondExistByDiamondID(); err != nil {
			msg = errors.GetMessage(err)
			return nil
		}

		q := d.composeUpdateQuery()
		if _, err := dbExec(q); err != nil {
			msg = errors.GetMessage(err)
			return nil
		}

		if msg == "" {
			fileHandleMsgMap[info.Name()] = msg
		} else {
			fileHandleMsgMap[info.Name()] = "uploaded"
		}
		return nil
	})
	return fileHandleMsgMap
}

//images foldname-stockid
func handleUploadedZipImagesJewelry(tempDir string) map[string]string {
	fileHandleMsgMap := make(map[string]string)
	filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		var msg string
		j := jewelry{}
		if !info.IsDir() {
			msg = fmt.Sprintf("SKIP file %s", info.Name())
			util.Println(msg)
			return nil
		}
		var filenames []string
		files, _ := ioutil.ReadDir(info.Name())
		for _, file := range files {
			if file.IsDir() {
				util.Printf("SKIP folder %s/%s", info.Name(), file)
			} else {
				bs, err := ioutil.ReadFile(path)
				if err != nil {
					msg = errors.GetMessage(err)
					return nil
				}
				if !filetype.IsImage(bs) {
					msg = fmt.Sprintf("%s  is not image", info.Name())
					util.Println(msg)
					return nil
				}
				ext := filepath.Ext(info.Name())
				if ext == "" {
					msg = fmt.Sprintf("SKIP file %s as the file has no extension", info.Name())
					util.Println(msg)
				} else {
					filename := fmt.Sprintf("beyoudiamond-image-%d.%s", time.Now().UnixNano(), filepath.Ext(info.Name()))
					if err := handleImage(path, filepath.Join(".image", "jewelry", filename)); err != nil {
						msg = errors.GetMessage(err)
					} else {
						filenames = append(filenames, filename)
					}
				}
			}
		}
		j.Images = filenames
		j.StockID = strings.Split(info.Name(), "-")[0]
		if err := j.isJewelryExistByStockID(); err != nil {
			msg = errors.GetMessage(err)
			return nil
		}

		q := j.composeUpdateQuery()
		if _, err := dbExec(q); err != nil {
			msg = errors.GetMessage(err)
			return nil
		}

		if msg == "" {
			fileHandleMsgMap[info.Name()] = msg
		} else {
			fileHandleMsgMap[info.Name()] = "uploaded"
		}
		return nil
	})
	return fileHandleMsgMap
}

//images foldname-stockid
func handleUploadedZipImagesGem(tempDir string) map[string]string {
	fileHandleMsgMap := make(map[string]string)
	filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		var msg string
		g := gem{}
		if !info.IsDir() {
			msg = fmt.Sprintf("SKIP file %s", info.Name())
			util.Println(msg)
			return nil
		}
		var filenames []string
		files, _ := ioutil.ReadDir(info.Name())
		for _, file := range files {
			if file.IsDir() {
				util.Printf("SKIP folder %s/%s", info.Name(), file)
			} else {
				bs, err := ioutil.ReadFile(path)
				if err != nil {
					msg = errors.GetMessage(err)
					return nil
				}
				if !filetype.IsImage(bs) {
					msg = fmt.Sprintf("%s  is not image", info.Name())
					util.Println(msg)
					return nil
				}
				ext := filepath.Ext(info.Name())
				if ext == "" {
					msg = fmt.Sprintf("SKIP file %s as the file has no extension", info.Name())
					util.Println(msg)
				} else {
					filename := fmt.Sprintf("beyoudiamond-image-%d.%s", time.Now().UnixNano(), filepath.Ext(info.Name()))
					if err := handleImage(path, filepath.Join(".image", "gem", filename)); err != nil {
						msg = errors.GetMessage(err)
					} else {
						filenames = append(filenames, filename)
					}
				}
			}
		}
		g.Images = filenames
		g.StockID = strings.Split(info.Name(), "-")[0]
		if err := g.isGemExistByStockID(); err != nil {
			msg = errors.GetMessage(err)
			return nil
		}

		q := g.composeUpdateQuery()
		if _, err := dbExec(q); err != nil {
			msg = errors.GetMessage(err)
			return nil
		}

		if msg == "" {
			fileHandleMsgMap[info.Name()] = msg
		} else {
			fileHandleMsgMap[info.Name()] = "uploaded"
		}
		return nil
	})
	return fileHandleMsgMap
}
