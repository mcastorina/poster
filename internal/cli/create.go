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
	Long: `Create a resource. Valid resource types:
	
    request            Method, url, and default environment to run the request
    environment        Name of an environment for variable scope
    const-variable     Environment dependent constant values
`,
}
var createRequestCmd = &cobra.Command{
	Use:     "request METHOD ALIAS [PATH]",
	Aliases: []string{"req", "r"},
	Short:   "A brief description of your command",
	Long: `Create request will create and save a request resource. A request resource
contains the following attributes:

    name                Name of the request for ease of use
    method              HTTP request method
    url                 The URL path
    environment         The default environment to run the request
`,
	Run:  createRequest,
	Args: createRequestArgs,
}
var createEnvironmentCmd = &cobra.Command{
	Use:     "environment NAME",
	Aliases: []string{"env", "e"},
	Short:   "Create an environment resource",
	Long: `Create environment will create and save an environment resource. An
environment resource contains the following attributes:

    name                Name of the environment
`,
	Run:  createEnvironment,
	Args: createEnvironmentArgs,
}
var createConstVariableCmd = &cobra.Command{
	Use:     "const-variable",
	Aliases: []string{"const-var", "cv"},
	Short:   "Create a constant variable resource",
	Long: `Create const-variable will create and save a constant variable resource.
Variables in a request are denoted by prefixing the name with a colon
(e.g. :variable-name).

A variable resource contains the following attributes:

    name                Name of the variable
    value               Current value of the variable
    type                Type of variable (const, request, script)
    environment         Environment this variable belongs to
    generator           How to generate the value
`,
	Run:  createConstVariable,
	Args: createConstVariableArgs,
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.AddCommand(createRequestCmd)
	createCmd.AddCommand(createEnvironmentCmd)
	createCmd.AddCommand(createConstVariableCmd)

	// create request flags
	createRequestCmd.Flags().StringP("name", "n", "", "Name of request for ease of use")
	createRequestCmd.MarkFlagRequired("name")
	createRequestCmd.Flags().StringP("environment", "e", "", "Default environment for this request")
	createRequestCmd.MarkFlagRequired("environment")

	// create const-variable flags
	createConstVariableCmd.Flags().StringP("environment", "e", "", "Environment to store variable in")
	createConstVariableCmd.MarkFlagRequired("environment")
}

// run functions
func createRequest(cmd *cobra.Command, args []string) {
	name, _ := cmd.Flags().GetString("name")
	environment, _ := cmd.Flags().GetString("environment")
	request := &models.Request{
		Name:        name,
		Method:      args[0],
		URL:         args[1],
		Environment: models.Environment{Name: environment},
	}
	request.Save()
}
func createEnvironment(cmd *cobra.Command, args []string) {
	env := &models.Environment{
		Name: args[0],
	}
	env.Save()
}
func createConstVariable(cmd *cobra.Command, args []string) {
	environment, _ := cmd.Flags().GetString("environment")
	variable := &models.Variable{
		Name:        args[0],
		Value:       args[1],
		Type:        models.ConstType,
		Environment: models.Environment{Name: environment},
	}
	variable.Save()
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
		return fmt.Errorf("expected args missing: NAME")
	}
	return nil
}
func createConstVariableArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("expected args missing: NAME VALUE")
	}
	return nil
}
