package main

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"time"

	"gopkg.in/russross/blackfriday.v2"
	"gopkg.in/yaml.v2"
)

var ErrBadFrontMatter = errors.New("bad front matter")

func NewFile(fname string) (*File, error) {

	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	ret, err := NewFileReader(fname, f)
	if err != nil {
		return nil, err
	}

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	ret.Mtime = stat.ModTime()

	return ret, nil
}

func NewFileReader(name string, r io.Reader) (*File, error) {
	byt, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	fm := map[string]interface{}{}
	var arr [][]byte
	if bytes.HasPrefix(byt, []byte("---\n")) {
		arr = bytes.SplitN(byt, []byte("---\n"), 3)
	} else if bytes.HasPrefix(byt, []byte("===\n")) {
		arr = bytes.SplitN(byt, []byte("===\n"), 3)
	} else {
		goto L
	}

	if len(arr) != 3 {
		return nil, ErrBadFrontMatter
	}

	byt = arr[2]
	err = yaml.Unmarshal(arr[1], &fm)
	if err != nil {
		return nil, err
	}
L:

	return &File{
		Name:  name,
		Body:  string(byt),
		Mtime: time.Now(),
		Fm:    fm,
	}, nil
}

type File struct {
	Name  string
	Body  string
	Mtime time.Time
	/* TODO
	Ctime time.Time
	Atime time.Time
	*/
	Fm map[string]interface{}
}

func (f *File) Md() string {
	return string(blackfriday.Run([]byte(f.Body)))
}
