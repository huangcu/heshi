package util

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func PathExists(p string) bool {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return false
	}
	return true
}

func GetFileMd5sum(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	md5h := md5.New()
	io.Copy(md5h, f)
	return fmt.Sprintf("%x", md5h.Sum(nil)), nil
}

func CreateMd5sumFile(path string) (err error) {
	fsum, err := GetFileMd5sum(path)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(path+".md5sum", []byte(fsum), 0644)
	if err != nil {
		return
	}

	return
}

func Base64EncodePng(path string) (string, error) {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	s := base64.StdEncoding.EncodeToString(bs)
	return s, nil
}

func GetFileSize(zf string) (int64, error) {
	f, err := os.Open(zf)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}
