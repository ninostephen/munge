// Package worker provides the core functionality for the Munge application.
package worker

//
// This package includes functions for mutating words based on different levels,
// sorting and removing duplicates from the mutated words, and performing various
// operations related to word manipulation.
//

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"unicode"

	"github.com/ninostephen/munge/models"
	"github.com/spf13/cobra"
)

var (
	activeAgents int32
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

// munge mutates the words based on the specified level, sorts them, and removes duplicates.
// It returns a slice of mutated words.
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

	// Sort the mutated words and remove duplicates
	sort.SliceStable(mutatedWord, func(i, j int) bool {
		return mutatedWord[i] < mutatedWord[j]
	})
	mutatedWord = removeDuplicateStr(mutatedWord)

	return mutatedWord
}

// swapcase toggles characters from lowercase to uppercase and vice versa in a word.
// It returns the word with swapped case.
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

// replace mutates the word by replacing characters in it with leetspeak characters
// and appends commonly used numbers to it. It returns a slice of mutated words.
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

// removeDuplicateStr removes all duplicate strings from a string slice.
// It returns a new slice with unique strings.
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

// basic performs basic level mutations on a word and returns a wordlist.
func basic(word string) []string {
	var wordlist []string

	wordlist = append(wordlist, word)
	wordlist = append(wordlist, strings.ToUpper(word))

	wordToTitleCase := strings.ToTitle(word)
	wordlist = append(wordlist, wordToTitleCase)
	wordlist = append(wordlist, swapcase(wordToTitleCase))

	return wordlist
}

// advanced performs advanced level mutations on a word and returns a wordlist.
func advanced(word string) []string {
	var wordlist []string

	wordlist = append(wordlist, basic(word)...)
	wordlist = append(wordlist, replace(word, leetSpeakMap, level2Postfix)...)

	return wordlist
}

// expert performs expert level mutations on a word and returns a wordlist.
func expert(word string) []string {
	var wordlist []string

	wordlist = append(wordlist, advanced(word)...)
	wordlist = append(wordlist, replace(word, leetSpeakMap, level3Postfix)...)

	return wordlist
}

// Start begins the Munge process with the specified flags.
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
		go addWordToQueue(flagvals.Word, flagvals.Level, &wg, &completedQueue)

	} else if _, err := os.Stat(flagvals.Input); err == nil {
		inputFile, err := os.Open(flagvals.Input)
		if err != nil {
			fmt.Printf("Failed to open input file: %v \n", err)
			return
		}

		wg.Add(1)
		go parseFile(inputFile, &wg, &taskQueue)

		maxGoroutines := runtime.NumCPU()
		atomic.AddInt32(&activeAgents, int32(maxGoroutines))

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

// parseFile reads words from the input file and adds them to the task queue.
func parseFile(inputFile *os.File, wg *sync.WaitGroup, taskQueue *chan string) {
	defer inputFile.Close()
	reader := bufio.NewReader(inputFile)
	for {
		word, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading from file:", err)
			}
			break
		}
		*taskQueue <- word
	}
	close(*taskQueue)
	wg.Done()
}

// writeFile writes the final mutated words to the output file.
func writeFile(outputFile *os.File, wg *sync.WaitGroup, completedQueue *chan string) {

	for finalWord := range *completedQueue {
		if finalWord == "<EOL>" {
			close(*completedQueue)
			break
		}
		if finalWord == "" {
			continue
		}
		_, err := outputFile.WriteString(finalWord)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}
	wg.Done()
}

// printFromQueue prints the final mutated words from the completed queue.
func printFromQueue(completedQueue *chan string, wg *sync.WaitGroup) {
	for finalWord := range *completedQueue {
		if finalWord == "<EOL>" {
			close(*completedQueue)
			break
		}
		fmt.Printf(finalWord)
	}
	wg.Done()
}

// addWordToQueue adds a single word to the task queue for mutation.
func addWordToQueue(word string, level int, wg *sync.WaitGroup, completedQueue *chan string) {
	var wordlist []string
	wordlist = munge(word, level)
	fmt.Println(wordlist)
	for _, finalWord := range wordlist {
		*completedQueue <- finalWord
	}
	*completedQueue <- "<EOL>"
	wg.Done()

}

// genWorkers generates worker goroutines that mutate words from the task queue.
func genWorkers(agentID int, level int, wg *sync.WaitGroup, taskQueue, completedQueue *chan string) {

	for word := range *taskQueue {
		mutatedList := munge(word, level)
		for _, mutation := range mutatedList {
			if mutation != "" {
				// fmt.Printf("Agent %d Pushed %v", agentID, mutation)
				*completedQueue <- mutation
			}
		}
	}

	atomic.AddInt32(&activeAgents, -1)

	if atomic.LoadInt32(&activeAgents) == 0 {
		*completedQueue <- "<EOL>"
	}
	wg.Done()
}
