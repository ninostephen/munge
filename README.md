# mudge: Golang port of Munge by Th3S3cr3tAg3nt

This program is a word munging tool that mutates words based on different levels of complexity. It provides functionalities to transform words using basic, advanced, and expert level mutations. The mutated words are then sorted and duplicates are removed, resulting in a list of unique and sorted mutated words.

## Features:

- Basic Level Mutations: Swapping case (lowercase to uppercase and vice versa) and title case variations.
- Advanced Level Mutations: Leetspeak transformations and appending commonly used numbers.
- Expert Level Mutations: Combining advanced level mutations with additional numbers and postfixes.
- Input Options: Mutate a single word or read words from an input file.
- Output Options: Write mutated words to an output file or print them to the console.

## Examples:

Mutate a Single Word:
```bash
$ munge -word "example" -level 2 -output "output.txt"
```

Mutate Words from an Input File:
```bash
$ munge -input "input.txt" -level 3 -output "output.txt"
```
Print Mutated Words to Console:
```bash
$ munge -input "input.txt" -level 1
```

## Installation

To install the Word Munging Tool, follow these steps:

- Make sure you have Go installed on your system. You can download and install Go from the official website: https://golang.org/

- Clone the repository to your local machine:

```bash
git clone https://github.com/ninostephen/munge
```

Navigate to the project directory:
```bash
cd munge
```
Install the required packages:
```bash
go mod tidy
```
Build the code:
```bash
go build -o munge main.go
```

Remember to add path to the executable to your path variable. 

## Usage:
```
munge [flags]

Flags:
-h, --help Show help information
-i, --input string Specify the input file
-l, --level int Set the munge level [1-3] (default 2)
-o, --output string Specify the output file
-w, --word string Specify the word to munge
```

## Detailed Description of Flags:
Munge supports the following options:

    Input File (-i, --input)
Specify the path to an input file containing words to be mutated.

    Munge Level (-l, --level)
Set the level of mutations to apply to the words.
Valid values: 1, 2, 3 (default: 2).

    Output File (-o, --output)
Specify the path to an output file to store the mutated words.

    Word to Munge (-w, --word)
Specify a single word to be mutated.

Use the appropriate flags to customize the behavior of the word munging tool according to your requirements.

Note: Make sure to provide the necessary input and output files and ensure the required packages are imported for the program to function correctly.
