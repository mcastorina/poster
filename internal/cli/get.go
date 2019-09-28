package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

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
var getVariableCmd = &cobra.Command{
	Use:     "variable",
	Aliases: []string{"variables", "var", "vars", "v"},
	Short:   "Print variable resources",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: getVariable,
}
var tabWriter = tabwriter.NewWriter(os.Stdout, 0, 0, 6, ' ', 0)

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.AddCommand(getRequestCmd)
	getCmd.AddCommand(getEnvironmentCmd)
	getCmd.AddCommand(getVariableCmd)
}

// run functions
func getRequest(cmd *cobra.Command, args []string) {
	fmt.Fprintf(tabWriter, "%s\t%s\t%s\t%s\t\n", "NAME", "METHOD", "URL", "DEFAULT ENVIRONMENT")
	for _, request := range models.GetAllRequests() {
		fmt.Fprintf(tabWriter, "%s\t%s\t%s\t%s\t\n", request.Name,
			request.Method, request.URL, request.Environment.Name)
	}
	tabWriter.Flush()
}
func getEnvironment(cmd *cobra.Command, args []string) {
	fmt.Fprintf(tabWriter, "%s\t\n", "NAME")
	for _, env := range models.GetAllEnvironments() {
		fmt.Fprintf(tabWriter, "%s\t\n", env.Name)
	}
	tabWriter.Flush()
}
func getVariable(cmd *cobra.Command, args []string) {
	fmt.Fprintf(tabWriter, "%s\t%s\t%s\t%s\t\n", "NAME", "VALUE", "ENVIRONMENT", "TYPE")
	for _, variable := range models.GetAllVariables() {
		fmt.Fprintf(tabWriter, "%s\t%s\t%s\t%s\t\n", variable.Name,
			variable.Value, variable.Environment.Name, variable.Type)
	}
	tabWriter.Flush()
}

// argument functions
