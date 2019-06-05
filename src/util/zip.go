package util

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func Unzip(zipPath, destDir string) (err error) {
	if !PathExists(zipPath) {
		return fmt.Errorf("no such zip file %s", zipPath)
	}

	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return
	}
	defer r.Close()

	for _, f := range r.File {
		err = copyZipFile(f, destDir)
		if err != nil {
			return
		}
	}
	return
}

func copyZipFile(f *zip.File, destDir string) (err error) {
	rc, err := f.Open()
	if err != nil {
		return
	}
	defer rc.Close()

	targetPath := filepath.Join(destDir, f.Name)

	if f.FileInfo().IsDir() {
		err = os.MkdirAll(targetPath, f.Mode())
		if err != nil {
			return
		}
		return
	}

	targetDir := filepath.Dir(targetPath)
	if !PathExists(targetDir) {
		err = os.MkdirAll(targetDir, 0755)
		if err != nil {
			return
		}
	}

	tf, err := os.OpenFile(targetPath, os.O_RDWR|os.O_CREATE, f.Mode())
	if err != nil {
		return
	}
	defer tf.Close()

	_, err = io.Copy(tf, rc)
	return
}

func Zip(src, dest string) error {
	if !PathExists(src) {
		return fmt.Errorf("no such dir %s", src)
	}

	fw, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer fw.Close()

	zw := zip.NewWriter(fw)
	defer zw.Close()

	return filepath.Walk(src, func(path string, fi os.FileInfo, _ error) error {
		if fi.IsDir() {
			return nil
		}
		return copyFileToZip(src, path, zw)
	})
}

func copyFileToZip(src, path string, zw *zip.Writer) error {
	rel, _ := filepath.Rel(src, path)

	if rel == "." {
		rel = path
	}

	fr, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fr.Close()

	f, err := zw.Create(rel)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, fr)
	return err
}
