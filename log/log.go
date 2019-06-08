package log

import (
	"fmt"
	"os"
)

var out = os.Stdout

func Fatal(v ...interface{}) {
	fmt.Print(v...)
	os.Exit(1)
}

func Fatalln(v ...interface{}) {
	fmt.Println(v...)
	os.Exit(1)
}

func Fataf(s string, args ...interface{}) {
	fmt.Fprintf(out, s, args...)
	os.Exit(1)
}

func Print(v ...interface{}) {
	fmt.Print(v...)
}

func Println(v ...interface{}) {
	fmt.Println(v...)
}

func Printf(s string, args ...interface{}) {
	fmt.Fprintf(out, s, args...)
}
