package cli

import (
	"fmt"

	"github.com/mcastorina/poster/internal/store"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}
var getRequestCmd = &cobra.Command{
	Use:     "request",
	Aliases: []string{"requests", "req", "reqs", "r"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: getRequest,
}
var getTargetCmd = &cobra.Command{
	Use:     "target",
	Aliases: []string{"targets", "t"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: getTarget,
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.AddCommand(getRequestCmd)
	getCmd.AddCommand(getTargetCmd)

	getRequestCmd.Flags().StringP("alias", "a", "", "Help message for alias")
}

// run functions
func getRequest(cmd *cobra.Command, args []string) {
	requests := store.GetAllRequests()
	fmt.Printf("%20s%20s%20s\n", "NAME", "METHOD", "URL")
	for _, request := range requests {
		fmt.Printf("%20s%20s%20s\n", request.Name, request.Method, request.Target.URL)
	}
}
func getTarget(cmd *cobra.Command, args []string) {
	targets := store.GetAllTargets()
	fmt.Printf("%30s%20s\n", "ALIAS", "URL")
	for _, target := range targets {
		fmt.Printf("%30s%20s\n", target.Alias, target.URL)
	}
}

// argument functions
