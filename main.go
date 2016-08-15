package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
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
	files, _ := filepath.Glob("/Users/joel/tmp/f*.js")
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
