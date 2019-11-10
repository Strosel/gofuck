package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var (
	filename string
	file     []byte
	noCells  int
	cells    []byte
	buff     bool
	err      error

	cellIndex = 0
	fileIndex = 0
)

func main() {
	log.SetFlags(0)

	flag.StringVar(&filename, "file", "", "Path to .bf file")
	flag.IntVar(&noCells, "cells", 30000, "Set n0. of cells used")
	flag.BoolVar(&buff, "v", false, "Verbose, prints the full buffer after execution")
	flag.Parse()

	if filename == "" {
		log.Fatal("No file!")
	} else if filename[len(filename)-3:] != ".bf" {
		log.Fatal("Given file is not .bf")
	}

	cells = make([]byte, noCells)

	file, err = ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("Error reading file!")
	}

	defer func() {
		r := recover()
		if r != nil {
			fmt.Println("Error on line", lineNo(), "of", filename)
			fmt.Println(r)
		}
		if buff {
			fmt.Println("\nCursor at:", cellIndex)
			fmt.Println(cells)
		}
	}()

	for fileIndex < len(file) {
		cc := file[fileIndex]
		cell := &(cells[cellIndex])
		switch cc {
		case '>':
			if cellIndex != noCells-1 {
				cellIndex++
			}
			fileIndex++
		case '<':
			if cellIndex != 0 {
				cellIndex--
			}
			fileIndex++
		case '+':
			*cell++
			fileIndex++
		case '-':
			*cell--
			fileIndex++
		case '.':
			fmt.Print(string(*cell))
			fileIndex++
		case ',':
			*cell = *input()
			fileIndex++
		case '[':
			if *cell == 0 {
				findClose()
			}
			fileIndex++
		case ']':
			if *cell != 0 {
				findOpen()
			}
			fileIndex++
		default:
			fileIndex++
		}
	}
}

func findClose() {
	track := 0
	for i := fileIndex + 1; i < len(file); i++ {
		c := file[i]
		if c == '[' {
			track++
		} else if c == ']' && track != 0 {
			track--
		} else if c == ']' {
			fileIndex = i
			return
		}
	}
}

func findOpen() {
	track := 0
	for i := fileIndex - 1; i > 0; i-- {
		c := file[i]
		if c == ']' {
			track++
		} else if c == '[' && track != 0 {
			track--
		} else if c == '[' {
			fileIndex = i
			return
		}
	}
}

func input() *byte {
	reader := bufio.NewReader(os.Stdin)
	bt, err := reader.ReadByte()
	if err != nil {
		log.Fatal("Can't read byte!")
	}
	return &bt
}

func lineNo() int {
	out := 1
	for i := 0; i < fileIndex; i++ {
		if file[i] == '\n' {
			out++
		}
	}
	return out
}
