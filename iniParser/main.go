package main

import (
	"flag"
	"log"
	"os"
)

var (
	iniFileName = flag.String("name", "", "name of the inifile")
)

func main() {
	flag.Parse()

	if *iniFileName == "" {
		log.Fatal("Argument name must be set")
	}

	log.Printf("Parse inifile: '%s'", *iniFileName)

	f, err := os.Open(*iniFileName)
	if err != nil {
		log.Fatalf("error opening file '%s': %v", *iniFileName, err)
	}
	defer f.Close()

	c, err := NewConfig(f)
	if err != nil {
		log.Fatalf("error creating config: %v", err)
	}
	log.Println("config read successfully!")

	for n, sec := range c.sections {
		log.Printf("Section '%s' with %d key-value pairs", n, len(sec.values))
	}

}
