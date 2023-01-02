package cli

import (
	"fmt"
	"io/ioutil"
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
		Args:    util.WrappedArgs(cobra.MinimumNArgs(1)),
		Run: func(cmd *cobra.Command, args []string) {
			if !multi && cmd.Flags().Changed("max-words") {
				fmt.Println("ERROR: flag --max-words is only supported with flag --multi")
				return
			}

			text := strings.Join(args, " ")

			tree := anagram.NewTree([]rune(anagram.Letters))
			fileContent, err := ioutil.ReadFile(dict)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			words := strings.Split(string(fileContent), "\n")
			tree.AddWords(words)

			if partial {
				ana := tree.PartialAnagrams(text)
				for _, res := range ana {
					fmt.Println(res)
				}
				return
			}
			if multi {
				ana := tree.MultiAnagrams(text, maxWords, false)
				for _, res := range ana {
					fmt.Println(res)
				}
				return
			}
			fmt.Println(tree.Anagrams(text))
		},
	}
	anagram.Flags().StringVarP(&dict, "dict", "d", "./data/german-700k.txt", "Path to the dictionary/word list to use.")

	anagram.Flags().BoolVarP(&partial, "partial", "p", false, "Find partial anagrams.")
	anagram.Flags().BoolVarP(&multi, "multi", "m", false, "Find combinations of multiple partial anagrams.")

	anagram.Flags().IntVarP(&maxWords, "max-words", "w", 0, "Word count limit for multi-anagrams.")

	anagram.MarkFlagsMutuallyExclusive("partial", "multi")

	return anagram
}
