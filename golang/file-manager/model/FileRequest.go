package model

import utils "github.com/saichler/utils/golang"

type FileRequest struct {
	path string
	dept int
}

func NewFileRequest(path string, dept int) *FileRequest {
	fr := &FileRequest{}
	fr.path = path
	fr.dept = dept
	return fr
}

func (fr *FileRequest) Path() string {
	return fr.path
}

func (fr *FileRequest) Dept() int {
	return fr.dept
}

func (fr *FileRequest) Marshal() []byte {
	bs := utils.NewByteSlice()
	bs.AddString(fr.path)
	bs.AddInt(fr.dept)
	return bs.Data()
}

func (fr *FileRequest) UnMarshal(data []byte) {
	bs := utils.NewByteSliceWithData(data, 0)
	fr.path = bs.GetString()
	fr.dept = bs.GetInt()
}
