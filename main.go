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
	// Parse parses the command-line flags from os.Args[1:]. Must be called after all flags are defined and before flags are accessed by the program.
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
	defer g.Close()

	// Then pass gzip reader to the tar reader
	t := tar.NewReader(g)

	messageCounter := 0
	headerValueCounter := 0                            // keep track of how many email messages are read
	csvDataOutter := make([][]string, messageCounter)  // in order to write to a csv we need a slice containing slices of strings. This is the outter slice
	csvDataInner := make([]string, headerValueCounter) // this is the inner slice

	// create an infinite loop that will break out once the end of the .tar file is reached EOF
	for {
		h, err := t.Next()

		if err == io.EOF {
			break // end of archive file
		}

		if err != nil {
			log.Fatal("next error", err)
		}

		if h.Typeflag == tar.TypeDir {
			continue // if the tar Type is a dir, continue
		}

		messageCounter++ // increment the counter
		// fmt.Printf("Contents of %+v:\n", h.Name)
		message, err := mail.ReadMessage(t)
		if err != nil {
			log.Fatal(err)
		}

		for h, val := range message.Header {
			for _, headerValue := range val {
				// fmt.Printf("%s, %+v \n", h, headerValue)
				headerValueCounter++
				csvDataInner = append(csvDataInner, h, headerValue)

			}
		}
	}

	fmt.Printf("\n messageeCounter, %d \n", messageCounter)
	fmt.Printf("\n headerValueCounter, %d \n", headerValueCounter)

	csvDataOutter = append(csvDataOutter, csvDataInner)

	// create the csv output to put parsed results from the input file
	csvfile, err := os.Create("output.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer csvfile.Close() // defer the close on the csvfile until the program is done executing
	w := csv.NewWriter(csvfile)
	if err := w.Error(); err != nil {
		log.Fatalln("error writing csv:", err)
	}

	fmt.Printf("\n This is the csvDataOutter Slice %+v", csvDataOutter)

	err = w.WriteAll(csvDataOutter)
	if err != nil {
		log.Fatalln("something went wrong writing file", err)
	}

}
