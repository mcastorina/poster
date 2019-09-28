package cli

import (
	"fmt"
	"os"

	"github.com/mcastorina/poster/internal/models"
	"github.com/spf13/cobra"
)

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
