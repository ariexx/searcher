package main

import (
	"bufio"
	"fmt"
	"github.com/go-playground/validator/v10"
	"os"
	"path/filepath"
	"strings"
)

type DomainCount struct {
	Domain string
	Count  int
}

func main() {
	files, err := filepath.Glob("*.txt")
	if err != nil {
		fmt.Println("Error reading files:", err)
		return
	}

	if len(files) == 0 {
		fmt.Println("No .txt files found in the current directory.")
		return
	}

	fmt.Println("Available files:")
	for i, filePath := range files {
		fmt.Printf("%d. %s\n", i+1, filePath)
	}
	fmt.Print("Enter the number of the file you want to process: ")
	var selectedFileIndex int
	fmt.Scanln(&selectedFileIndex)
	if selectedFileIndex < 1 || selectedFileIndex > len(files) {
		fmt.Println("Invalid file selection. Exiting.")
		return
	}
	selectedFilePath := files[selectedFileIndex-1]

	fmt.Print("Enter the domain name you want to search for: ")
	var searchDomain string
	fmt.Scanln(&searchDomain)

	domainFound, countDomain := checkDomainInFile(selectedFilePath, searchDomain)

	if domainFound {

		fmt.Printf("The domain %s appears %d times in the file.\n", countDomain.Domain, countDomain.Count)

		searchResults := searchDomainInFile(selectedFilePath, searchDomain)

		fmt.Print("Enter the number of search results you want to save: ")
		var searchResultsToSave int
		fmt.Scanln(&searchResultsToSave)

		if searchResultsToSave > len(searchResults) {
			fmt.Println("Number of search results to save is greater than the search results.")
			return
		}

		searchResults = searchResults[:searchResultsToSave]
		saveSearchResults(searchResults, searchDomain)
		fmt.Println("Search results saved to result-" + searchDomain + ".txt")
	} else {
		fmt.Println("Domain not found. No results saved.")
	}
}

func extractDomainCounts(fileContent string) []DomainCount {
	scanner := bufio.NewScanner(strings.NewReader(fileContent))
	scanner.Split(bufio.ScanWords)

	domainCountMap := make(map[string]int)

	for scanner.Scan() {
		text := scanner.Text()
		if isValidDomain(text) {
			domainCountMap[text]++
		}
	}

	domainCounts := make([]DomainCount, 0, len(domainCountMap))
	for domain, count := range domainCountMap {
		domainCounts = append(domainCounts, DomainCount{Domain: domain, Count: count})
	}

	return domainCounts
}

func isValidDomain(text string) bool {
	v := validator.New()
	if err := v.Var(text, "required,url"); err != nil {
		return false
	}

	return true
}

func searchDomainInFile(filePath string, searchDomain string) []string {
	searchResults := make([]string, 0)

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", filePath, err)
		return searchResults
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, searchDomain) {
			searchResults = append(searchResults, line)
		}
	}

	return searchResults
}

func saveSearchResults(searchResults []string, searchDomain string) {
	resultFolderPath := "result"
	if _, err := os.Stat(resultFolderPath); os.IsNotExist(err) {
	}

	resultFilePath := filepath.Join(resultFolderPath, "result-"+searchDomain+".txt")

	file, err := os.Create(resultFilePath)
	if err != nil {
		fmt.Printf("Error creating result file: %v\n", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range searchResults {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			fmt.Printf("Error writing to result file: %v\n", err)
			return
		}
	}
}

func checkDomainInFile(filePath, domainName string) (result bool, count DomainCount) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return false, DomainCount{}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	domainExists := false
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), domainName) {
			domainExists = true
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return false, DomainCount{}
	}

	if domainExists {
		fmt.Println("Domain exists in the file.")
	} else {
		fmt.Println("Domain does not exist in the file.")
	}

	count = DomainCount{Domain: domainName, Count: 1}

	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), domainName) {
			count.Count++
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			return false, DomainCount{}
		}
	}

	return domainExists, count
}
