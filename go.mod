module github.com/ear7h/tmpl

go 1.12

require (
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	gopkg.in/russross/blackfriday.v2 v2.0.0-00010101000000-000000000000
	gopkg.in/yaml.v2 v2.2.2
)

replace gopkg.in/russross/blackfriday.v2 => github.com/russross/blackfriday/v2 v2.0.1
