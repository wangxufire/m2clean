package main

import (
	"github.com/wangxufire/m2clean/args"
	"github.com/wangxufire/m2clean/cleaner"
)

func main() {
	args.Parse()
	cleaner.Process()
}
