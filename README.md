
# `ccwc` Command Line Tool

## Overview

`wc` is a command-line tool written in Go for counting **lines**, **words**, **characters**, and **bytes** in text files. It replicates and extends the functionality of the Unix `wc` utility, offering enhanced features like comparison mode and custom flag options.

## Features

- **Line Counting**: Counts the number of lines in a file.
- **Word Counting**: Counts the number of words in a file.
- **Byte Counting**: Counts the number of bytes in a file.
- **Character Counting**: Counts the number of characters in a file.
- **Compare Mode**: Provides an additional pass to compare different counting methods.
- **File or Standard Input Support**: Accepts both file paths and standard input.

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/omarahr/ccwc.git
   cd ccwc
   ```
2. Build the binary:
   ```bash
   go build -o ccwc main.go
   ```

## Usage

The tool provides several flags to control the counting functionality:

```bash
wc [flags] [file...]
```

### Flags

- `-c, --count bytes`    : Count bytes.
- `-l, --count lines`    : Count lines.
- `-w, --count words`    : Count words.
- `-m, --count chars`    : Count characters.
- `-s, --compare mode`   : Enable comparison mode for additional validation.

### Examples

1. Count lines, words, and bytes in a file:
   ```bash
   ./ccwc -l -w -c example.txt
   ```

2. Use standard input:
   ```bash
   echo "Hello world" | ./ccwc -w
   ```

3. Compare counting methods for a file:
   ```bash
   ./ccwc -s example.txt
   ```

## License

This project is licensed under the MIT License.

---

## Acknowledgments

- Inspired the solution from [Coding Challenges by John Crickett](https://codingchallenges.fyi/challenges/challenge-wc)