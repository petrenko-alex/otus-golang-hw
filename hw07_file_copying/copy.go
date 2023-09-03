package main

import (
	"errors"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

type FileCopier struct {
	fromPath, toPath string
	offset, limit    int64

	fromFile *os.File
}

func NewFileCopier(fromPath, toPath string, offset, limit int64) *FileCopier {
	return &FileCopier{
		fromPath: fromPath,
		toPath:   toPath,
		offset:   offset,
		limit:    limit,
	}
}

func (fc *FileCopier) Copy() error {
	fromFile, openingErr := fc.openFromFile()
	if openingErr != nil {

	}
	defer fromFile.Close()

	buffer := fc.getBufferForFile(fromFile)

	_, readErr := fromFile.ReadAt(buffer, fc.offset)

	if readErr != nil {
		//todo: realize
		if readErr == io.EOF {
			// что если не дочитали ?
		}
	}

	toFile, fileCreateErr := fc.createToFile()
	if fileCreateErr != nil {
		// todo: realize
	}
	defer toFile.Close()

	_, writeErr := toFile.Write(buffer)
	if writeErr != nil {
		// todo: realize
	}

	return nil
}

func (fc *FileCopier) openFromFile() (*os.File, error) {
	file, err := os.Open(fc.fromPath)
	if err != nil {
		if os.IsNotExist(err) {
			// файл не найден
		}
		// другие ошибки, например нет прав
		// todo: wrap error open file
		return nil, err
	}

	return file, nil
}

func (fc *FileCopier) createToFile() (*os.File, error) {
	toFile, err := os.Create(fc.toPath)
	if err != nil {
		// todo: wrap error
		return nil, err
	}

	return toFile, nil
}

func (fc *FileCopier) getBufferForFile(file *os.File) []byte {
	fileInfo, err := file.Stat()
	if err != nil {
		// todo: realize
	}

	fileSize := fileInfo.Size()
	buffSize := fc.limit
	if buffSize <= 0 {
		buffSize = fileInfo.Size() - offset
	} else if buffSize > fileSize {
		buffSize = fileSize
	}

	if offset+limit > fileSize {
		buffSize = fileSize - offset
	}

	return make([]byte, buffSize)
}
