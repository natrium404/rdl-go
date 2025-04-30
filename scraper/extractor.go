package scraper

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func extract(src, dest string) error {
	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		path := filepath.Join(dest, file.Name)

		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return err
		}

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			return err
		}

		output, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		readContent, err := file.Open()
		if err != nil {
			output.Close()
			return err

		}

		_, err = io.Copy(output, readContent)

		output.Close()
		readContent.Close()

		if err != nil {
			return err
		}
	}

	err = os.Remove(src)
	if err != nil {
		return err
	}
	return nil
}
