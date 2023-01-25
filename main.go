package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var totalUrls int
var doneUrls int
var cpm int

func checkWpVersion(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	if !strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		resp.Body.Close()
		return "", fmt.Errorf("%s is not a HTML page", url)
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), `<meta name="generator" content="WordPress `) {
			version := strings.Split(scanner.Text(), `<meta name="generator" content="WordPress `)[1]
			version = strings.Split(version, `"`)[0]
			fmt.Printf("The version of WordPress on %s is %s\n", url, version)
			resp.Body.Close()
			return url + "," + version, nil
		}
	}

	resp.Body.Close()

	if err := scanner.Err(); err != nil {
		return "", err
	}

	fmt.Printf("%s does not use WordPress.\n", url)
	return "", nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Not enough arguments provided")
		return
	}

	fileName := os.Args[1]
	threadCount, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Second argument should be an integer")
		return
	}

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var urls []string
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return
	}

	var wg sync.WaitGroup
	validFile, err := os.Create("valid.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer validFile.Close()

	invalidFile, err := os.Create("invalid.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer invalidFile.Close()

	var validUrls []string
	var invalidUrls []string
	for i := 0; i < len(urls); i += threadCount {
		end := i + threadCount
		if end > len(urls) {
			end = len(urls)
		}

		wg.Add(end - i)
		for _, url := range urls[i:end] {
			go func(url string) {
				defer wg.Done()
				result, err := checkWpVersion(url)
				if err != nil {
					fmt.Println(err)
					return
				}
				if result != "" {
					validUrls = append(validUrls, result)
				} else {
					invalidUrls = append(invalidUrls, url)
				}
			}(url)
		}
		wg.Wait()
	}

	for _, url := range validUrls {
		validFile.WriteString(url + "\n")
	}
	for _, url := range invalidUrls {
		invalidFile.WriteString(url + "\n")
	}
	totalUrls = len(urls)
	fmt.Printf("WordPress Version checker | %d/%d | CPM: %d\n", totalUrls, doneUrls, cpm)
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				cpm = doneUrls - cpm
				fmt.Printf("WordPress Version checker | %d/%d | CPM: %d\n", totalUrls, doneUrls, cpm)
				cpm = doneUrls
			}
		}
	}()

	for i := 0; i < len(urls); i += threadCount {
		end := i + threadCount
		if end > len(urls) {
			end = len(urls)
		}

		wg.Add(end - i)
		for _, url := range urls[i:end] {
			go func(url string) {
				defer wg.Done()
				result, err := checkWpVersion(url)
				if err != nil {
					fmt.Println(err)
					return
				}
				if result != "" {
					validUrls = append(validUrls, result)
				} else {
					invalidUrls = append(invalidUrls, url)
				}
				doneUrls++
			}(url)
		}
		wg.Wait()
	}
	ticker.Stop()
	fmt.Printf("WordPress Version checker | %d/%d | CPM: %d\n", totalUrls, doneUrls, cpm)

}
