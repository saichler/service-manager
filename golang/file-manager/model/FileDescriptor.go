package model

import (
	utils "github.com/saichler/utils/golang"
	"io/ioutil"
	"os"
)

type FileDescriptor struct {
	name  string
	size  int64
	hash  string
	files []*FileDescriptor
}

func (fd *FileDescriptor) Name() string {
	return fd.name
}

func (fd *FileDescriptor) Size() int64 {
	return fd.size
}

func (fd *FileDescriptor) Files() []*FileDescriptor {
	return fd.files
}

func (fd *FileDescriptor) Marshal() []byte {
	bs := utils.NewByteSlice()
	return bs.Data()
}

func (fd *FileDescriptor) marshal(bs *utils.ByteSlice) {
	bs.AddString(fd.name)
	bs.AddInt64(fd.size)
	bs.AddString(fd.hash)
	if fd.files == nil {
		bs.AddInt(0)
	} else {
		bs.AddInt(len(fd.files))
		for _, fdChild := range fd.files {
			fdChild.marshal(bs)
		}
	}
}

func (fd *FileDescriptor) Unmarshal(data []byte) {
	bs := utils.NewByteSliceWithData(data, 0)
	fd.unmarshal(bs)
}

func (fd *FileDescriptor) unmarshal(bs *utils.ByteSlice) {
	fd.name = bs.GetString()
	fd.size = bs.GetInt64()
	fd.hash = bs.GetString()
	fd.files = make([]*FileDescriptor, 0)
	size := bs.GetInt()
	for i := 0; i < size; i++ {
		fdChild := &FileDescriptor{}
		fdChild.unmarshal(bs)
		fd.files = append(fd.files, fdChild)
	}
}

func Create(path string) (*FileDescriptor, error) {
	fi, e := os.Stat(path)
	if e != nil {
		return nil, e
	}
	fd := &FileDescriptor{}
	fd.name = fi.Name()
	fd.size = fi.Size()
	if fi.IsDir() {
		fd.files = make([]*FileDescriptor, 0)
		files, e := ioutil.ReadDir(path)
		if e == nil {
			for _, file := range files {
				fdChild, e := Create(path + "/" + file.Name())
				if e == nil {
					fd.files = append(fd.files, fdChild)
				}
			}
		}
	}
	return fd, nil
}
