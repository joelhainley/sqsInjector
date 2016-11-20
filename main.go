package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type siConfig struct {
	MessagePath string
	QueueURL    string
	Region      string
}

/*
- get list of files from dir
- read each file & output after loading into single string
- connect to aws sqs
- connect to queue
- put message in queue
- stitch it all together
- delay should be configurable?
- add wait time to loop?
- cli configurations should override config file settings
? possible to sync and not just append all to the queue
? opt libraries for go apps
*/

var showHelp = false
var configFile string

func init() {
	flag.BoolVar(&showHelp, "h", false, "--help")
	flag.StringVar(&configFile, "c", "./sqsinjector.cfg", "--c pathToConfig")
}

func main() {
	flag.Parse()

	fmt.Println("Help Flag Value : ", showHelp)

	// check if config file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		fmt.Println("Configuration File not found: ", configFile)
		os.Exit(1)
	}
	configData := fileToString(configFile)

	var config siConfig

	if _, err := toml.Decode(configData, &config); err != nil {
		fmt.Println("Problem reading configuration file: ", err.Error())
		os.Exit(1)
	}

	fmt.Println("Message Path : ", config.MessagePath)
	fmt.Println("Queue URL : ", config.QueueURL)
	fmt.Println("Region : ", config.Region)

	files, _ := filepath.Glob(config.MessagePath)
	fileCount := len(files)
	if fileCount > 0 {
		conn := getConnection(config)
		fmt.Printf("# of Message Files Found : %v\n", len(files))
		fmt.Println("-------------------------------------------")
		for _, f := range files {
			contents := fileToString(f)
			injectMessageIntoQueue(conn, config.QueueURL, contents)
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

func getConnection(config siConfig) *sqs.SQS {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return nil
	}

	svc := sqs.New(sess, &aws.Config{Region: aws.String("us-west-2"),
		CredentialsChainVerboseErrors: aws.Bool(true)})

	return svc
}

func injectMessageIntoQueue(conn *sqs.SQS, queueURL string, message string) {
	var smi sqs.SendMessageInput
	smi.MessageBody = aws.String(message)
	smi.QueueUrl = aws.String(queueURL)
	smi.DelaySeconds = aws.Int64(1)

	resp, err := conn.SendMessage(&smi)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println("it would appear that we have encountered an error...")
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	//fmt.Println("Response from SQS:")
	fmt.Println(resp.GoString())
}
