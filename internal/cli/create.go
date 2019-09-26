package cli

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/mcastorina/poster/internal/models"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:     "create RESOURCE",
	Aliases: []string{"c", "add", "a"},
	Short:   "Create a resource",
	Long:    `Create a resource. Valid resource types: [target, request, environment]`,
}
var createRequestCmd = &cobra.Command{
	Use:     "request METHOD ALIAS [PATH]",
	Aliases: []string{"req", "r"},
	Short:   "A brief description of your command",
	Long: `Create request will create and save a request resource. A request resource
contains the following attributes:

    name                Name of the request for ease of use
    method              HTTP request method
    target              The target alias
    path                The URL path
`,
	Run:  createRequest,
	Args: createRequestArgs,
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
	createCmd.AddCommand(createEnvironmentCmd)

	// create request flags
	createRequestCmd.Flags().StringP("name", "n", "", "Name of request for ease of use")
	createRequestCmd.MarkFlagRequired("name")
}

// run functions
func createRequest(cmd *cobra.Command, args []string) {
	// add '/' as default arg
	name, _ := cmd.Flags().GetString("name")
	request := &models.Request{
		Name:   name,
		Method: args[0],
		URL:    args[1],
	}
	request.Save()
}
func createEnvironment(cmd *cobra.Command, args []string) {
	env := &models.Environment{
		Name: args[0],
	}
	env.Save()
}

// argument functions
func createRequestArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("expected args missing: METHOD URL")
	}
	// check method is valid
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
	// check url is valid
	if !strings.Contains(args[1], "//") {
		args[1] = fmt.Sprintf("//%s", args[1])
	}
	urlObj, err := url.Parse(args[1])
	if err != nil {
		return err
	}
	if urlObj.Scheme == "" {
		urlObj.Scheme = "http"
	}
	args[1] = urlObj.String()
	return nil
}
func createEnvironmentArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("expected args missing: name")
	}
	return nil
}
