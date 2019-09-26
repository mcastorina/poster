package cli

import (
	"fmt"

	"github.com/mcastorina/poster/internal/models"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Print resources",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}
var getRequestCmd = &cobra.Command{
	Use:     "request",
	Aliases: []string{"requests", "req", "reqs", "r"},
	Short:   "Print request resources",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: getRequest,
}
var getTargetCmd = &cobra.Command{
	Use:     "target [alias | url] ...",
	Aliases: []string{"targets", "t"},
	Short:   "Print target resources",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: getTarget,
}
var getEnvironmentCmd = &cobra.Command{
	Use:     "environment",
	Aliases: []string{"environments", "env", "envs", "e"},
	Short:   "Print environment resources",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: getEnvironment,
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.AddCommand(getRequestCmd)
	getCmd.AddCommand(getTargetCmd)
	getCmd.AddCommand(getEnvironmentCmd)
}

// run functions
func getRequest(cmd *cobra.Command, args []string) {
	requests := models.GetAllRequests()
	fmt.Printf("%20s%20s%20s%20s\n", "NAME", "METHOD", "URL", "PATH")
	for _, request := range requests {
		fmt.Printf("%20s%20s%20s%20s\n", request.Name,
			request.Method, request.Target.URL, request.Path)
	}
}
func getTarget(cmd *cobra.Command, args []string) {
	fmt.Printf("%30s%20s\n", "ALIAS", "URL")
	if len(args) == 0 {
		targets := models.GetAllTargets()
		for _, target := range targets {
			fmt.Printf("%30s%20s\n", target.Alias, target.URL)
		}
		return
	}

	for _, arg := range args {
		if target, err := models.GetTargetByAlias(arg); err == nil {
			fmt.Printf("%30s%20s\n", target.Alias, target.URL)
		} else if target, err := models.GetTargetByURL(arg); err == nil {
			fmt.Printf("%30s%20s\n", target.Alias, target.URL)
		}
	}
}
func getEnvironment(cmd *cobra.Command, args []string) {
	fmt.Printf("%20s\n", "NAME")
	envs := models.GetAllEnvironments()
	for _, env := range envs {
		fmt.Printf("%20s\n", env.Name)
	}
}

// argument functions
