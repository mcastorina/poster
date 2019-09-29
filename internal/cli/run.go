package cli

import (
	"fmt"
	"os"

	"github.com/mcastorina/poster/internal/models"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:     "run RESOURCE_NAME [RESOURCE_NAME ...]",
	Aliases: []string{"execute", "exec", "r"},
	Short:   "Execute the named resource",
	Long: `Run the resource.

For request types, this will send the HTTP request in the default environment,
unless overridden with the --env flag.  For suite types, it will send all HTTP
requests in the suite.
`,
	Run: run,
}

func init() {
	rootCmd.AddCommand(runCmd)

	// run flags
	runCmd.Flags().StringP("env", "e", "", "Run the resources in the specified environment")
}

func run(cmd *cobra.Command, args []string) {
	// Override environment if set
	env := models.Environment{}
	e, err := cmd.Flags().GetString("env")
	if e != "" {
		env, err = models.GetEnvironmentByName(e)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %+v\n", err)
			os.Exit(1)
		}
	}

	flags := models.FlagPrintResponseCode
	if printHeaders, _ := cmd.Flags().GetBool("verbose"); printHeaders {
		flags |= models.FlagPrintHeaders
		flags |= models.FlagPrintBody
	}

	for _, arg := range args {
		if resource, err := models.GetRunnableResourceByName(arg); err == nil {
			if env.Name == "" {
				// Use default environment
				err = resource.Run(flags)
			} else {
				// Override environment
				err = resource.RunEnv(env, flags)
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %+v\n", err)
				os.Exit(1)
			}
		}
	}
}
