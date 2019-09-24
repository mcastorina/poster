package cli

import (
	"fmt"
	"strings"

	"github.com/mcastorina/poster/internal/store"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:     "create RESOURCE",
	Aliases: []string{"c", "add", "a"},
	Short:   "Create a resource",
	Long:    `Create a resource. Valid resource types: [target, request, environment]`,
}
var createRequestCmd = &cobra.Command{
	Use:     "request METHOD ALIAS",
	Aliases: []string{"req", "r"},
	Short:   "A brief description of your command",
	Long: `Create request will create and save a request resource. A request resource
contains the following attributes:

    name                Name of the request for ease of use
    method              HTTP request method
    target              The target alias
`,
	Run:  createRequest,
	Args: createRequestArgs,
}
var createTargetCmd = &cobra.Command{
	Use:     "target URL",
	Aliases: []string{"t"},
	Short:   "Create a target resource",
	Long: `Create target will create and save a target resource. A target resource
contains the following attributes:

    url                 Endpoint URL
    alias               Name of the target for ease of use
`,
	Run:  createTarget,
	Args: createTargetArgs,
}
var createEnvironmentCmd = &cobra.Command{
	Use:     "environment name",
	Aliases: []string{"env", "e"},
	Short:   "Create an environment resource",
	Long: `Create environment will create and save an environment resource. An
environment resource contains the following attributes:

    name                Name of the environment
`,
	Run:  createEnvironment,
	Args: createEnvironmentArgs,
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.AddCommand(createRequestCmd)
	createCmd.AddCommand(createTargetCmd)
	createCmd.AddCommand(createEnvironmentCmd)

	// create request flags
	createRequestCmd.Flags().StringP("name", "n", "", "Name of request for ease of use")
	createRequestCmd.MarkFlagRequired("name")

	// create target flags
	createTargetCmd.Flags().StringP("alias", "a", "", "Help message for alias")
}

// run functions
func createRequest(cmd *cobra.Command, args []string) {
	name, _ := cmd.Flags().GetString("name")
	store.StoreRequest(store.RequestType{
		Name:   name,
		Method: args[0],
		Target: store.TargetType{
			Alias: args[1],
		},
	})
}
func createTarget(cmd *cobra.Command, args []string) {
	alias, _ := cmd.Flags().GetString("alias")
	if alias == "" {
		alias = args[0]
	}
	store.StoreTarget(store.TargetType{
		URL:   args[0],
		Alias: alias,
	})
}
func createEnvironment(cmd *cobra.Command, args []string) {
	store.StoreEnvironment(store.EnvironmentType{
		Name: args[0],
	})
}

// argument functions
func createRequestArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("expected args missing: METHOD ALIAS")
	}
	validMethods := map[string]bool{
		"GET":     true,
		"HEAD":    true,
		"POST":    true,
		"PUT":     true,
		"DELETE":  true,
		"CONNECT": true,
		"OPTIONS": true,
		"TRACE":   true,
	}
	if _, ok := validMethods[strings.ToUpper(args[0])]; !ok {
		validMethodsArray := make([]string, 0, len(validMethods))
		for key := range validMethods {
			validMethodsArray = append(validMethodsArray, key)
		}
		return fmt.Errorf("METHOD \"%s\" not recognized. valid methods: %+v",
			args[0], validMethodsArray)
	}
	args[0] = strings.ToUpper(args[0])
	return nil
}
func createTargetArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("expected args missing: URL")
	}
	return nil
}
func createEnvironmentArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("expected args missing: name")
	}
	return nil
}
