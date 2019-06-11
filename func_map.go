package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

var FuncMap = template.FuncMap{
	"sh": func(cmd string, args ...string) (string, error) {
		out, err := exec.Command(cmd, args...).Output()
		return string(out), err
	},
	"split": strings.Split,
	"open":  NewFile,
	"ls": func(fnames ...string) ([]string, error) {
		ret := []string{}
		for _, v := range fnames {
			f, err := os.Open(fname)
			if err != nil {
				return nil, err
			}
			stat, err := f.Stat()
			if err != nil {
				return nil, err
			}

			if !stat.IsDir() {
				return []string{fname}, nil
			}

			arr, err := f.Readdirnames(0)
			if err != nil {
				return nil, err
			}

			ret = append(ret, arr...)
		}
		return ret, nil
	},
	"noext": func(fname string) string {
		return strings.TrimSuffix(fname, filepath.Ext(fname))
	},
}
