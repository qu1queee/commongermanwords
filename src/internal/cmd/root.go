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
		if data, err := parser.GetArticle(word); err == nil {
			file, err := os.Create(fmt.Sprintf("%v.md", word))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprintf(file, "# %s \n", word)
			fmt.Fprintf(file, "## Type \n")
			for _, ipa := range data.Type {
				fmt.Fprintf(file, "- _**%s**_ \n", ipa)
			}
			fmt.Fprintf(file, "## Pronunciation \n")
			for _, ipa := range data.IPA {
				fmt.Fprintf(file, "- _**%s**_ \n", ipa)
			}
			fmt.Fprintf(file, "## Meaning \n")
			for _, meaning := range data.Meaning {
				fmt.Fprintf(file, "- **%s** \n", meaning)
			}
			fmt.Fprintf(file, "## Examples \n")
			for _, examples := range data.Examples {
				fmt.Fprintf(file, "- **%s** \n", examples)
			}
			fmt.Fprintf(file, "## Translations \n")
			for _, translation := range data.Translation {
				fmt.Fprintf(file, "- **%s** \n", translation)
			}
			fmt.Fprintf(file, "## Features \n")
			for _, feature := range data.Features {
				fmt.Fprintf(file, "- **%s** \n", feature)
			}
		}
		os.Exit(0)
	},
}

func init() {
	rootCmd.PersistentFlags().StringP("word", "w", "", "word to look up!")
	viper.BindPFlag("word", rootCmd.PersistentFlags().Lookup("word"))

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
