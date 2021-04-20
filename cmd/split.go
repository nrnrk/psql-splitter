package cmd

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/nrnrk/psql-splitter/config"
	"github.com/nrnrk/psql-splitter/domain/split"
)

var splitBy int
var outputDir string

// splitCmd represents the split command
var splitCmd = &cobra.Command{
	Use:   "split",
	Short: "Split a SQL file into multiple files.",
	Long: `Split a SQL file into multiple files.
If you want to split a file by 1,000 statements, Run

psql-splitter {target file} -n 1000

and then you can get broken down files which include 1000 statements.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires split target file name")
		}
		if len(args) >= 20 {
			return errors.New("the number of split target files must be smaller than 20")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			log.WithFields(log.Fields{
				"file": arg,
			}).Info("Splitting sqls")
			config.OutputDir = outputDir
			if err := split.Split(arg, splitBy); err != nil {
				panic(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(splitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// splitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// splitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	splitCmd.Flags().IntVarP(&splitBy, "split by", "n", 200, "the number of statements in each file")
	splitCmd.Flags().StringVarP(&outputDir, "output directory", "o", ".", "the directory to output")
}
