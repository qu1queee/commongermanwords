package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/qu1queee/commongermanwords/src/pkg/goword/parser"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "goword",
	Short: "Get the meaning of a German word",
	Long: `👾 👾 👾 👾 👾 👾 👾 👾
	
Goword generates an extensive markdown document 
of a word(German only), in order for people to 
learn new words.	
	`,
	Run: func(cmd *cobra.Command, args []string) {
		word := viper.GetString("word")
		compactEnable := viper.GetBool("compact")

		if data, err := parser.GetArticle(word); err == nil {
			if compactEnable {
				compactFile, err := os.Create(fmt.Sprintf("%v.txt", word))
				if err != nil {
					log.Fatal(err)
				}
				// todo:  Improve on full url, as its hardcoded
				// todo: this code needs to be converted to a func and minimized
				if len(data.Meaning) == 0 {
					fmt.Fprintf(compactFile, "%s /%s/\n %s\n  \nSee full definition at https://github.com/qu1queee/commongermanwords/blob/main/german/words/%v.md\n", word, data.IPA[0], data.Type[0], word)
				} else {
					fmt.Fprintf(compactFile, "%s /%s/\n %s\n  %s\n\nSee full definition at https://github.com/qu1queee/commongermanwords/blob/main/german/words/%v.md\n", word, data.IPA[0], data.Type[0], data.Meaning[data.Type[0]][0], word)
				}
				os.Exit(0)
			}

			file, err := os.Create(fmt.Sprintf("%v.md", word))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprintf(file, "# %s\n", word)
			if len(data.Translation) > 0 {
				fmt.Fprintf(file, "## Translations\n")
				for _, translation := range data.Translation {
					fmt.Fprintf(file, "- %s\n", translation)
				}
			}

			if len(data.Type) > 0 {
				fmt.Fprintf(file, "## Type\n")
				for _, ipa := range data.Type {
					fmt.Fprintf(file, "- _%s_\n", ipa)
				}
			}

			if len(data.IPA) > 0 {
				fmt.Fprintf(file, "## Pronunciation\n")
				for _, ipa := range data.IPA {
					fmt.Fprintf(file, "- **_%s_**\n", ipa)
				}
			}

			if len(data.Meaning) > 0 {
				fmt.Fprintf(file, "## Meaning\n")
				for key, meaning := range data.Meaning {
					fmt.Fprintf(file, "### %s\n", key)
					for _, lines := range meaning {
						fmt.Fprintf(file, "- %s\n", lines)
					}
				}
			}

			if len(data.Examples) > 0 {
				fmt.Fprintf(file, "## Examples\n")
				for key, examples := range data.Examples {
					fmt.Fprintf(file, "### %s\n", key)
					for _, lines := range examples {
						fmt.Fprintf(file, "- %s\n", lines)
					}
				}
			}

			if len(data.Features) > 0 {
				fmt.Fprintf(file, "## Features\n")
				for key, features := range data.Features {
					fmt.Fprintf(file, "### %s\n", key)
					for _, lines := range features {
						fmt.Fprintf(file, "- %s\n", lines)
					}
				}
			}
		}
		os.Exit(0)
	},
}

func init() {
	rootCmd.PersistentFlags().StringP("word", "w", "", "word to look up!")
	rootCmd.PersistentFlags().BoolP("compact", "c", false, "generate compact output")
	viper.BindPFlag("word", rootCmd.PersistentFlags().Lookup("word"))
	viper.BindPFlag("compact", rootCmd.PersistentFlags().Lookup("compact"))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
