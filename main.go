package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
)

var (
	cat    = flag.Bool("append", false, "append the output to the files rather than rewriting them")
	ignore = flag.Bool("ignore-interrupts", false, "ignore the SIGINT signal")
)

func handleFlags() int {
	flag.Parse()

	oflags := os.O_WRONLY | os.O_CREATE

	if *cat {
		oflags |= os.O_APPEND
	}

	if *ignore {
		signal.Ignore(os.Interrupt)
	}

	return oflags
}

const name = "tee"

func main() {
	oflags := handleFlags()

	files := make([]*os.File, 0, flag.NArg())
	writers := make([]io.Writer, 0, flag.NArg()+1)
	for _, fname := range flag.Args() {
		f, err := os.OpenFile(fname, oflags, 0o666)
		if err != nil {
			log.Fatalf("%s: error opening %s: %v", name, fname, err)
		}
		files = append(files, f)
		writers = append(writers, f)
	}
	writers = append(writers, os.Stdout)

	mw := io.MultiWriter(writers...)
	if _, err := io.Copy(mw, os.Stdin); err != nil {
		log.Fatalf("%s: error: %v", name, err)
	}

	for _, f := range files {
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "%s: error closing file %q: %v\n", name, f.Name(), err)
		}
	}
}