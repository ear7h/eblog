package main

import (
	"errors"
	"flag"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ear7h/eblog/log"
)

var ErrTmplNotFound = errors.New("template not found")

var tmplFlag = flag.String("t", "", "tempate file path")
var helpFlag = flag.Bool("h", false, "show help")

const usage = "%s\t[-t template] [infile [outfile]]\n%[1]s\t-h\n"

func main() {
	var tmpl io.Reader
	var tmplFile *os.File
	var err error

	flag.Parse()


	if *helpFlag {
		log.Printf(usage, os.Args[0])
		log.Println("flags:")
		flag.PrintDefaults()
		os.Exit(0)
	}

	if len(*tmplFlag) != 0 {
		tmplFile, err = os.Open(*tmplFlag)
	} else {
		tmplFile, err = findTmpl()
	}

	if err != nil {
		log.Fatalln(err)
	}

	defer tmplFile.Close()
	tmpl = tmplFile

	post := Post{}
	var out io.Writer

	out = os.Stdout
	post.body = os.Stdin
	post.Title = "stdin"

	args := flag.Args()
	switch len(args) {
	case 0:
		break
	case 2:
		outf, err := os.Open(args[1])
		if err != nil {
			log.Fatalln(err)
		}
		outf.Close()
		out = outf
		fallthrough
	case 1:
		inf, err := os.Open(args[0])
		if err != nil {
			log.Fatalln(err)
		}
		defer inf.Close()
		post.body = inf
		base := filepath.Base(args[0])
		base = strings.TrimSuffix(base, filepath.Ext(base))
		post.Title = base
	default:
		log.Printf(usage, os.Args[0])
	}

	err = post.Compile(tmpl, out)
	if err != nil {
		log.Fatalln(err)
	}
}

func findTmpl() (*os.File, error) {
	path, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	for len(path) > 1 {
		file, err := os.Open(filepath.Join(path, "e7blog.html"))
		if err == nil {
			return file, nil
		}

		path = filepath.Dir(path)
	}

	return nil, ErrTmplNotFound
}

type Post struct {
	Title string
	body  io.Reader
}

func (p *Post) Body() (string, error) {
	byt, err := ioutil.ReadAll(p.body)
	return string(byt), err
}

func (p *Post) Compile(tmplr io.Reader, out io.Writer) error {
	byt, err := ioutil.ReadAll(tmplr)
	if err != nil {
		return err
	}

	tmpl, err := template.New("page").Parse(string(byt))
	if err != nil {
		return err
	}

	return tmpl.Execute(out, &p)
}
