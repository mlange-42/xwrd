package cli

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/mlange-42/xwrd/anagram"
	"github.com/mlange-42/xwrd/core"
	"github.com/mlange-42/xwrd/util"
	"github.com/spf13/cobra"
)

type anagramOptions struct {
	partial    bool
	multi      bool
	maxWords   uint
	minLength  uint
	filter     string
	pattern    *regexp.Regexp
	unknown    []uint
	minUnknown uint
	maxUnknown uint
}

func anagramCommand(config *core.Config) *cobra.Command {
	op := anagramOptions{}
	var dict string

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
			if !op.multi && cmd.Flags().Changed("max-words") {
				fmt.Print("ERROR: flag --max-words is only supported with flag --multi")
				return
			}
			if !op.multi && !op.partial && op.minLength > 0 {
				fmt.Print("ERROR: flag --min-length is only supported with flag --multi or --partial")
				return
			}

			var err error
			op.minUnknown, op.maxUnknown, err = parseUnknown(op.unknown, op.multi)
			if err != nil {
				fmt.Printf("ERROR: %s", err.Error())
			}

			dictionary := config.GetDict()
			if dict != "" {
				dictionary = util.NewDict(dict)
			}
			words, err := util.LoadDictionary(dictionary)
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

			if op.filter != "" {
				op.pattern, err = createPattern(op.filter)
				if err != nil {
					fmt.Printf("failed to find anagrams: %s", err.Error())
					return
				}
			}

			commands := map[string]bool{
				"filter":     true,
				"max-words":  true,
				"min-length": true,
				"unknown":    true,
				"f":          true,
				"w":          true,
				"l":          true,
				"u":          true,
			}

			if interactive {
				if op.multi {
					fmt.Print("Find multi-word anagrams.")
				} else if op.partial {
					fmt.Print("Find partial anagrams.")
				} else {
					fmt.Print("Find anagrams.")
				}
				fmt.Println(" Enter ? for help.")
			}
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
					if str, ok := interactiveFlags(answer, &op, commands); ok {
						fmt.Println(str)
						continue
					}
					text = []string{answer}
				} else {
					text = args
				}

				for _, word := range text {
					if !interactive {
						fmt.Printf("%s:\n", word)
					}
					if op.partial {
						ana := tree.PartialAnagramsWithUnknown(word, op.minLength, op.minUnknown, op.maxUnknown)
						for _, res := range ana {
							line := printFiltered(res, op.pattern)
							if line != "" {
								fmt.Printf("  %s\n", line)
							}
						}
					} else if op.multi {
						ana := tree.MultiAnagrams(word, op.maxWords, op.minLength, false)

						for _, res := range ana {
							found := op.pattern == nil
							foundIndex := -1
							if op.pattern != nil {
							FindMatch:
								for b, block := range res {
									for _, word := range block {
										if op.pattern.MatchString(word) {
											found = true
											foundIndex = b
											break FindMatch
										}
									}
								}
							}
							if !found {
								continue
							}

							fmt.Print("  ")
							for b, block := range res {
								if b == foundIndex {
									fmt.Print(printFiltered(block, op.pattern))
								} else {
									fmt.Print(strings.Join(block, "  "))
								}
								if b < len(res)-1 {
									fmt.Print("  |  ")
								}
							}
							fmt.Println()
						}
					} else {
						ana := tree.AnagramsWithUnknown(word, op.minUnknown, op.maxUnknown)
						for _, res := range ana {
							line := printFiltered(res, op.pattern)
							if line != "" {
								fmt.Printf("  %s\n", line)
							}
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

	anagram.Flags().BoolVarP(&op.partial, "partial", "p", false, "Find partial anagrams.")
	anagram.Flags().BoolVarP(&op.multi, "multi", "m", false, "Find combinations of multiple partial anagrams.")

	anagram.Flags().UintVarP(&op.maxWords, "max-words", "w", 0, "Word count limit for multi-anagrams.")
	anagram.Flags().UintVarP(&op.minLength, "min-length", "l", 0, "Minimum word length for partial and multi-anagrams.")

	anagram.Flags().UintSliceVarP(&op.unknown, "unknown", "u", []uint{}, "Number of unknown/open letters ([min,]max).\nUse a single number like '1' for an exact number of unknowns.\nOtherwise, use a range like '0,2'")

	anagram.Flags().StringVarP(&op.filter, "filter", "f", "", "Pattern for filtering anagrams.")

	anagram.MarkFlagsMutuallyExclusive("partial", "multi")

	return anagram
}

func interactiveFlags(answer string, op *anagramOptions, commands map[string]bool) (string, bool) {
	if answer == "?" {
		sb := strings.Builder{}
		fmt.Fprintln(&sb, "")
		fmt.Fprintln(&sb, "To change the mode, enter #normal, #partial or #multi.")
		fmt.Fprintln(&sb, "")
		fmt.Fprintln(&sb, "To change flags, enter the flag's name and the value, separated by '='")
		fmt.Fprintln(&sb, "Available flags with current setting:")
		fmt.Fprintln(&sb, "")
		fmt.Fprintf(&sb, "  filter = %s\n", op.filter)
		if !op.multi {
			fmt.Fprintf(&sb, "  unknown = %d,%d\n", op.minUnknown, op.maxUnknown)
		} else {
			fmt.Fprintf(&sb, "  unknown = %d,%d    (*)\n", op.minUnknown, op.maxUnknown)
		}
		if op.multi {
			fmt.Fprintf(&sb, "  max-words = %d\n", op.maxWords)
		} else {
			fmt.Fprintf(&sb, "  max-words = %d    (*)\n", op.maxWords)
		}
		if op.multi || op.partial {
			fmt.Fprintf(&sb, "  min-length = %d\n", op.minLength)
		} else {
			fmt.Fprintf(&sb, "  min-length = %d   (*)\n", op.minLength)
		}

		mode := "#normal"
		if op.partial {
			mode = "#partial"
		} else if op.multi {
			mode = "#multi"
		}

		fmt.Fprintf(&sb, "  (*)...ignored in mode %s\n", mode)
		fmt.Fprintln(&sb, "")
		fmt.Fprintln(&sb, "To quit, enter nothing or press Ctrl+C")
		return sb.String(), true
	}
	if strings.HasPrefix(answer, "#") {
		switch answer {
		case "#normal", "#n":
			op.partial = false
			op.multi = false
			return "switched to mode #normal", true
		case "#partial", "#p":
			op.partial = true
			op.multi = false
			return "switched to mode #partial", true
		case "#multi", "#m":
			op.partial = false
			op.multi = true
			return "switched to mode #multi", true
		default:
			return fmt.Sprintf("failed to set mode: unknown mode #%s", answer), true
		}

	}
	if strings.Contains(answer, "=") {
		parts := strings.SplitN(answer, "=", 2)
		if len(parts) < 2 {
			return "failed to set flags: no value provided", true
		}
		command := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if !commands[command] {
			return fmt.Sprintf("failed to set flags: unknown flag #%s", command), true
		}
		switch command {
		case "filter", "f":
			op.filter = value
			pat, err := createPattern(op.filter)
			if err != nil {
				return fmt.Sprintf("failed to set filter: %s", err.Error()), true
			}
			if op.filter == "" {
				op.pattern = nil
			} else {
				op.pattern = pat
			}
			return fmt.Sprintf("set filter=%s", op.filter), true
		case "max-words", "w":
			max, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Sprintf("failed to set max-words: %s", err.Error()), true
			}
			op.maxWords = uint(max)
			return fmt.Sprintf("set max-words=%d", op.maxWords), true
		case "min-length", "l":
			min, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Sprintf("failed to set min-length: %s", err.Error()), true
			}
			op.minLength = uint(min)
			return fmt.Sprintf("set min-length=%d", op.minLength), true
		case "unknown", "u":
			min, max, err := parseUnknownStr(value, op.multi)
			if err != nil {
				return fmt.Sprintf("failed to set min-length: %s", err.Error()), true
			}
			op.minUnknown, op.maxUnknown = min, max
			return fmt.Sprintf("set unknown=%d,%d", op.minUnknown, op.maxUnknown), true
		default:
			return fmt.Sprintf("failed to set flags: unknown flag #%s", command), true
		}
	}
	return "", false
}

func parseUnknownStr(unknownStr string, multi bool) (uint, uint, error) {
	parts := strings.Split(unknownStr, ",")
	unknown := make([]uint, len(parts), len(parts))
	for i, p := range parts {
		val, err := strconv.Atoi(p)
		if err != nil {
			return 0, 0, err
		}
		unknown[i] = uint(val)
	}
	return parseUnknown(unknown, multi)
}
func parseUnknown(unknown []uint, multi bool) (uint, uint, error) {
	var minUnknown uint = 0
	var maxUnknown uint = 0

	if len(unknown) > 0 {
		if multi {
			return 0, 0, fmt.Errorf("flag --unknown is not supported with flags --multi")
		}
		switch len(unknown) {
		case 0:
		case 1:
			minUnknown = unknown[0]
			maxUnknown = unknown[0]
		case 2:
			minUnknown = unknown[0]
			maxUnknown = unknown[1]
			if minUnknown > maxUnknown {
				return 0, 0, fmt.Errorf("flag --unknown - 2nd argument must not be larger than 1st argument")
			}
		default:
			return 0, 0, fmt.Errorf("flag --unknown expects one or two arguments")
		}
	}
	return minUnknown, maxUnknown, nil
}

func printFiltered(leaf anagram.Leaf, pattern *regexp.Regexp) string {
	temp := []string{}
	for _, word := range leaf {
		if pattern == nil || pattern.MatchString(word) {
			temp = append(temp, word)
		}
	}
	if len(temp) == 0 {
		return ""
	}
	return strings.Join(temp, "  ")
}
