package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

func NewFile(fname string) (*File, error) {
	byt, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}

	fm := map[string]interface{}{}
	arr := bytes.SplitN(byt, []byte("---\n"), 1)
	if len(arr) == 1 {
		arr = bytes.SplitN(byt, []byte("===\n"), 1)
	}
	if len(arr) == 2 {
		err = yaml.Unmarshal(arr[1], &fm)
		if err != nil {
			return nil, err
		}
	}

	stat, err := os.Stat(fname)
	if err != nil {
		return nil, err
	}

	return &File{
		Name:  fname,
		Body:  string(byt),
		Mtime: stat.ModTime(),
		Fm:    fm,
	}, nil
}

func NewFileStdin() (*File, error) {
	byt, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return nil, err
	}

	fm := map[string]interface{}{}
	arr := bytes.SplitN(byt, []byte("---\n"), 1)
	if len(arr) == 1 {
		arr = bytes.SplitN(byt, []byte("===\n"), 1)
	}
	if len(arr) == 2 {
		err = yaml.Unmarshal(arr[1], &fm)
		if err != nil {
			return nil, err
		}
	}

	return &File{
		Name:  "stdin",
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
