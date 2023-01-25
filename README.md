# GoVersion
### _The only Mass URL wordpress version checker you need_

A mass wordpress url version checker written in Golang.

## Features

- Unlimited Threads
- Custom URL Files processing.
- Checks per second counter
- Speeds over 50,000 CPM
- Proxyless
- Simple CLI

Introducing GoVersion, the ultimate WordPress version checking tool written in Golang. With GoVersion, you can easily and efficiently check the version of multiple WordPress sites in one go. Our tool is designed to save you time and effort by automating the process of checking for updates and ensuring that your sites are always running the latest version. 

## Installation

Dillinger requires [Golang](https://go.dev/) to run.

To Build
```sh
cd {directory}
go build main.go
```

To run
```sh
cd {directory}
go run main.go
```

## How to use

```sh
cd {directory}
main.exe {url_file} {thread_amount}

You can run the script by passing the file name and thread count as command line arguments. For example: `go run script.go file.txt 10`

This will process the URLs in the file `file.txt` 10 at a time, using goroutines. The URLs that use WordPress and their versions will be written to a file named "valid.txt", and the URLs that do not use WordPress will be written to a file named "invalid.txt".
```
