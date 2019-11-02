package model

import (
	utils "github.com/saichler/utils/golang"
	"os"
)

const (
	MAX_PART_SIZE = 1024 * 1024 * 5
)

type FileData struct {
	part int
	data []byte
}

func NewFileData(descriptor *FileDescriptor) *FileData {
	loc := MAX_PART_SIZE * int64(descriptor.part)
	fileData := &FileData{}
	fileData.part = descriptor.part

	dataSize := MAX_PART_SIZE
	if descriptor.size-loc < MAX_PART_SIZE {
		dataSize = int(descriptor.size - loc)
	}
	data := make([]byte, dataSize)

	file, _ := os.Open(descriptor.path)
	_, e := file.Seek(loc, 0)
	if e != nil {
		panic(e)
	}

	_, e = file.Read(data)
	if e != nil {
		panic(e)
	}
	file.Close()
	fileData.data = data

	dest, _ := os.Create("/tmp/test.zip")
	dest.Write(data)
	dest.Close()
	return fileData
}

func (fileData *FileData) Marshal() []byte {
	bs := utils.NewByteSlice()
	bs.AddInt(fileData.part)
	bs.AddByteSlice(fileData.data)
	return bs.Data()
}

func (fileData *FileData) Unmarshal(data []byte) {
	bs := utils.NewByteSliceWithData(data, 0)
	fileData.part = bs.GetInt()
	fileData.data = bs.GetByteSlice()
}

func (fileData *FileData) Part() int {
	return fileData.part
}

func (fileData *FileData) Data() []byte {
	return fileData.data
}
