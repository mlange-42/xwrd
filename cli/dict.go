package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"github.com/mlange-42/xwrd/anagram"
	"github.com/mlange-42/xwrd/core"
	"github.com/mlange-42/xwrd/util"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
)

func dictCommand(config *core.Config) *cobra.Command {
	root := &cobra.Command{
		Use:     "dict",
		Short:   "Handle dictionaries",
		Aliases: []string{"d"},
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	root.AddCommand(showDictsCommand(config))
	root.AddCommand(setDictCommand(config))
	root.AddCommand(listDictsCommand(config))
	root.AddCommand(installDictCommand(config))
	root.AddCommand(analyzeDictCommand(config))

	return root
}

func showDictsCommand(config *core.Config) *cobra.Command {
	download := &cobra.Command{
		Use:     "info",
		Short:   "Shows the currently set dictionary",
		Aliases: []string{"i"},
		Args:    util.WrappedArgs(cobra.NoArgs),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf(config.Dict)
		},
	}
	return download
}

func listDictsCommand(config *core.Config) *cobra.Command {
	download := &cobra.Command{
		Use:     "list",
		Short:   "List installable and installed dictionaries",
		Aliases: []string{"l"},
		Args:    util.WrappedArgs(cobra.NoArgs),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Available:")

			for lang, dicts := range util.Dictionaries {
				for _, d := range dicts {
					fmt.Printf("  %s/%s\n", lang, d.Name)
				}
			}

			allDicts, err := util.AllDictionaries()
			if err != nil {
				fmt.Printf("failed to list dictionaries: %s", err.Error())
				return
			}

			keys := maps.Keys(allDicts)
			sort.Strings(keys)
			fmt.Println("Installed:")
			for _, key := range keys {
				dict := allDicts[key]
				fmt.Printf("  %s/%s\n", dict.Language, strings.TrimSuffix(dict.Name, filepath.Ext(dict.Name)))
			}
			if len(allDicts) == 0 {
				fmt.Printf("  None\n")
			}

		},
	}
	return download
}

func setDictCommand(config *core.Config) *cobra.Command {
	download := &cobra.Command{
		Use:     "set",
		Short:   "Set the default dictionary",
		Aliases: []string{"s"},
		Args:    util.WrappedArgs(cobra.ExactArgs(1)),
		Run: func(cmd *cobra.Command, args []string) {
			dictionary := util.NewDict(args[0])
			if !util.HasDictionary(dictionary) {
				fmt.Printf("failed to set dictionary: dictionary %s is not installed.\nTry: xwrd dict install %[1]s", dictionary.FullName())
				return
			}

			config.Dict = dictionary.FullName()
			err := core.SaveConfig(*config)
			if err != nil {
				fmt.Printf("failed to set dictionary: %s", err.Error())
				return
			}

			fmt.Printf("dictionary set to %s", dictionary.FullName())
		},
	}
	return download
}

func installDictCommand(config *core.Config) *cobra.Command {
	install := &cobra.Command{
		Use:     "install DICT",
		Short:   "Install dictionaries",
		Aliases: []string{"i"},
		Args:    util.WrappedArgs(cobra.ExactArgs(1)),
		Run: func(cmd *cobra.Command, args []string) {
			dictionary := util.NewDict(args[0])
			if util.HasDictionary(dictionary) {
				fmt.Printf("failed to install dictionary: dictionary already exists")
				return
			}

			lang, ok := util.Dictionaries[dictionary.Language]
			if !ok {
				fmt.Printf("failed to install dictionary: no dictionaries in language '%s'", dictionary.Language)
				return
			}

			found := false
			for _, dict := range lang {
				if dict.Name == dictionary.Name {
					dictionary = dict
					found = true
					break
				}
			}
			if !found {
				fmt.Printf("failed to install dictionary: dictionary '%s' not found in language '%s'", dictionary.Name, dictionary.Language)
				return
			}

			fmt.Printf("installing dictionary %s/%s from %s...\n", dictionary.Language, dictionary.Name, dictionary.URL)
			err := util.DownloadDictionary(dictionary)
			if err != nil {
				fmt.Printf("failed to install dictionary: %s", err.Error())
				return
			}
		},
	}
	return install
}

func analyzeDictCommand(config *core.Config) *cobra.Command {
	analyze := &cobra.Command{
		Use:     "analyze [DICT]",
		Short:   "Analyze dictionaries",
		Aliases: []string{"a"},
		Args:    util.WrappedArgs(cobra.MaximumNArgs(1)),
		Run: func(cmd *cobra.Command, args []string) {
			dictName := config.Dict
			if len(args) > 0 {
				dictName = args[0]
			}
			dictionary := util.NewDict(dictName)

			if !util.HasDictionary(dictionary) {
				fmt.Printf("failed to analyze dictionary: dictionary '%s' does not exist", dictName)
				return
			}

			words, err := util.LoadDictionary(dictionary)
			if err != nil {
				fmt.Printf("failed to analyze dictionary: %s", err.Error())
				return
			}

			analyze(words)
		},
	}
	return analyze
}

func analyze(words []string) {
	tree := anagram.NewTree([]rune(anagram.Letters))
	progress := make(chan int, 8)
	go tree.AddWords(words, progress)

	for pr := range progress {
		bar := strings.Repeat("#", pr/2)
		fmt.Fprintf(os.Stderr, "\rBuilding tree: [%-50s]", bar)
	}
	fmt.Fprintln(os.Stderr)

	numWords := len(words)
	numNonAnagrams := len(tree.Leaves)
	lengthHist := []int{}
	totalRunes := map[int]int{}
	maxRunes := map[int]int{}
	wordsWithRune := map[int]int{}
	totalRuneCount := 0

	maxWordLength := 0
	longestWords := []string{}

	for _, word := range words {
		runes := []rune(word)
		if len(runes) == 0 {
			continue
		}

		l := len(runes)
		for len(lengthHist) <= l {
			lengthHist = append(lengthHist, 0)
		}
		lengthHist[l]++
		totalRuneCount += l

		if l > maxWordLength {
			longestWords = longestWords[:0]
			longestWords = append(longestWords, word)
			maxWordLength = l
		} else if l == maxWordLength {
			longestWords = append(longestWords, word)
		}

		numRunes := map[int]int{}
		for _, r := range runes {
			rn := unicode.ToLower(r)
			if _, ok := numRunes[int(rn)]; ok {
				numRunes[int(rn)]++
			} else {
				numRunes[int(rn)] = 1
			}
		}
		for r, cnt := range numRunes {
			if _, ok := totalRunes[r]; !ok {
				totalRunes[r] = 0
				maxRunes[r] = 0
				wordsWithRune[r] = 0
			}
			totalRunes[r] += cnt
			if maxRunes[r] < cnt {
				maxRunes[r] = cnt
			}
			wordsWithRune[r]++
		}
	}

	maxAnagrams := 0
	var maxLeafs []anagram.Leaf
	for _, leaf := range tree.Leaves {
		if len(leaf) > maxAnagrams {
			maxLeafs = maxLeafs[:0]
			maxLeafs = append(maxLeafs, leaf)
			maxAnagrams = len(leaf)
		} else if len(leaf) == maxAnagrams {
			maxAnagrams = len(leaf)
		}
	}

	allRunes := []string{}
	for _, v := range maps.Keys(totalRunes) {
		allRunes = append(allRunes, string(rune(v)))
	}
	sort.Strings(allRunes)

	fmt.Printf("Words  : %d (%d)\n\n", numWords, numNonAnagrams)

	fmt.Printf("Words length:\n")
	for i, l := range lengthHist {
		fmt.Printf("%2d: %8d\n", i, l)
	}

	if len(longestWords) > 10 {
		fmt.Printf("Longest words: %s...\n\n", strings.Join(longestWords[:10], ", "))
	} else {
		fmt.Printf("Longest words: %s\n\n", strings.Join(longestWords, ", "))
	}

	fmt.Printf("Letters: max    total   percent    words   percent\n")
	for _, r := range allRunes {
		rn := []rune(r)[0]
		fmt.Printf(
			"  %s %8d %8d  (%5.02f%%) %8d  (%5.02f%%)\n",
			r, maxRunes[int(rn)], totalRunes[int(rn)],
			100.0*float64(totalRunes[int(rn)])/float64(totalRuneCount),
			wordsWithRune[int(rn)],
			100.0*float64(wordsWithRune[int(rn)])/float64(len(words)),
		)
	}

	fmt.Printf("\n")
	fmt.Printf("Most anagrams:\n")
	for _, leaf := range maxLeafs {
		fmt.Printf("%s\n", strings.Join(leaf, "  "))
	}
}
