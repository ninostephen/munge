// Package cmd is the root of all thing Munge Go Brrrrr!!!!!
/*
Copyright © 2023 Nino Stephen <ninostephen.me>
*/
package cmd

import (
	"os"

	"github.com/ninostephen/munge/models"
	"github.com/ninostephen/munge/worker"
	"github.com/spf13/cobra"
)

var (
	flagvals     models.Flags
	wordlist     []string
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

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "mgbr",
	Aliases: []string{"mgb", "mungegobrrrrr"},
	Short:   "Dirty little word munger",
	Long: `
 ______                              
|  ___ \                             
| | _ | | _   _  ____    ____   ____ 
| || || || | | ||  _ \  / _  | / _  )
| || || || |_| || | | |( ( | |( (/ / 
|_||_||_| \____||_| |_| \_|| | \____)
                       (_____|       
  ______           ______                                      _  _  _  _  _ 
 / _____)         (____  \                                    | || || || || |
| /  ___   ___     ____)  )  ____   ____   ____   ____   ____ | || || || || |
| | (___) / _ \   |  __  (  / ___) / ___) / ___) / ___) / ___)|_||_||_||_||_|
| \____/|| |_| |  | |__)  )| |    | |    | |    | |    | |     _  _  _  _  _ 
 \_____/  \___/   |______/ |_|    |_|    |_|    |_|    |_|    |_||_||_||_||_|
                              Copyright © 2023 Nino Stephen <ninostephen.me>   
	
	A faster version of Mudge by Th3S3cr3tAg3nt
	
	`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {

		worker.Start(cmd, flagvals)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&flagvals.Word, "word", "w", "", "word to munge")
	rootCmd.Flags().IntVarP(&flagvals.Level, "level", "l", 2, "munge level [1-3] (default 2)")
	rootCmd.Flags().StringVarP(&flagvals.Input, "input", "i", "", "input file")
	rootCmd.Flags().StringVarP(&flagvals.Output, "output", "o", "", "output file")
}

// // mudge command triggers functions based on level which was passed as argument.
// func munge(word string, level int) {
// 	switch level {
// 	case 1:
// 		basic(word)
// 	case 2:
// 		advanced(word)
// 	case 3:
// 		expert(word)
// 	}
// }

// // swapcase function toggles characters from lowercase to uppercase and wiseversa
// func swapcase(word string) string {
// 	swapped := ""
// 	for _, char := range word {
// 		if unicode.IsUpper(char) {
// 			swapped += strings.ToLower(string(char))
// 		} else if unicode.IsLower(char) {
// 			swapped += strings.ToUpper(string(char))
// 		} else {
// 			swapped += string(char)
// 		}
// 	}
// 	return swapped
// }

// // replace function replaces characters in the word with leetSpeak and appends postfix numbers
// func replace(word string, chars map[string]string, nums []string) {
// 	for char, val := range chars {
// 		word = strings.ReplaceAll(word, char, val)
// 		wordlist = append(wordlist, word)
// 		for _, val = range nums {
// 			wordlist = append(wordlist, word+val)

// 		}
// 	}
// }

// // removeDuplicateStr removes all duplicate strings
// func removeDuplicateStr(strSlice []string) []string {
// 	allKeys := make(map[string]bool)
// 	list := []string{}
// 	for _, item := range strSlice {
// 		if _, value := allKeys[item]; !value {
// 			allKeys[item] = true
// 			list = append(list, item)
// 		}
// 	}
// 	return list
// }

// // basic function does all the magic for level 1
// func basic(word string) {
// 	wordlist = append(wordlist, word)
// 	wordlist = append(wordlist, strings.ToUpper(word))
// 	wordlist = append(wordlist, strings.ToTitle(word))
// 	temp := strings.ToTitle(word)
// 	wordlist = append(wordlist, swapcase(temp))
// }

// // advanced function does the calculations for level 2
// func advanced(word string) {
// 	basic(word)
// 	replace(word, leetSpeakMap, level2Postfix)

// }

// // expert function does all the work for level 3
// func expert(word string) {
// 	advanced(word)
// 	replace(word, leetSpeakMap, level3Postfix)
// }

// var scannerPool = sync.Pool{
// 	New: func() interface{} {
// 		return bufio.NewScanner(nil)
// 	},
// }

// func start2(cmd *cobra.Command) {

// 	var wg sync.WaitGroup

// 	if flagvals.input != "" {
// 		file, err := os.Open(flagvals.input)
// 		if err != nil {
// 			fmt.Println("Error opening file:", err)
// 			return
// 		}
// 		defer file.Close()

// 		// Obtain a scanner from the pool
// 		scanner := scannerPool.Get().(*bufio.Scanner)
// 		// scanner.Reset(file)

// 		// Start a goroutine to read the file line by line
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()

// 			for scanner.Scan() {
// 				word := scanner.Text()
// 				munge(strings.TrimSpace(word), flagvals.level)
// 			}

// 			if err := scanner.Err(); err != nil {
// 				fmt.Println("Error reading file:", err)
// 			}

// 			// Put the scanner back into the pool
// 			scannerPool.Put(scanner)
// 		}()
// 	} else {
// 		cmd.Help()
// 		return
// 	}
// 	// Create a channel to receive the processed words
// 	wordChan := make(chan string)

// 	// Start a goroutine to write the words to the output file or print to stdout
// 	if flagvals.output != "" {
// 		file, err := os.Create(flagvals.output)
// 		if err != nil {
// 			fmt.Println("Error creating file:", err)
// 			return
// 		}
// 		defer file.Close()

// 		// Create a writer for the file
// 		writer := bufio.NewWriter(file)

// 		// Start a goroutine to write the words to the file
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()

// 			for word := range wordChan {
// 				_, err := writer.WriteString(word + "\n")
// 				if err != nil {
// 					fmt.Println("Error writing to file:", err)
// 					return
// 				}
// 			}

// 			// Flush the writer to ensure all data is written to the file
// 			err := writer.Flush()
// 			if err != nil {
// 				fmt.Println("Error flushing writer:", err)
// 				return
// 			}
// 		}()
// 	} else {
// 		// Start a goroutine to print the words to stdout
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()

// 			for word := range wordChan {
// 				fmt.Println(word)
// 			}
// 		}()
// 	}

// 	// Process the string if word is provided
// 	if flagvals.word != "" {
// 		munge(
// 			strings.ToLower(flagvals.word),
// 			flagvals.level,
// 		)
// 	}

// 	// Wait for all goroutines to finish
// 	wg.Wait()

// 	fmt.Println("Word list processed successfully.")

// }

// func start(cmd *cobra.Command) {
// 	// sanitize level val
// 	if flagvals.level > 3 {
// 		flagvals.level = 3
// 	} else if flagvals.level < 0 {
// 		flagvals.level = 0
// 	}
// 	// if a word is passed in, just process that else process everything in the file
// 	if flagvals.word != "" {
// 		munge(
// 			strings.ToLower(flagvals.word),
// 			flagvals.level,
// 		)
// 	} else if flagvals.input != "" {
// 		file, err := os.Open(flagvals.input)
// 		if err != nil {
// 			fmt.Println("Error opening file:", err)
// 			return
// 		}
// 		defer file.Close()

// 		// Create a scanner to read the file line by line
// 		scanner := bufio.NewScanner(file)

// 		// Read the file line by line
// 		for scanner.Scan() {
// 			word := scanner.Text()
// 			munge(strings.TrimSpace(word), flagvals.level)
// 		}

// 		// Check for any scanner errors
// 		if err := scanner.Err(); err != nil {
// 			fmt.Println("Error reading file:", err)
// 		}

// 	} else {
// 		cmd.Help()
// 		return
// 	}
// 	sort.SliceStable(wordlist, func(i, j int) bool {
// 		return wordlist[i] < wordlist[j]
// 	})
// 	wordlist = removeDuplicateStr(wordlist)
// 	if flagvals.output != "" {
// 		file, err := os.Create(flagvals.output)
// 		if err != nil {
// 			fmt.Println("Error creating file:", err)
// 			return
// 		}
// 		defer file.Close()

// 		// Create a writer for the file
// 		writer := bufio.NewWriter(file)

// 		// Write each word to the file
// 		for _, word := range wordlist {
// 			_, err := writer.WriteString(word + "\n")
// 			if err != nil {
// 				fmt.Println("Error writing to file:", err)
// 				return
// 			}
// 		}

// 		// Flush the writer to ensure all data is written to the file
// 		err = writer.Flush()
// 		if err != nil {
// 			fmt.Println("Error flushing writer:", err)
// 			return
// 		}

// 		fmt.Println("Word list written to file successfully.")
// 	} else {
// 		fmt.Println(strings.Join(wordlist, "\n"))
// 	}

// }

// func startOld(cmd *cobra.Command) {
// 	if flagvals.level > 3 {
// 		flagvals.level = 3
// 	} else if flagvals.level < 0 {
// 		flagvals.level = 0
// 	}
// 	if flagvals.word != "" {
// 		munge(
// 			strings.ToLower(flagvals.word),
// 			flagvals.level,
// 		)
// 	} else if flagvals.input != "" {
// 		file, err := os.Open(flagvals.input)
// 		if err != nil {
// 			fmt.Println("Error opening file:", err)
// 			return
// 		}
// 		defer file.Close()

// 		// Create a scanner to read the file line by line
// 		scanner := bufio.NewScanner(file)

// 		// Read the file line by line
// 		for scanner.Scan() {
// 			word := scanner.Text()
// 			munge(strings.TrimSpace(word), flagvals.level)
// 		}

// 		// Check for any scanner errors
// 		if err := scanner.Err(); err != nil {
// 			fmt.Println("Error reading file:", err)
// 		}

// 	} else {
// 		cmd.Help()
// 		return
// 	}
// 	sort.SliceStable(wordlist, func(i, j int) bool {
// 		return wordlist[i] < wordlist[j]
// 	})
// 	wordlist = removeDuplicateStr(wordlist)
// 	if flagvals.output != "" {
// 		file, err := os.Create(flagvals.output)
// 		if err != nil {
// 			fmt.Println("Error creating file:", err)
// 			return
// 		}
// 		defer file.Close()

// 		// Create a writer for the file
// 		writer := bufio.NewWriter(file)

// 		// Write each word to the file
// 		for _, word := range wordlist {
// 			_, err := writer.WriteString(word + "\n")
// 			if err != nil {
// 				fmt.Println("Error writing to file:", err)
// 				return
// 			}
// 		}

// 		// Flush the writer to ensure all data is written to the file
// 		err = writer.Flush()
// 		if err != nil {
// 			fmt.Println("Error flushing writer:", err)
// 			return
// 		}

// 		fmt.Println("Word list written to file successfully.")
// 	} else {
// 		fmt.Println(strings.Join(wordlist, "\n"))
// 	}
// }
