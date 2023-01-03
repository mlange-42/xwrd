package cli

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"unicode/utf8"

	"github.com/mlange-42/xwrd/util"
	"github.com/spf13/cobra"
)

func matchCommand() *cobra.Command {
	var dict string

	match := &cobra.Command{
		Use:   "match [WORDS...]",
		Short: "Find words matching a pattern of letter positions",
		Long: `Find words matching a pattern of letter positions.

Enters interactive mode if called without position arguments (i.e. words).

Patterns
--------

'.' (period) strands for one arbitrary letter
'*' (asterisk) stands for 0 or more arbitrary letters

Examples
--------

a....  - find all 5-letter words starting with 'a'
*pf    - find all words ending with 'pf'
a....b - find all words of length 6 stat start with 'a' and end with 'b'
`,
		Aliases: []string{"m"},
		Args:    util.WrappedArgs(cobra.ArbitraryArgs),
		Run: func(cmd *cobra.Command, args []string) {
			if dict == "" {
				dict = util.DictPath()
			}
			words, err := util.ReadWordList(dict)
			if err != nil {
				fmt.Printf("failed to find matching words: %s", err.Error())
				return
			}

			interactive := len(args) == 0

			for {
				var text []string
				if interactive {
					fmt.Print("Enter a word: ")
					var answer string
					scanner := bufio.NewScanner(os.Stdin)
					if scanner.Scan() {
						answer = scanner.Text()
					}
					if len(answer) == 0 {
						break
					}
					text = []string{answer}
				} else {
					text = args
				}

				for _, word := range text {
					if !interactive {
						fmt.Printf("%s:\n", word)
					}

					pattern, err := createPattern(word)
					if err != nil {
						fmt.Printf("failed to find matching words: %s", err.Error())
						return
					}
					res := findWords(words, pattern)
					for _, r := range res {
						fmt.Println("  " + r)
					}
				}

				if !interactive {
					break
				}
			}
		},
	}
	match.Flags().StringVarP(&dict, "dict", "d", "", "Path to the dictionary/word list to use.")

	return match
}

func createPattern(word string) (*regexp.Regexp, error) {
	exp1, err := regexp.Compile("\\.+")
	if err != nil {
		return nil, err
	}
	exp2, err := regexp.Compile("\\*+")
	if err != nil {
		return nil, err
	}
	pattern := exp1.ReplaceAllStringFunc(
		word,
		func(m string) string {
			return fmt.Sprintf("\\p{L}{%d}", utf8.RuneCountInString(m))
		},
	)
	pattern = exp2.ReplaceAllStringFunc(
		pattern,
		func(m string) string {
			return "\\p{L}*"
		},
	)
	pattern = fmt.Sprintf("(?i)^%s$", pattern)
	return regexp.Compile(pattern)
}

func findWords(words []string, pattern *regexp.Regexp) []string {
	results := []string{}

	for _, word := range words {
		if pattern.MatchString(word) {
			results = append(results, word)
		}
	}

	return results
}
