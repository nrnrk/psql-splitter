package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// scCmd represents the sc command
var scCmd = &cobra.Command{
	Use:   "sc",
	Short: "Count the number of statements in a SQL file.",
	Long: `Count the number of statements in a SQL file.
eg.) Check the count of statements in a SQL file
> psql-splitter sc {target file}
1000
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("sc called")
	},
}

func init() {
	rootCmd.AddCommand(scCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
