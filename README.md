# tmpl

A unix command for generating files from templates. The syntax and
purpose are very similar to `jekyll`, but this cli is made with the
unix philosphy of doing one thing and doing it well. This cli
simply interprets templates and generates files from them.

The program is written in `Go` and uses it's
[standard library templating engine](https://golang.org/pkg/text/template/).
The `Go` docs are a good starting point for those unfamiliar with the template
syntax and semantics. Some extra functions are also provided for use within
the templates (namely `sh`).

## Running

```text
tmpl    [-t template] [infile [outfile]]
tmpl    -e [-t template] [outfile]
```

In general the program needs 3 files: the template, an input, and an output.
When run with no arguments, file named `template` with any extension is
searched for in all parent directories, and input and output is done
through stdin and stdout. A user defined template file can be given via
the `-t` flag.

## Example

This examples generates a site map of a git repository.
`tmpl.html`
```html
<!DOCTYPE html>
<html>
<head>
  <title>index for {{ sh `pwd` | join "" | base }}</title>
</head>
<body>
  {{ range sh `git ls-files` -}}
    <a href="{{ . }}">{{ . }}</a> </br>
  {{ end }}
</body>
</html>
```

Since this template does not take any file as input, make sure to use
the `-e` flag.

```bash
tmpl -e
```

## The File object

The file object is the default value of `.` (when not run with the `-e` flag).
due to limitations in the templating library and complexity, it should be
noted that the entire file is read into memory.

* `Name` - the name of the file
* `Body` - the contents of the file without front matter
* `Mtime` - the modification time
* `Fm` - the front matter as a `map[string]interface{}`
* `Md` - the markdown rendering of the `Body`

### Front matter

Front matter is `yaml` formated meta data at the top of a markdown file.
Currently, `tmpl` looks for front matter in all files and it must
be the very first thing in the file in betwen `---` or `===`.

## Functions

This is likely not updated so, the file `func_map.go` is the ultimate source of
truth. Some important ones are listed below:
* `func sh(string) []string` - takes in a shell script and returns the
  individual lines as elements in an array.
* `func ls([]string) []string` - behaves like the ls shell command, but faster
  since it doesn't call the shell.
* `func open(string) *File` - returns a file object, yes it leaks file
  yes it leaks file descriptors :(
* `func split(sep, a string)` - splits a at every sep

**note on string functions:** the normal convention of having `(a, sp)` params
are swapped to make use of template pipelines. The `a` will usually be some
other function's output and pipelines are tacked as the last parameter to a
function.

