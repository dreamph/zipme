package main

import (
	//"archive/zip"
	"github.com/alexmullins/zip"
	"io"
	"os"
	"path/filepath"
)

func addFileToZip(w *zip.Writer, file string, baseDir string, password string) error {
	fileInfo, err := os.Stat(file)
	if err != nil {
		return err
	}

	// Create a header based on file info
	header, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		return err
	}

	// Create relative path for the file inside the zip
	relPath, err := filepath.Rel(baseDir, file)
	if err != nil {
		return err
	}
	header.Name = relPath

	// If it's a directory, add a trailing slash
	if fileInfo.IsDir() {
		header.Name += "/"
	}

	if password != "" {
		header.SetPassword(password)
	}

	writer, err := w.CreateHeader(header)
	if err != nil {
		return err
	}

	// If it's a file, copy its contents
	if !fileInfo.IsDir() {
		fileReader, err := os.Open(file)
		if err != nil {
			return err
		}
		defer fileReader.Close()

		_, err = io.Copy(writer, fileReader)
		if err != nil {
			return err
		}
	}

	return nil
}

// addDirectoryToZip recursively adds files from a directory to the zip archive.
func addDirectoryToZip(w *zip.Writer, dir string, password string) error {
	return filepath.Walk(dir, func(file string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return addFileToZip(w, file, dir, password)
	})
}

func newZipWriter(zipFile *os.File) *zip.Writer {
	zipWriter := zip.NewWriter(zipFile)
	return zipWriter
}

func zipFiles(zipFilePath string, filePaths []string, password string) error {
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := newZipWriter(zipFile)
	defer zipWriter.Close()

	for _, filePath := range filePaths {
		info, err := os.Stat(filePath)
		if err != nil {
			return err
		}

		if info.IsDir() {
			err := addDirectoryToZip(zipWriter, filePath, password)
			if err != nil {
				return err
			}
		} else {
			err := addFileToZip(zipWriter, filePath, filepath.Dir(filePath), password)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
