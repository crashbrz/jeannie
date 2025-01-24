package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
)

var (
	apiKey      string
	filePath    string
	endpoint    string
	removeColor bool
	showInvalid bool
	debug       bool
	numThreads  int
)

func init() {
	flag.StringVar(&apiKey, "k", "", "Single API key to validate")
	flag.StringVar(&filePath, "f", "", "File containing API keys, one per line")
	flag.StringVar(&endpoint, "e", "https://api.opsgenie.com/v1/json/cloudwatch?apiKey=", "Opsgenie API endpoint")
	flag.BoolVar(&removeColor, "remove-color", false, "Disable color output")
	flag.BoolVar(&showInvalid, "d", false, "Display invalid API keys")
	flag.BoolVar(&debug, "debug", false, "Display HTTP response body for each API key")
	flag.IntVar(&numThreads, "t", 1, "Number of goroutines to use for validation")
}

func validateAPIKey(apiKey string, endpoint string) (bool, string) {
	fullEndpoint := endpoint + apiKey
	req, err := http.NewRequest("GET", fullEndpoint, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return false, ""
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return false, ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return false, ""
	}

	return strings.Contains(string(body), "Could not authenticate"), string(body)
}

func printResult(apiKey string, isValid bool, body string) {
	if removeColor {
		if isValid {
			fmt.Printf("Invalid: %s\n", apiKey)
		} else {
			fmt.Printf("Valid: %s\n", apiKey)
		}
	} else {
		if isValid {
			fmt.Printf("\033[31mInvalid: %s\033[0m\n", apiKey)
		} else {
			fmt.Printf("\033[32mValid: %s\033[0m\n", apiKey)
		}
	}
	if debug {
		fmt.Printf("Response Body: %s\n", body)
	}
}

func worker(apiKeys <-chan string, results chan<- struct {
	apiKey  string
	isValid bool
	body    string
}) {
	for apiKey := range apiKeys {
		isValid, body := validateAPIKey(apiKey, endpoint)
		isValid = !isValid
		results <- struct {
			apiKey  string
			isValid bool
			body    string
		}{apiKey, isValid, body}
	}
}

func main() {
	flag.Parse()

	if apiKey == "" && filePath == "" {
		fmt.Println("Either -k or -f must be provided.")
		flag.Usage()
		os.Exit(1)
	}

	var validCount, invalidCount int

	if apiKey != "" {
		isValid, body := validateAPIKey(apiKey, endpoint)
		isValid = !isValid
		if isValid || showInvalid {
			printResult(apiKey, !isValid, body)
		}
		if isValid {
			validCount++
		} else {
			invalidCount++
		}
	}

	if filePath != "" {
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		apiKeys := make([]string, 0)
		for scanner.Scan() {
			key := strings.TrimSpace(scanner.Text())
			if key != "" {
				apiKeys = append(apiKeys, key)
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			os.Exit(1)
		}

		apiKeyChan := make(chan string, len(apiKeys))
		results := make(chan struct {
			apiKey  string
			isValid bool
			body    string
		}, len(apiKeys))

		var wg sync.WaitGroup
		for i := 0; i < numThreads; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				worker(apiKeyChan, results)
			}()
		}

		for _, key := range apiKeys {
			apiKeyChan <- key
		}
		close(apiKeyChan)

		wg.Wait()
		close(results)

		for result := range results {
			if result.isValid || showInvalid {
				printResult(result.apiKey, !result.isValid, result.body)
			}
			if result.isValid {
				validCount++
			} else {
				invalidCount++
			}
		}
	}

	fmt.Printf("\nSummary:\n")
	fmt.Printf("Valid API keys: %d\n", validCount)
	fmt.Printf("Invalid API keys: %d\n", invalidCount)
}
