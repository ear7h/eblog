package main

import (
	"os"
	"os/exec"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

func splitNoEmpty(s, sep string) []string {

		ret := strings.Split(s, sep)
		i := 0
		for _, v := range ret {
			if v == "" {
				continue
			}
			ret[i] = v
			i++
		}
		if i == 0 {
			return nil
		} else {
			return ret[:i]
		}
}

var FuncMap = template.FuncMap{
	"sh": func(cmds string) ([]string, error) {
		cmd := exec.Command("sh", "-c", cmds)
		cmd.Env = os.Environ()
		out, err := cmd.Output()
		if err != nil {
			errv, ok := err.(*exec.ExitError)
			if ok {
				return nil, fmt.Errorf(string(errv.Stderr))
			} else {
				return nil, err
			}
		}
		return splitNoEmpty(string(out), "\n"), nil
	},
	"split": func (sep, a string) []string {
      return splitNoEmpty(a, sep)
  },
	"open":  NewFile,
	"ls": func(fnames ...string) ([]string, error) {
		if len(fnames) == 0 {
			fnames = []string{"."}
		}
		ret := []string{}
		for _, fname := range fnames {
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

		sort.StringSlice(ret).Sort()
		return ret, nil
	},
	"noext": func(fname string) string {
		return strings.TrimSuffix(fname, filepath.Ext(fname))
	},
	"env": os.Getenv,
	"dir": filepath.Dir,
  "base": filepath.Base,
  "ext": filepath.Ext,
	"pjoin": filepath.Join,
  "join": sepFirst(strings.Join),
}

func sepFirst(f func([]string, string) string) func(string, []string) string {
    return func(sep string, a []string) string {
        return f(a, sep)
    }
}
