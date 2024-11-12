package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/spf13/cobra"
)

type Counters struct {
	Lines int64
	Words int64
	Bytes int64
	Chars int64
}

type Flags struct {
	CountBytes  bool
	CountLines  bool
	CountWords  bool
	CountChars  bool
	FileMode    bool
	CompareMode bool
}

var state = &Flags{
	CountBytes: false,
	CountLines: false,
	CountWords: false,
	CountChars: false,
	FileMode:   false,
}

func (f *Flags) AnySet() bool {
	return f.CountBytes || f.CountLines || f.CountWords || f.CountChars
}

func (f *Flags) SetDefaults() {
	f.CountBytes = true
	f.CountLines = true
	f.CountWords = true
}

func printOutput(arg string, c *Counters) {
	var output []string

	if state.CountLines {
		output = append(output, fmt.Sprintf("%d", c.Lines))
	}

	if state.CountWords {
		output = append(output, fmt.Sprintf("%d", c.Words))
	}

	if state.CountBytes {
		output = append(output, fmt.Sprintf("%d", c.Bytes))
	}

	if state.CountChars {
		output = append(output, fmt.Sprintf("%d", c.Chars))
	}

	output = append(output, arg)

	fmt.Printf("\t%s\n", strings.Join(output, " "))
}

func getReaderFromFilePath(filePath string) (*bufio.Reader, func() error, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	return bufio.NewReader(file), file.Close, nil
}

func getReaderFromStdIn() (*bufio.Reader, func() error, error) {
	return bufio.NewReader(os.Stdin), func() error { return nil }, nil
}

func generalRuneCounter(reader *bufio.Reader) (*Counters, error) {
	inWord := false
	counters := &Counters{}

	for {
		r, n, err := reader.ReadRune()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		counters.Bytes += int64(n)
		counters.Chars++

		// Count lines by checking for newline characters
		if r == '\n' {
			counters.Lines++
		}

		// Count words by identifying word boundaries
		if unicode.IsSpace(r) {
			inWord = false
		} else if !inWord {
			inWord = true
			counters.Words++
		}

	}

	return counters, nil
}

func generalBufferCounter(reader *bufio.Reader) (*Counters, error) {
	buf := make([]byte, 8192)
	inWord := false
	counters := &Counters{}

	for {
		n, err := reader.Read(buf)

		if n > 0 {
			counters.Bytes += int64(n)
			chunk := buf[:n]

			i := 0
			for i < len(chunk) {
				r, size := utf8.DecodeRune(chunk[i:])
				if r == utf8.RuneError && size == 1 {
					return counters, fmt.Errorf("invalid UTF-8")
				}

				counters.Chars++

				if r == '\n' {
					counters.Lines++
				}

				// Word counting
				if unicode.IsSpace(r) {
					inWord = false
				} else if !inWord {
					inWord = true
					counters.Words++
				}

				i += size // Move to the next rune
			}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return counters, err
		}
	}

	return counters, nil
}

func wc(args []string) {
	if len(args) != 0 {
		state.FileMode = true
	}

	if !state.AnySet() {
		state.SetDefaults()
	}

	if state.FileMode {
		for _, arg := range args {
			reader, closer, err := getReaderFromFilePath(arg)
			if err != nil {
				log.Fatal(err)
			}
			process(arg, reader, closer)
		}

		if state.CompareMode {
			for _, arg := range args {
				reader, closer, err := getReaderFromFilePath(arg)
				if err != nil {
					log.Fatal(err)
				}
				processBuffer(arg, reader, closer)
			}
		}

	} else {
		reader, closer, err := getReaderFromStdIn()
		if err != nil {
			log.Fatal(err)
		}
		process("", reader, closer)
	}
}

func process(arg string, reader *bufio.Reader, closer func() error) {
	defer func() { _ = closer() }()

	var counters *Counters
	var err error

	timeElapsed(func() error {
		counters, err = generalRuneCounter(reader)
		return err
	})

	if err != nil {
		log.Fatal(err)
	}

	printOutput(arg, counters)
}

func processBuffer(arg string, reader *bufio.Reader, closer func() error) {
	defer func() { _ = closer() }()

	var counters *Counters
	var err error

	timeElapsed(func() error {
		counters, err = generalBufferCounter(reader)
		return err
	})

	if err != nil {
		log.Fatal(err)
	}
	printOutput(arg, counters)
}

var rootCmd = &cobra.Command{
	Use:   "wc",
	Short: "wc is a command line tool for counting lines, words and bytes of text files",
	Long:  `wc is a command line tool for counting lines, words and bytes of text files`,
	Run: func(cmd *cobra.Command, args []string) {
		wc(args)
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&state.CountBytes, "count bytes", "c", false, "count bytes")
	rootCmd.PersistentFlags().BoolVarP(&state.CountLines, "count lines", "l", false, "count lines")
	rootCmd.PersistentFlags().BoolVarP(&state.CountWords, "count words", "w", false, "count words")
	rootCmd.PersistentFlags().BoolVarP(&state.CountChars, "count chars", "m", false, "count chars")
	rootCmd.PersistentFlags().BoolVarP(&state.CompareMode, "compare mode", "s", false, "compare mode")
}

func timeElapsed(f func() error) {
	currenTime := time.Now()
	_ = f()
	elapsedTime := time.Since(currenTime)
	fmt.Printf("Elapsed time: %s\n", elapsedTime)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
