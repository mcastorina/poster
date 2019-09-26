package cli

import (
	"fmt"

	"github.com/mcastorina/poster/internal/models"
	"github.com/spf13/cobra"
)

func run(cmd *cobra.Command, args []string) {
	for _, arg := range args {
		if resource, err := models.GetRunnableResourceByName(arg); err == nil {
			// TODO: move this message to within model
			fmt.Printf(" * Running %s...", arg)
			resource.Run()
			fmt.Println("Done")
		}
		// TODO: notify user of failure
	}
}
