package model

import (
	utils "github.com/saichler/utils/golang"
	"io/ioutil"
	"os"
)

type FileDescriptor struct {
	name  string
	path  string
	size  int64
	hash  string
	part  int
	files []*FileDescriptor
}

func (fd *FileDescriptor) Name() string {
	return fd.name
}

func (fd *FileDescriptor) Path() string {
	return fd.path
}

func (fd *FileDescriptor) Size() int64 {
	return fd.size
}

func (fd *FileDescriptor) Part() int {
	return fd.part
}

func (fd *FileDescriptor) SetPart(p int) {
	fd.part = p
}

func (fd *FileDescriptor) Files() []*FileDescriptor {
	return fd.files
}

func (fd *FileDescriptor) Marshal() []byte {
	bs := utils.NewByteSlice()
	fd.marshal(bs)
	return bs.Data()
}

func (fd *FileDescriptor) marshal(bs *utils.ByteSlice) {
	bs.AddString(fd.name)
	bs.AddString(fd.path)
	bs.AddInt64(fd.size)
	bs.AddString(fd.hash)
	bs.AddInt(fd.part)
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
	fd.path = bs.GetString()
	fd.size = bs.GetInt64()
	fd.hash = bs.GetString()
	fd.part = bs.GetInt()
	fd.files = make([]*FileDescriptor, 0)
	size := bs.GetInt()
	for i := 0; i < size; i++ {
		fdChild := &FileDescriptor{}
		fdChild.unmarshal(bs)
		fd.files = append(fd.files, fdChild)
	}
}

func Create(path string, dept, current int) (*FileDescriptor, error) {
	fi, e := os.Stat(path)
	if e != nil {
		return nil, e
	}
	fd := &FileDescriptor{}
	fd.name = fi.Name()
	fd.size = fi.Size()
	if current == 0 {
		fd.path = path
	}
	if fi.IsDir() && current < dept {
		fd.files = make([]*FileDescriptor, 0)
		files, e := ioutil.ReadDir(path)
		if e == nil {
			current++
			for _, file := range files {
				fdChild, e := Create(path+"/"+file.Name(), dept, current)
				if e == nil {
					fd.files = append(fd.files, fdChild)
				}
			}
		}
	} else {
		fd.part = int(fd.size/MAX_PART_SIZE) + 1
	}
	return fd, nil
}
