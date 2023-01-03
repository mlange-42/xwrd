package cli

import (
	"fmt"

	"github.com/mlange-42/xwrd/core"
	"github.com/mlange-42/xwrd/util"
	"github.com/spf13/cobra"
)

func downloadCommand(config *core.Config) *cobra.Command {
	download := &cobra.Command{
		Use:     "download DICT",
		Short:   "Download dictionaries",
		Aliases: []string{"d"},
		Args:    util.WrappedArgs(cobra.ExactArgs(1)),
		Run: func(cmd *cobra.Command, args []string) {
			dictionary := util.NewDict(args[0])
			if util.HasDictionary(dictionary) {
				fmt.Printf("failed to download dictionary: dictionary already exists")
				return
			}

			lang, ok := util.Dictionaries[dictionary.Language]
			if !ok {
				fmt.Printf("failed to download dictionary: no dictionaries in language '%s'", dictionary.Language)
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
				fmt.Printf("failed to download dictionary: dictionary '%s' not found in language '%s'", dictionary.Name, dictionary.Language)
				return
			}

			fmt.Printf("downloading dictionary %s/%s from %s...\n", dictionary.Language, dictionary.Name, dictionary.URL)
			err := util.DownloadDictionary(dictionary)
			if err != nil {
				fmt.Printf("failed to download dictionary: %s", err.Error())
				return
			}
		},
	}
	return download
}
