package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/mlange-42/xwrd/anagram"
	"github.com/mlange-42/xwrd/util"
	"github.com/spf13/cobra"
)

func anagramCommand() *cobra.Command {
	var dict string
	var partial bool
	var multi bool
	var maxWords int

	anagram := &cobra.Command{
		Use:     "anagram [WORDS...]",
		Short:   "Determines anagrams",
		Aliases: []string{"a"},
		Args:    util.WrappedArgs(cobra.ArbitraryArgs),
		Run: func(cmd *cobra.Command, args []string) {
			if !multi && cmd.Flags().Changed("max-words") {
				fmt.Println("ERROR: flag --max-words is only supported with flag --multi")
				return
			}

			tree := anagram.NewTree([]rune(anagram.Letters))
			fileContent, err := ioutil.ReadFile(dict)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			words := strings.Split(string(fileContent), "\n")

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
					fmt.Scanln(&answer)
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
	anagram.Flags().StringVarP(&dict, "dict", "d", "./data/german-700k.txt", "Path to the dictionary/word list to use.")

	anagram.Flags().BoolVarP(&partial, "partial", "p", false, "Find partial anagrams.")
	anagram.Flags().BoolVarP(&multi, "multi", "m", false, "Find combinations of multiple partial anagrams.")

	anagram.Flags().IntVarP(&maxWords, "max-words", "w", 0, "Word count limit for multi-anagrams.")

	anagram.MarkFlagsMutuallyExclusive("partial", "multi")

	return anagram
}
