// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		name := os.Args[0]
		fmt.Printf("Usage: %v \"path/to/target/file\"\n", name)
		os.Exit(1)
	}
	fpath := os.Args[1]
	dat, err := ioutil.ReadFile(fpath)
	if err != nil {
		log.Fatal(err)
	}
	parser := &Parser{Buffer: string(dat)}
	parser.Init()
	if err := parser.Parse(); err != nil {
		log.Fatal(err)
	}
	parser.PrintSyntaxTree()
}
