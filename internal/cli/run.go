package cli

import (
	"fmt"
	"os"

	"github.com/mcastorina/poster/internal/models"
	"github.com/spf13/cobra"
)

func run(cmd *cobra.Command, args []string) {
	for _, arg := range args {
		if resource, err := models.GetRunnableResourceByName(arg); err == nil {
			if err = resource.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "error: %+v\n", err)
				os.Exit(1)
			}
		}
	}
}
