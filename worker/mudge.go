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

	"github.com/spf13/cobra"
)

type flags struct {
	word   string
	input  string
	output string
	level  int
}

var (
	flagvals     flags
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
func Start(cmd *cobra.Command) {
	var wordlist []string
	// Create a task queue channel to hold the mutation tasks
	taskQueue := make(chan string)

	// Create a completed queue channel to hold the results of muation tasks
	completedQueue := make(chan string)

	if flagvals.level > 3 {
		flagvals.level = 3
	} else if flagvals.level < 0 {
		flagvals.level = 0
	}

	if flagvals.word != "" {
		// if wa word was passed in, just munge that
		wordlist = munge(strings.ToLower(flagvals.word), flagvals.level)
		for _, finalWord := range wordlist {
			completedQueue <- finalWord
		}

	} else if _, err := os.Stat(flagvals.input); err == nil {
		// Now that we confirmed that the input file exist, we can now
		// open the input file and start reading
		inputFile, err := os.Open(flagvals.input)
		if err != nil {
			fmt.Printf("Failed to open input file: %v \n", err)
			return
		}
		defer inputFile.Close()

		// Create a WaitGroup to wait for all Goroutines to finish
		var wg sync.WaitGroup

		// Create a reader to read the file line by line
		reader := bufio.NewReader(inputFile)

		// Read the file and distribute mutation tasks among Goroutines
		for {
			// Read a line from the file
			word, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					fmt.Println("Error reading from file:", err)
				}
				break // Exit the loop when done reading
			}

			// Push the word into the task queue for processing by a Goroutine
			taskQueue <- word
		}

		// Close the task queue to signal that no more tasks will be added
		close(taskQueue)

		// Wait for all Goroutines to finish before exiting the program
		wg.Wait()

		// Determine the maximum number of Goroutines based on the available CPU cores
		maxGoroutines := runtime.NumCPU()

		// Create the worker pool
		for i := 0; i < maxGoroutines; i++ {
			// Increment the WaitGroup counter for each Goroutine
			wg.Add(1)

			// Launch a Goroutine as a worker
			go func() {
				defer wg.Done() // Decrement the WaitGroup counter when the Goroutine completes

				// Process tasks from the task queue
				for word := range taskQueue {
					mutatedList := munge(word, flagvals.level) // Call the mutation function on the word
					for _, mutation := range mutatedList {
						completedQueue <- mutation
					}
				}
			}()
		}

	} else {
		cmd.Help()
		return
	}

	if flagvals.output != "" {
		// create an output outputFile or just print the data using an if else.
		outputFile, err := os.Create(flagvals.output)
		if err != nil {
			fmt.Printf("Error creating output file: %v \n", err)
			return
		}
		defer outputFile.Close()

		go func() {
			for finalWord := range completedQueue {
				_, err := outputFile.WriteString(finalWord + "\n")
				if err != nil {
					fmt.Println("Error writing to file:", err)
					return
				}
			}
		}()

		fmt.Println("Word list written to file successfully.")
		close(completedQueue)
	} else {
		// Read the mutated words from the mutations channel
		for finalWord := range completedQueue {
			fmt.Println("Mutated word:", finalWord)
		}
		close(completedQueue)
	}

}
