package main

import (
	"fmt"
	"runtime"
)

func debug(arg any, name string) {
	_, file, line, _ := runtime.Caller(1)
	fmt.Println(file, line)
	fmt.Println(name, ": ", arg)
}
