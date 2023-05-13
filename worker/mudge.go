// Package worker provides all the Mudge functions.
package worker

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
	"unicode"

	"github.com/ninostephen/munge/models"
	"github.com/spf13/cobra"
)

var (
	leetSpeakMap = map[string]string{
		"e": "3",
		"a": "4",
		"i": "!",
		"o": "0",
		"s": "$",
	}
	level2Postfix = []string{
		"1", "123456", "12", "2", "123", "!", ".",
		"?", "_", "0", "01", "69", "21", "22", "23", "1234",
		"8", "9", "10", "11", "13", "3", "4", "5", "6", "7",
	}
	level3Postfix = []string{
		"07", "08", "09", "14", "15", "16", "17", "18", "19", "24", "77",
		"88", "99", "12345", "123456789", "00", "02", "03", "04", "05", "06",
		"19", "20", "25", "26", "27", "28", "007", "1234567", "12345678", "111111",
		"111", "777", "666", "101", "33", "44", "55", "66", "2008", "2009", "2010",
		"2011", "86", "87", "89", "90", "91", "92", "93", "94", "95", "98",
	}
)

// mudge mutates the words based on level, sorts and removes duplicates before returning
// a slice of mutated words
func munge(word string, level int) []string {
	var mutatedWord []string
	switch level {
	case 1:
		mutatedWord = basic(word)
	case 2:
		mutatedWord = advanced(word)
	case 3:
		mutatedWord = expert(word)
	}

	// sort the mutated words and remove duplicates
	sort.SliceStable(mutatedWord, func(i, j int) bool {
		return mutatedWord[i] < mutatedWord[j]
	})
	mutatedWord = removeDuplicateStr(mutatedWord)

	return mutatedWord
}

// swapcase function toggles characters from lowercase to uppercase and vise-versa
func swapcase(word string) string {
	swapped := ""
	for _, char := range word {
		if unicode.IsUpper(char) {
			swapped += strings.ToLower(string(char))
		} else if unicode.IsLower(char) {
			swapped += strings.ToUpper(string(char))
		} else {
			swapped += string(char)
		}
	}
	return swapped
}

// replace mutates strings to leetspeak and appends commonly used numbers to it before returning a
// slice of mutated words
func replace(word string, chars map[string]string, nums []string) []string {
	var wordlist []string
	for char, val := range chars {
		word = strings.ReplaceAll(word, char, val)
		wordlist = append(wordlist, word)
		for _, val = range nums {
			wordlist = append(wordlist, word+val)

		}
	}
	return wordlist
}

// removeDuplicateStr removes all duplicate strings
func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

// basic function does basic level mutiations and returns a wordlist
func basic(word string) []string {
	var wordlist []string

	wordlist = append(wordlist, word)
	wordlist = append(wordlist, strings.ToUpper(word))

	wordToTitleCase := strings.ToTitle(word)
	wordlist = append(wordlist, wordToTitleCase)
	wordlist = append(wordlist, swapcase(wordToTitleCase))

	return wordlist
}

// advanced function does level 2 mutations
func advanced(word string) []string {
	var wordlist []string

	wordlist = append(wordlist, basic(word)...)
	wordlist = append(wordlist, replace(word, leetSpeakMap, level2Postfix)...)

	return wordlist
}

// expert function does level 3 mutations
func expert(word string) []string {
	var wordlist []string

	wordlist = append(wordlist, advanced(word)...)
	wordlist = append(wordlist, replace(word, leetSpeakMap, level3Postfix)...)

	return wordlist
}

// Start function does all the mudging
func Start(cmd *cobra.Command, flagvals models.Flags) {

	var wg sync.WaitGroup
	taskQueue := make(chan string)
	completedQueue := make(chan string)

	if flagvals.Level > 3 {
		flagvals.Level = 3
	} else if flagvals.Level < 0 {
		flagvals.Level = 0
	}

	if flagvals.Word != "" {
		wg.Add(1)
		go addToQueueSingular(flagvals.Word, flagvals.Level, &wg, &completedQueue)

	} else if _, err := os.Stat(flagvals.Input); err == nil {

		println("Initiating input file read: ", flagvals.Input)
		inputFile, err := os.Open(flagvals.Input)
		if err != nil {
			fmt.Printf("Failed to open input file: %v \n", err)
			return
		}

		wg.Add(1)
		go parseFile(inputFile, &wg, &taskQueue)

		maxGoroutines := runtime.NumCPU()
		println("Identified Max number of Go routines: ", maxGoroutines)

		for agentID := 0; agentID < maxGoroutines; agentID++ {
			wg.Add(1)
			go genWorkers(agentID, flagvals.Level, &wg, &taskQueue, &completedQueue)
		}

	} else {
		close(completedQueue)
		close(taskQueue)
		cmd.Help()
		return
	}

	if flagvals.Output != "" {

		outputFile, err := os.Create(flagvals.Output)
		if err != nil {
			fmt.Printf("Error creating output file: %v \n", err)
			return
		}
		defer outputFile.Close()
		wg.Add(1)
		go writeFile(outputFile, &wg, &completedQueue)
	} else {
		wg.Add(1)
		go printFromQueue(&completedQueue, &wg)
	}

	wg.Wait()
}

func parseFile(inputFile *os.File, wg *sync.WaitGroup, taskQueue *chan string) {
	defer inputFile.Close()
	reader := bufio.NewReader(inputFile)
	defer wg.Done()
	for {
		word, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading from file:", err)
			}
			break
		}
		println("Pushed: ", word)
		*taskQueue <- word
	}

}
func writeFile(outputFile *os.File, wg *sync.WaitGroup, completedQueue *chan string) {
	defer wg.Done()

	for finalWord := range *completedQueue {
		println("Wrote: ", finalWord)
		_, err := outputFile.WriteString(finalWord + "\n")
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}
	println("completed!")
	close(*completedQueue)
}

func printFromQueue(completedQueue *chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(*completedQueue)

	for finalWord := range *completedQueue {
		fmt.Println("Mutated word:", finalWord)
	}

}

func addToQueueSingular(word string, level int, wg *sync.WaitGroup, completedQueue *chan string) {
	var wordlist []string
	wordlist = munge(word, level)
	defer wg.Done()

	for _, finalWord := range wordlist {
		*completedQueue <- finalWord
	}
	println("Completed adding words to completed queue")
}

func genWorkers(agentID int, level int, wg *sync.WaitGroup, taskQueue, completedQueue *chan string) {
	defer wg.Done()
	println("Agent launched: ", agentID)

	for word := range *taskQueue {
		mutatedList := munge(word, level)
		for _, mutation := range mutatedList {
			if mutation != "" {
				fmt.Printf("Agent %d Pushed %v", agentID, mutation)
				*completedQueue <- mutation
			}
		}
	}
	close(*taskQueue)
}
