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

var noinFlag = flag.Bool("e", false, "generate with no input file")
var tmplFlag = flag.String("t", "", "template file path")
var helpFlag = flag.Bool("h", false, "show help")

const usage = `%s	[-t template] [infile [outfile]]
%[1]s	-e [-t template] [outfile]
%[1]s	-h
`

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
	if err != nil {
		log.Fatalln(err)
	}

	var out io.Writer = os.Stdout
	var in *File

	args := flag.Args()
	if *noinFlag {
		switch len(args) {
		case 1:
			outf, err := os.OpenFile(args[0],
				os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644)
			if err != nil {
				log.Fatalln(err)
			}
			defer outf.Close()
			out = outf
		case 0:
			break
		default:
			log.Printf(usage, os.Args[0])
		}
	} else {
		switch len(args) {
		case 2:
			outf, err := os.OpenFile(args[1],
				os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644)
			if err != nil {
				log.Fatalln(err)
			}
			defer outf.Close()
			out = outf
			fallthrough
		case 1:
			in, err = NewFile(args[0])
			if err != nil {
				if len(args) == 2 {
					os.Remove(args[1])
				}
				log.Fatalln(err)
			}
		case 0:
			in, err = NewFileReader("stdin", os.Stdin)
			if err != nil {
				log.Fatalln(err)
			}
		default:
			log.Printf(usage, os.Args[0])
		}
	}

	err = tmpl.Execute(out, in)
	if err != nil {
		if len(args) == 2 {
			os.Remove(args[1])
		}
		log.Fatalln(err)
	}
}

func findTmpl() (*os.File, error) {
	path, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	for len(path) > 1 {
		file, err := os.Open(filepath.Join(path, "template"))
		if err == nil {
			return file, nil
		}

		path = filepath.Dir(path)
	}

	return nil, ErrTmplNotFound
}
