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
When run with no arguments, file named `template` is
searched for in all parent directories, and input and output is done
through stdin and stdout. A user defined template file can be given via
the `-t` flag.

## Example

This examples generates a site map of a git repository.
`template`
```html
<!DOCTYPE html>
<html>
<head>
  <title>index for {{ sh `pwd` | join "" | base }}</title>
</head>
<body>
  <ul>
  {{ range sh `git ls-files` -}}
    <li><a href="{{ . }}"> {{ . }} </a></li>
  {{ end -}}
  </ul>
</body>
</html>
```

A few things to note here:
* `sh` runs the `sh -C` command with a string argument and returns and *array*
  of strings, each string being a line (this makes iterating commands like `ls`
  simpler)
* `join ""` joins the strings in the array, making it a single string
* `base` returns the name of the directory, with other parts of the path
  removed.
* `range sh ..` here is an example where sh returning arrays is helpful

Since this template does not take any file as input, make sure to use
the `-e` flag.

```bash
tmpl -e
```

The output of running the above command on this directory should be
something like:

```html
<!DOCTYPE html>
<html>
<head>
  <title>index for tmpl</title>
</head>
<body>
  <ul>
  <li><a href=".gitignore"> .gitignore </a></li>
  <li><a href="README.md"> README.md </a></li>
  <li><a href="file.go"> file.go </a></li>
  <li><a href="func_map.go"> func_map.go </a></li>
  <li><a href="go.mod"> go.mod </a></li>
  <li><a href="go.sum"> go.sum </a></li>
  <li><a href="log/log.go"> log/log.go </a></li>
  <li><a href="main.go"> main.go </a></li>
  <li><a href="test/exif_test.jpg"> test/exif_test.jpg </a></li>
  <li><a href="test/exif_test.txt.tmpl"> test/exif_test.txt.tmpl </a></li>
  </ul>
</body>
</html>
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

