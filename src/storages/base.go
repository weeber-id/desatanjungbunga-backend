package storages

import (
	"mime/multipart"

	"github.com/google/uuid"
)

// BaseObject minio storage
type BaseObject struct {
	FileName string
	PathName string
	File     []byte
	Size     int64
}

// LoadFromFileHeader to this struct
// Usually used from gin formdata
func (b *BaseObject) LoadFromFileHeader(fileHeader *multipart.FileHeader, pathName string, fileName string) error {
	b.FileName = fileName
	b.PathName = pathName
	b.Size = fileHeader.Size

	file, err := fileHeader.Open()
	defer file.Close()
	if err != nil {
		return err
	}

	buf := make([]byte, fileHeader.Size)
	_, err = file.Read(buf)
	if err != nil {
		return err
	}
	b.File = buf

	return nil
}

// LoadFromFileHeaderRandomName from gin file header with random filename
func (b *BaseObject) LoadFromFileHeaderRandomName(fileHeader *multipart.FileHeader, pathName string) error {
	return b.LoadFromFileHeader(fileHeader, pathName, uuid.New().String())
}
