package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"unicode/utf8"

	"github.com/pborman/getopt"
)

/*
- get list of files from dir
- read each file & output after loading into single string
- connect to aws
- put message in queue
- stitch it all together
? possible to sync and not just append all to the queue
? opt libraries for go apps

*/
func main() {
	fmt.Println("running sqsInjector ")

	helpFlag := getopt.Bool('h', "display help")
	src := flag.String("src", "", "path to source file or directory, supports globbing")
	config := flag.String("config", "./default.config", "the configuration file to use for this")

	getopt.Parse()

	fmt.Println("help flag value : ", *helpFlag)

	//flag.Parse()

	if utf8.RuneCountInString(*src) == 0 {
		fmt.Println("you must supply a src argument to run this application")
	}

	fmt.Println("reading from : ", *src)
	fmt.Println("non-flag args : ", flag.Args())
	fmt.Println("config file : ", *config)

	files, _ := filepath.Glob(*src)
	fileCount := len(files)
	if fileCount > 0 {
		fmt.Printf("%v files found\n", len(files))
		for _, f := range files {
			fmt.Printf("\n\nPROCESSING: %v\n", f)
			contents := fileToString(f)
			injectMessageIntoQueue("testqueue", contents)
		}
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func fileToString(filename string) string {
	dat, err := ioutil.ReadFile(filename)
	check(err)
	return (string(dat))
}

func injectMessageIntoQueue(queueName string, message string) {
	fmt.Printf("QUEUE NAME: \t %v\nMESSAGE:\n%v", queueName, message)
}
