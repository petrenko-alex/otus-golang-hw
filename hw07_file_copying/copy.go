package main

import (
	"errors"
	"fmt"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrWithSrcFile           = errors.New("source file problem")
	ErrWithDestFile          = errors.New("destination file problem")
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
		return fmt.Errorf(ErrWithSrcFile.Error()+": %w", openingErr)
	}
	defer fromFile.Close()

	buffer, buffSizeErr := fc.getBufferForFile(fromFile)
	if buffSizeErr != nil {
		return fmt.Errorf(ErrWithSrcFile.Error()+": %w", openingErr)
	}

	_, readErr := fromFile.ReadAt(buffer, fc.offset)
	if readErr != nil {
		return fmt.Errorf(ErrWithSrcFile.Error()+": %w", openingErr)
	}

	toFile, fileCreateErr := fc.createToFile()
	if fileCreateErr != nil {
		return fmt.Errorf(ErrWithDestFile.Error()+": %w", fileCreateErr)
	}
	defer toFile.Close()

	_, writeErr := toFile.Write(buffer)
	if writeErr != nil {
		return fmt.Errorf(ErrWithDestFile.Error()+": %w", fileCreateErr)
	}

	return nil
}

func (fc *FileCopier) openFromFile() (*os.File, error) {
	file, err := os.Open(fc.fromPath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (fc *FileCopier) createToFile() (*os.File, error) {
	toFile, err := os.Create(fc.toPath)
	if err != nil {
		return nil, err
	}

	return toFile, nil
}

func (fc *FileCopier) getBufferForFile(file *os.File) ([]byte, error) {
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	fileSize := fileInfo.Size()
	if fileSize == 0 {
		return nil, ErrUnsupportedFile
	} else if fileSize < fc.offset {
		return nil, ErrOffsetExceedsFileSize
	}

	buffSize := fc.limit
	if buffSize <= 0 {
		buffSize = fileInfo.Size() - fc.offset
	} else if buffSize > fileSize {
		buffSize = fileSize
	}

	if fc.offset+fc.limit > fileSize {
		buffSize = fileSize - fc.offset
	}

	return make([]byte, buffSize), nil
}
