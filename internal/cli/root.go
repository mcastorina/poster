package cli

import (
	"fmt"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "poster",
	Short: "API testing aid",
	Long: `poster is a service to quickly and easily send HTTP requests.
`,
	Run: run,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// Check if the subcommand is found; if not, add "--" to execute the run command
	// This is a hack, but it seems to be the best solution
	prevStringWasFlag := false
	for i, arg := range os.Args {
		if i == 0 || strings.HasPrefix(arg, "-") || prevStringWasFlag {
			if strings.HasPrefix(arg, "--") {
				// --flag=value should not be marked
				prevStringWasFlag = strings.Index(arg, "=") == -1
			} else {
				// -fvalue should not be marked
				prevStringWasFlag = len(arg) == 2
			}
			if arg == "-v" || arg == "--verbose" {
				// boolean flags don't have args afterwards
				prevStringWasFlag = false
			}
			continue
		}
		if _, _, err := rootCmd.Find(os.Args[i:]); err != nil {
			var args []string
			args = append(args, os.Args[:i]...)
			args = append(args, "--")
			args = append(args, os.Args[i:]...)
			os.Args = args
		}
		break
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"config file (default is $HOME/.poster.yaml)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Print verbose output")
	rootCmd.Flags().StringP("env", "e", "", "Run the resources in the specified environment")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".poster" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".poster")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
