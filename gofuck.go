package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

var (
	filename string
	file     []byte
	noCells  int
	cells    []byte
	buff     bool
	intp     bool
	err      error

	cellIndex = 0
	fileIndex = 0
)

func main() {
	log.SetFlags(0)

	flag.StringVar(&filename, "file", "", "Path to .bf file")
	flag.IntVar(&noCells, "cells", 30000, "Set no. of cells to be used")
	flag.BoolVar(&buff, "v", false, "Verbose, prints the full buffer after execution")
	flag.BoolVar(&intp, "i", false, "Integer print, prints cell integer value instead of ascii")
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
			if !intp {
				fmt.Print(string(*cell))
			} else {
				fmt.Print(*cell, " ")
			}
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
	bts, err := reader.ReadBytes('\n') //TODO: handle for windowws \r\n
	if err != nil {
		log.Fatal("Can't read byte!")
	}
	if l := len(bts); l > 2 && bts[0] == '\\' {
		sesc, err := strconv.Unquote("\"" + string(bts[:l-1]) + "\"")
		if err != nil {
			log.Fatal("Can't read escape sequence!")
		}
		esc := sesc[0]
		return &esc
	}
	return &bts[0]
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
