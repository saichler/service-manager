package model

import (
	"bytes"
	"github.com/saichler/security"
	utils "github.com/saichler/utils/golang"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

type FileDescriptor struct {
	name      string
	parent    *FileDescriptor
	size      int64
	hash      string
	parts     int
	files     map[string]*FileDescriptor
	pathCache string
}

func (fileDescriptor *FileDescriptor) Name() string {
	return fileDescriptor.name
}

func (fileDescriptor *FileDescriptor) Hash() string {
	return fileDescriptor.hash
}

func (fileDescriptor *FileDescriptor) Parent() *FileDescriptor {
	return fileDescriptor.parent
}

func (fileDescriptor *FileDescriptor) Size() int64 {
	return fileDescriptor.size
}

func (fileDescriptor *FileDescriptor) Parts() int {
	return fileDescriptor.parts
}

func (fileDescriptor *FileDescriptor) SetPart(p int) {
	fileDescriptor.parts = p
}

func (fileDescriptor *FileDescriptor) Files() map[string]*FileDescriptor {
	return fileDescriptor.files
}

func (fileDescriptor *FileDescriptor) Root() *FileDescriptor {
	if fileDescriptor.parent == nil {
		return fileDescriptor
	}
	return fileDescriptor.parent.Root()
}

func (fileDescriptor *FileDescriptor) Get(path string) *FileDescriptor {
	index := strings.Index(path, "/")
	if index == 0 {
		path = path[1:]
		index = strings.Index(path, "/")
	}
	if index == -1 {
		index = len(path)
	}
	name := path[0:index]
	child := fileDescriptor.files[name]
	if child != nil && index != len(path) {
		return child.Get(path[index+1:])
	}
	return child
}

func (fileDescriptor *FileDescriptor) FullPath() string {
	if fileDescriptor.pathCache != "" {
		return fileDescriptor.pathCache
	}
	if fileDescriptor.parent == nil {
		return "/"
	}
	buff := &bytes.Buffer{}
	fileDescriptor.fullPath(buff)
	fileDescriptor.pathCache = buff.String()
	return fileDescriptor.pathCache
}

func (fileDescriptor *FileDescriptor) fullPath(buff *bytes.Buffer) {
	if fileDescriptor.parent == nil {
		return
	}
	fileDescriptor.parent.fullPath(buff)
	buff.WriteString("/")
	buff.WriteString(fileDescriptor.name)
}

func (fileDescriptor *FileDescriptor) Marshal() []byte {
	bs := utils.NewByteSlice()
	bs.AddString(fileDescriptor.FullPath())
	fileDescriptor.Root().marshal(bs)
	return bs.Data()
}

func (fileDescriptor *FileDescriptor) marshal(bs *utils.ByteSlice) {
	bs.AddString(fileDescriptor.name)
	bs.AddInt64(fileDescriptor.size)
	bs.AddString(fileDescriptor.hash)
	bs.AddInt(fileDescriptor.parts)
	if fileDescriptor.files == nil {
		bs.AddInt(0)
	} else {
		bs.AddInt(len(fileDescriptor.files))
		for name, fdChild := range fileDescriptor.files {
			bs.AddString(name)
			fdChild.marshal(bs)
		}
	}
}

func UnmarshalFileDescriptor(data []byte) *FileDescriptor {
	bs := utils.NewByteSliceWithData(data, 0)
	path := bs.GetString()
	root := &FileDescriptor{}
	root.unmarshal(bs)
	child := root.Get(path)
	return child
}

func (fileDescriptor *FileDescriptor) unmarshal(bs *utils.ByteSlice) {
	fileDescriptor.name = bs.GetString()
	fileDescriptor.size = bs.GetInt64()
	fileDescriptor.hash = bs.GetString()
	fileDescriptor.parts = bs.GetInt()
	fileDescriptor.files = make(map[string]*FileDescriptor)
	size := bs.GetInt()
	for i := 0; i < size; i++ {
		key := bs.GetString()
		child := &FileDescriptor{}
		child.unmarshal(bs)
		fileDescriptor.files[key] = child
		child.parent = fileDescriptor
	}
}

func createEmpty(path string) *FileDescriptor {
	index := strings.Index(path, "/")
	descriptor := &FileDescriptor{}
	if index == 0 {
		descriptor.name = "/"
	} else if index == -1 {
		descriptor.name = path
		return descriptor
	} else {
		descriptor.name = path[0:index]
	}
	descriptor.files = make(map[string]*FileDescriptor)
	child := createEmpty(path[index+1:])
	child.parent = descriptor
	descriptor.files[child.name] = child
	return descriptor

}

func NewFileDescriptor(path string, dept int) *FileDescriptor {
	root := createEmpty(path)
	descriptor := root.Get(path)
	hashJobs := &HashJobs{}
	hashJobs.cond = sync.NewCond(&sync.Mutex{})
	fill(descriptor, dept, 0, hashJobs)
	//hashJobs.wait()
	return descriptor
}

func fill(descriptor *FileDescriptor, dept, current int, hashJobs *HashJobs) {
	file, e := os.Stat(descriptor.FullPath())
	if e != nil {
		return
	}

	descriptor.size = file.Size()

	if file.IsDir() && current < dept {
		descriptor.files = make(map[string]*FileDescriptor)
		files, e := ioutil.ReadDir(descriptor.FullPath())
		if e == nil {
			for _, file := range files {
				child := &FileDescriptor{}
				child.parent = descriptor
				child.name = file.Name()
				fill(child, dept, current+1, hashJobs)
				descriptor.files[child.name] = child
			}
		}
	} else if !file.IsDir() {
		descriptor.parts = int(descriptor.size/MAX_PART_SIZE) + 1
		/*
			hj := &HashJob{}
			hj.hashJobs = hashJobs
			hj.descriptor = descriptor

			hashJobs.cond.L.Lock()
			hashJobs.total++
			go hj.Run()
			hashJobs.cond.L.Unlock()
		*/
	}
}

type HashJobs struct {
	total int
	done  int
	cond  *sync.Cond
}

func (hashJobs *HashJobs) wait() {
	hashJobs.cond.L.Lock()
	if hashJobs.total > hashJobs.done {
		hashJobs.cond.Wait()
	}
	hashJobs.cond.L.Unlock()
}

type HashJob struct {
	hashJobs   *HashJobs
	descriptor *FileDescriptor
}

func (hj *HashJob) Run() {
	hj.descriptor.hash, _ = security.FileHash256(hj.descriptor.FullPath())
	hj.hashJobs.cond.L.Lock()
	defer hj.hashJobs.cond.L.Unlock()
	hj.hashJobs.done++
	if hj.hashJobs.done == hj.hashJobs.total {
		hj.hashJobs.cond.Broadcast()
	}
}
