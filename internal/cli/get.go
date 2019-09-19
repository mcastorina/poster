package cli

import (
	"fmt"

	"github.com/mcastorina/poster/internal/store"
	"github.com/urfave/cli"
)

func GetTarget(c *cli.Context) error {
	targets := store.GetAllTargets()
	fmt.Printf("%30s%20s\n", "ALIAS", "URL")
	for _, target := range targets {
		fmt.Printf("%30s%20s\n", target.Alias, target.URL)
	}
	return nil
}
func GetRequest(c *cli.Context) error {
	fmt.Println("Not implemented")
	return nil
}
