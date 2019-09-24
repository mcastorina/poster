package main

import (
	"github.com/mcastorina/poster/internal/cli"
	"github.com/mcastorina/poster/internal/store"
)

func main() {
	store.InitDB()
	cli.Execute()
}
