package cli

import (
	"fmt"
	"os"

	"github.com/mcastorina/poster/internal/models"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"remove", "rm"},
	Short:   "Delete resources",
	Long: `Delete resources.
`,
}
var deleteRequestCmd = &cobra.Command{
	Use:     "request [REQUEST_NAME ...]",
	Aliases: []string{"requests", "req", "reqs", "r"},
	Short:   "Delete request resources",
	Long: `Delete request resources.
`,
	Run: deleteRequest,
}
var deleteEnvironmentCmd = &cobra.Command{
	Use:     "environment [ENVIRONMENT_NAME ...]",
	Aliases: []string{"environments", "env", "envs", "e"},
	Short:   "Delete environment resources",
	Long: `Delete environment resources.
`,
	Run: deleteEnvironment,
}
var deleteVariableCmd = &cobra.Command{
	Use:     "variable [VARIABLE_NAME ...]",
	Aliases: []string{"variables", "var", "vars", "v"},
	Short:   "Delete variable resources",
	Long: `Delete variable resources.
`,
	Run: deleteVariable,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.AddCommand(deleteRequestCmd)
	deleteCmd.AddCommand(deleteEnvironmentCmd)
	deleteCmd.AddCommand(deleteVariableCmd)
}

// run functions
func deleteRequest(cmd *cobra.Command, args []string) {
	for _, request := range args {
		mReq, err := models.GetRequestByName(request)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: failed to delete %s: %+v\n", request, err)
			os.Exit(1)
		}
		if err := mReq.Delete(); err != nil {
			fmt.Fprintf(os.Stderr, "error: failed to delete %s: %+v\n", request, err)
			os.Exit(1)
		}
	}
}
func deleteEnvironment(cmd *cobra.Command, args []string) {
	for _, environment := range args {
		mEnv, err := models.GetEnvironmentByName(environment)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: failed to delete %s: %+v\n", environment, err)
			os.Exit(1)
		}
		if err := mEnv.Delete(); err != nil {
			fmt.Fprintf(os.Stderr, "error: failed to delete %s: %+v\n", environment, err)
			os.Exit(1)
		}
	}
}
func deleteVariable(cmd *cobra.Command, args []string) {
	for _, arg := range args {
		for _, variable := range models.GetVariablesByName(arg) {
			if err := variable.Delete(); err != nil {
				fmt.Fprintf(os.Stderr, "error: failed to delete %s: %+v\n", variable, err)
				os.Exit(1)
			}
		}
	}
}

// argument functions
