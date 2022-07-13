package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/adrium/goheif/heic2jpg"
)

func main() {
	flag.Parse()
	if flag.NArg() != 2 {
		fmt.Fprintf(os.Stderr, "usage: heic2jpg <in-file> <out-file> \n")
		os.Exit(1)
	}
	fin, fout := flag.Arg(0), flag.Arg(1)

	converter := heic2jpg.NewConverter()
	if err := converter.HEIC2JPG(fin, fout); err != nil {
		log.Fatal(err)
	}
	log.Printf("Converted %q to %q successfully\n", fin, fout)
}
