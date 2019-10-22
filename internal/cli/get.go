package cli

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/mcastorina/poster/internal/models"
	"github.com/spf13/cobra"
)

const (
	wideFormat = "wide"
)

var getCmd = &cobra.Command{
	Use:     "get RESOURCE",
	Aliases: []string{"print", "g", "p"},
	Short:   "Print resources",
	Long: `Print resources.
`,
}
var getRequestCmd = &cobra.Command{
	Use:     "request",
	Aliases: []string{"requests", "req", "reqs", "r"},
	Short:   "Print request resources",
	Long: `Print request resources.
`,
	Run: getRequest,
}
var getEnvironmentCmd = &cobra.Command{
	Use:     "environment",
	Aliases: []string{"environments", "env", "envs", "e"},
	Short:   "Print environment resources",
	Long: `Print environment resources.
`,
	Run: getEnvironment,
}
var getVariableCmd = &cobra.Command{
	Use:     "variable",
	Aliases: []string{"variables", "var", "vars", "v"},
	Short:   "Print variable resources",
	Long: `Print variable resources.
`,
	Run: getVariable,
}
var tabWriter = tabwriter.NewWriter(os.Stdout, 0, 0, 6, ' ', 0)

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.AddCommand(getRequestCmd)
	getCmd.AddCommand(getEnvironmentCmd)
	getCmd.AddCommand(getVariableCmd)

	// get flags
	getCmd.PersistentFlags().StringP("output", "o", "", "Output format")
}

// run functions
func getRequest(cmd *cobra.Command, args []string) {
	outputFormat, _ := cmd.Flags().GetString("output")
	header := []interface{}{"NAME", "METHOD", "URL", "DEFAULT ENVIRONMENT"}
	if outputFormat == wideFormat {
		header = append(header, "HEADERS", "BODY")
	}
	printTableRow(header...)
	for _, request := range models.GetAllRequests() {
		row := []interface{}{request.Name, request.Method, request.URL, request.Environment.Name}
		if outputFormat == wideFormat {
			row = append(row, request.Headers, request.Body)
		}
		printTableRow(row...)
	}
	tabWriter.Flush()
}
func getEnvironment(cmd *cobra.Command, args []string) {
	printTableRow("NAME", "VARIABLES")
	for _, env := range models.GetAllEnvironments() {
		printTableRow(env.Name, env.GetVariableNames())
	}
	tabWriter.Flush()
}
func getVariable(cmd *cobra.Command, args []string) {
	outputFormat, _ := cmd.Flags().GetString("output")
	header := []interface{}{"NAME", "VALUE", "ENVIRONMENT", "TYPE"}
	if outputFormat == wideFormat {
		header = append(header, "GENERATOR", "TIMEOUT", "LAST GENERATED")
	}
	printTableRow(header...)
	for _, variable := range models.GetAllVariables() {
		value := variable.Value
		if len(value) > 50 {
			value = value[:48] + ".."
		}
		row := []interface{}{variable.Name, value, variable.Environment.Name, variable.Type}
		if outputFormat == wideFormat {
			generator := ""
			timeout := ""
			lastGenerated := ""
			if varGen := variable.Generator; varGen != nil {
				switch variable.Type {
				case models.ScriptType:
					generator = varGen.Script
				case models.RequestType:
					generator = fmt.Sprintf("%s(%s): %s", varGen.RequestName,
						varGen.RequestEnvironment, varGen.RequestPath)
				case models.ConstType:
				}
				timeout = strconv.FormatInt(varGen.Timeout, 10)
				lastGenerated = varGen.LastGenerated.Format("01/02/06 15:04:05")
			}
			row = append(row, generator, timeout, lastGenerated)
		}
		printTableRow(row...)
	}
	tabWriter.Flush()
}

// argument functions

// helper functions
func printTableRow(cols ...interface{}) {
	formatStr := ""
	for i := 0; i < len(cols); i++ {
		formatStr += "%s\t"
	}
	formatStr += "\n"

	fmt.Fprintf(tabWriter, formatStr, cols...)
}
