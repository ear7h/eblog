package main

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
	"gopkg.in/russross/blackfriday.v2"
	"gopkg.in/yaml.v2"
)

var ErrBadFrontMatter = errors.New("bad front matter")


type fileError struct {
	fname string
	err error
}

func (fe *fileError) Error() string {
	return fmt.Sprintf("%s: %s", fe.fname, fe.err.Error())
}

func NewFile(fname string) (*File, error) {

	f, err := os.Open(fname)
	if err != nil {
		return nil, &fileError{fname, err}
	}
	defer f.Close()

	ret, err := NewFileReader(fname, f)
	if err != nil {
		return nil, &fileError{fname, err}
	}

	stat, err := f.Stat()
	if err != nil {
		return nil, &fileError{fname, err}
	}

	ret.Mtime = stat.ModTime()

	return ret, nil
}

type walkerFunc func(name exif.FieldName, tag *tiff.Tag) error

func (f walkerFunc) Walk(name exif.FieldName, tag *tiff.Tag) error {
	return f(name, tag)
}


func NewFileReader(name string, r io.Reader) (*File, error) {
	byt, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	meta := map[string]interface{}{}
	var arr [][]byte

	// if image
	switch filepath.Ext(name) {
	case "png", "jpg", "tif":
		break
	default:
		goto PostExif
	}

	{
		// parse exif
		x, err := exif.Decode(bytes.NewReader(byt))
		if err != nil {
			return nil, err
		}

		f := func(name exif.FieldName, tag *tiff.Tag) error {
			meta[string(name)] = tag.String()
			return nil
		}

		err = x.Walk(walkerFunc(f))
		if err != nil {
			return nil, err
		}

		goto PostMeta
	}
PostExif:


	if bytes.HasPrefix(byt, []byte("---\n")) {
		// parse front matter
		arr = bytes.SplitN(byt, []byte("---\n"), 3)
	} else if bytes.HasPrefix(byt, []byte("===\n")) {
		// parse front matter
		arr = bytes.SplitN(byt, []byte("===\n"), 3)
	} else {
		goto PostMeta
	}

	if len(arr) != 3 {
		return nil, ErrBadFrontMatter
	}

	byt = arr[2]
	err = yaml.Unmarshal(arr[1], &meta)
	if err != nil {
		return nil, err
	}

PostMeta:

	return &File{
		Name:  name,
		Body:  string(byt),
		Mtime: time.Now(), // TODO
		Meta:  meta,
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
	Meta map[string]interface{}
}

func (f *File) Md() string {
	return string(blackfriday.Run([]byte(f.Body)))
}
