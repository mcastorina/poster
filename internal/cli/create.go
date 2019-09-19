package cli

import (
	"fmt"

	"github.com/mcastorina/poster/internal/store"
	"github.com/urfave/cli"
)

func CreateTarget(c *cli.Context) error {
	return store.StoreTarget(store.Target{
		URL:   c.Args().Get(0),
		Alias: c.Args().Get(1),
	})
}
func CreateRequest(c *cli.Context) error {
	fmt.Println("Not implemented")
	return nil
}
