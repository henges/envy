package cli

import (
	"github.com/henges/envy/internal"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "envy envfile process [args...]",
	Short: "TODO",
	Long:  `TODO`,
	Args:  cobra.MatchAll(cobra.MinimumNArgs(2), cobra.OnlyValidArgs),
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {

		envPath := args[0]
		procPath := args[1]
		var procArgs []string
		if len(args) > 2 {
			procArgs = args[2:]
		}
		return internal.RunWithEnvFile(envPath, procPath, procArgs)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.envy.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}
