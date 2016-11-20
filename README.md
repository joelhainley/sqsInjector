# SQS Injector

This tool pushes text-files into an AWS SQS Message Queue, one message per file.

Written in Go

It took its inspiration from another java-based tool I wrote called amqInjector available here (<https://github.com/joelhainley/amqInjector>)

Todo:

- better usage documentation and help
- ability to pass additional information in config file
- ability to pass all configuration options as command line arguments
- probably some better error handling around certain methods
- shouldn't die if a config file is not found
