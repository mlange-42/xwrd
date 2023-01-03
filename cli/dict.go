package cli

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mlange-42/xwrd/core"
	"github.com/mlange-42/xwrd/util"
	"github.com/spf13/cobra"
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

	root.AddCommand(setDictCommand(config))
	root.AddCommand(listDictsCommand(config))
	root.AddCommand(installDictCommand(config))

	return root
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

			fmt.Println("Installed:")
			for _, dict := range allDicts {
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
