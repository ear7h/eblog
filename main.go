package main

import (
	"errors"
	"flag"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/ear7h/tmpl/log"
)

var ErrTmplNotFound = errors.New("template not found")

var tmplFlag = flag.String("t", "", "tempate file path")
var helpFlag = flag.Bool("h", false, "show help")

const usage = "%s\t[-t template] [infile [outfile]]\n%[1]s\t-h\n"

func main() {

	flag.Parse()

	if *helpFlag {
		log.Printf(usage, os.Args[0])
		log.Println("flags:")
		flag.PrintDefaults()
		os.Exit(0)
	}

	var tmplFile *os.File
	var err error

	if len(*tmplFlag) != 0 {
		tmplFile, err = os.Open(*tmplFlag)
	} else {
		tmplFile, err = findTmpl()
	}
	if err != nil {
		log.Fatalln(err)
	}
	byt, err := ioutil.ReadAll(tmplFile)
	if err != nil {
		log.Fatalln(err)
	}

	tmpl, err := template.New("template").
		Funcs(FuncMap).
		Parse(string(byt))

	var out io.Writer = os.Stdout
	var in *File

	args := flag.Args()
	switch len(args) {
	case 2:
		outf, err := os.Open(args[1])
		if err != nil {
			log.Fatalln(err)
		}
		defer outf.Close()
		out = outf
		fallthrough
	case 1:
		in, err = NewFile(args[0])
		if err != nil {
			log.Fatalln(err)
		}
	case 0:
		in, err = NewFileStdin()
		if err != nil {
			log.Fatalln(err)
		}
	default:
		log.Printf(usage, os.Args[0])
	}

	err = tmpl.Execute(out, in)
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
		file, err := os.Open(filepath.Join(path, ".template"))
		if err == nil {
			return file, nil
		}

		path = filepath.Dir(path)
	}

	return nil, ErrTmplNotFound
}
