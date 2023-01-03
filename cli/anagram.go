package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/mlange-42/xwrd/anagram"
	"github.com/mlange-42/xwrd/core"
	"github.com/mlange-42/xwrd/util"
	"github.com/spf13/cobra"
)

func anagramCommand(config *core.Config) *cobra.Command {
	var dict string
	var partial bool
	var multi bool
	var maxWords int

	anagram := &cobra.Command{
		Use:   "anagram [WORDS...]",
		Short: "Find anagrams",
		Long: `Find anagrams.

Finds anagrams, incl. partial anagrams as well as combined/multi-anagrams of multiple words.

Enters interactive mode if called without position arguments (i.e. words).
`,
		Aliases: []string{"a"},
		Args:    util.WrappedArgs(cobra.ArbitraryArgs),
		Run: func(cmd *cobra.Command, args []string) {
			if !multi && cmd.Flags().Changed("max-words") {
				fmt.Println("ERROR: flag --max-words is only supported with flag --multi")
				return
			}

			dictionary := config.GetDict()
			if dict != "" {
				dictionary = util.NewDict(dict)
			}
			words, err := util.ReadWordList(dictionary)
			if err != nil {
				fmt.Printf("failed to find anagrams: %s", err.Error())
				return
			}

			tree := anagram.NewTree([]rune(anagram.Letters))
			progress := make(chan int, 8)
			go tree.AddWords(words, progress)

			for pr := range progress {
				bar := strings.Repeat("#", pr/2)
				fmt.Fprintf(os.Stderr, "\rBuilding tree: [%-50s]", bar)
			}
			fmt.Fprintln(os.Stderr)

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
					if partial {
						ana := tree.PartialAnagrams(word)
						for _, res := range ana {
							fmt.Print("  ")
							fmt.Println(strings.Join(res, "  "))
						}
					} else if multi {
						ana := tree.MultiAnagrams(word, maxWords, false)
						for _, res := range ana {
							fmt.Print("  ")
							for b, block := range res {
								fmt.Print(strings.Join(block, "  "))
								if b < len(res)-1 {
									fmt.Print("  |  ")
								}
							}
							fmt.Println()
						}
					} else {
						ana := tree.Anagrams(word)
						if len(ana) > 0 {
							fmt.Printf("  %s\n", strings.Join(ana, "  "))
						}
					}
				}

				if !interactive {
					break
				}
			}
		},
	}
	anagram.Flags().StringVarP(&dict, "dict", "d", "", "Path to the dictionary/word list to use.")

	anagram.Flags().BoolVarP(&partial, "partial", "p", false, "Find partial anagrams.")
	anagram.Flags().BoolVarP(&multi, "multi", "m", false, "Find combinations of multiple partial anagrams.")

	anagram.Flags().IntVarP(&maxWords, "max-words", "w", 0, "Word count limit for multi-anagrams.")

	anagram.MarkFlagsMutuallyExclusive("partial", "multi")

	return anagram
}
