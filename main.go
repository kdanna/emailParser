package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"net/mail"
	"os"
)

func main() {
	fmt.Println("entering mailparse package")
	//Parse parses the command-line flags from os.Args[1:]. Must be called after all flags are defined and before flags are accessed by the program.
	flag.Parse()
	filename := flag.Arg(0) // this is the file that will be passed as a cli argument ( this .tar.gz file containing emails )
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	// Then creates a New Reader that can access the gzip contents
	g, err := gzip.NewReader(file)
	if err != nil {
		log.Fatal(err)
	}

	// Then pass gzip reader to the tar reader
	t := tar.NewReader(g)

	// create the csv output to put parsed results from the input file
	csvfile, err := os.Create("output.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer csvfile.Close() // defer the close on the csvfile until the program is done executing

	messageCounter := 0                               // keep track of how many email messages we read
	csvDataOutter := make([][]string, messageCounter) // inorder to write to a csv we need a slice containing slices of strings. This is outter slice
	csvDataInner := make([]string, 100)               // this is inner slice

	for {
		h, err := t.Next()
		if err == io.EOF && h != nil {
			break
		}
		messageCounter++
		if err != nil {
			log.Fatal(err)
		}
		if h.Typeflag == tar.TypeDir {
			continue
		}

		message, err := mail.ReadMessage(t)
		if err != nil {
			log.Fatal(err)
		}

		for h, val := range message.Header {
			for _, headerValue := range val {
				// fmt.Printf("%s, %+v \n", h, headerValue)
				csvDataInner = append(csvDataInner, h, headerValue)

			}
		}

		csvDataOutter = append(csvDataOutter, csvDataInner)

		w := csv.NewWriter(csvfile)
		defer w.Flush()
		if err := w.Error(); err != nil {
			log.Fatalln("error writing csv:", err)
		}

		fmt.Printf("\n This is the csvDataOutter Slice %+v", csvDataOutter)

		err = w.WriteAll(csvDataOutter)
		if err != nil {
			log.Fatalln("something went wrong writing file", err)
		}

		// fmt.Printf("\n This is the csvDataOutter Slice %+v", csvDataOutter)
		// fmt.Printf("\n This is the csvDataInner Slice %+v", len(csvDataInner))
		// fmt.Println("\n count of messages", messageCounter)
	}

}
