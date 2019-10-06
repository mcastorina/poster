package cli

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/mcastorina/poster/internal/models"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
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
	Use:     "request METHOD URL",
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
	Use:     "const-variable NAME VALUE",
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
var createScriptVariableCmd = &cobra.Command{
	Use:     "script-variable NAME GENERATOR",
	Aliases: []string{"script-var", "sv"},
	Short:   "Create a script variable resource",
	Long: `Create script-variable will create and save a script variable resource.
Variables in a request are denoted by prefixing the name with a colon
(e.g. :variable-name).

A variable resource contains the following attributes:

    name                Name of the variable
    value               Current value of the variable
    type                Type of variable (const, request, script)
    environment         Environment this variable belongs to
    generator           How to generate the value
`,
	Run:  createScriptVariable,
	Args: createScriptVariableArgs,
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.AddCommand(createRequestCmd)
	createCmd.AddCommand(createEnvironmentCmd)
	createCmd.AddCommand(createConstVariableCmd)
	createCmd.AddCommand(createScriptVariableCmd)

	// create flags
	createCmd.PersistentFlags().BoolP("interactive", "i", false, "Interactively create the resource")

	// create request flags
	createRequestCmd.Flags().StringP("name", "n", "", "Name of request for ease of use")
	createRequestCmd.Flags().StringP("environment", "e", "", "Default environment for this request")
	createRequestCmd.Flags().StringP("data", "d", "", "Request body")
	createRequestCmd.Flags().StringArrayP("header", "H", []string{}, "Request header")

	// create const-variable flags
	createConstVariableCmd.Flags().StringP("environment", "e", "", "Environment to store variable in")

	// create script-variable flags
	createScriptVariableCmd.Flags().StringP("environment", "e", "", "Environment to store variable in")
}

// run functions
func createRequest(cmd *cobra.Command, args []string) {
	if interactive, _ := cmd.Flags().GetBool("interactive"); interactive {
		createRequestI(cmd, args)
		return
	}
	name, _ := cmd.Flags().GetString("name")
	environment, _ := cmd.Flags().GetString("environment")
	body, _ := cmd.Flags().GetString("body")
	rawHeaders, _ := cmd.Flags().GetStringArray("header")

	headers := []models.Header{}
	for _, rawHeader := range rawHeaders {
		header, _ := rawHeaderToSlice(rawHeader)
		headers = append(headers, models.Header{
			Key:   header[0],
			Value: header[1],
		})
	}

	request := &models.Request{
		Name:        name,
		Method:      args[0],
		URL:         args[1],
		Environment: models.Environment{Name: environment},
		Body:        body,
		Headers:     headers,
	}
	if err := request.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "error: could not save request: %+v\n", err)
		os.Exit(1)
	}
}
func createEnvironment(cmd *cobra.Command, args []string) {
	if interactive, _ := cmd.Flags().GetBool("interactive"); interactive {
		createEnvironmentI(cmd, args)
		return
	}
	env := &models.Environment{
		Name: args[0],
	}
	if err := env.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "error: could not save environment: %+v\n", err)
		os.Exit(1)
	}
}
func createConstVariable(cmd *cobra.Command, args []string) {
	if interactive, _ := cmd.Flags().GetBool("interactive"); interactive {
		createConstVariableI(cmd, args)
		return
	}
	environment, _ := cmd.Flags().GetString("environment")
	variable := &models.Variable{
		Name:        args[0],
		Value:       args[1],
		Type:        models.ConstType,
		Environment: models.Environment{Name: environment},
	}
	if err := variable.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "error: could not save variable: %+v\n", err)
		os.Exit(1)
	}
}
func createScriptVariable(cmd *cobra.Command, args []string) {
	if interactive, _ := cmd.Flags().GetBool("interactive"); interactive {
		createScriptVariableI(cmd, args)
		return
	}
	environment, _ := cmd.Flags().GetString("environment")
	variable := &models.Variable{
		Name:        args[0],
		Type:        models.ScriptType,
		Environment: models.Environment{Name: environment},
		Generator:   args[1],
	}
	if err := variable.GenerateValue(); err != nil {
		fmt.Fprintf(os.Stderr, "error: could not generate value: %+v\n", err)
		os.Exit(1)
	}
	if err := variable.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "error: could not save variable: %+v\n", err)
		os.Exit(1)
	}
}

// argument functions
func createRequestArgs(cmd *cobra.Command, args []string) error {
	if interactive, _ := cmd.Flags().GetBool("interactive"); interactive {
		return nil
	}
	if len(args) != 2 {
		return fmt.Errorf("expected args missing: METHOD URL")
	}
	if !flagsAreSet(cmd, "name", "environment") {
		return fmt.Errorf("expected flags missing: --name, --environment")
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
	// check headers are valid (key:value)
	headers, _ := cmd.Flags().GetStringArray("header")
	for _, header := range headers {
		if _, err := rawHeaderToSlice(header); err != nil {
			return err
		}
	}
	return nil
}
func createEnvironmentArgs(cmd *cobra.Command, args []string) error {
	if interactive, _ := cmd.Flags().GetBool("interactive"); interactive {
		return nil
	}
	if len(args) != 1 {
		return fmt.Errorf("expected args missing: NAME")
	}
	return nil
}
func createConstVariableArgs(cmd *cobra.Command, args []string) error {
	if interactive, _ := cmd.Flags().GetBool("interactive"); interactive {
		return nil
	}
	if len(args) != 2 {
		return fmt.Errorf("expected args missing: NAME VALUE")
	}
	if !flagsAreSet(cmd, "environment") {
		return fmt.Errorf("expected flag missing: --environment")
	}
	return nil
}
func createScriptVariableArgs(cmd *cobra.Command, args []string) error {
	if interactive, _ := cmd.Flags().GetBool("interactive"); interactive {
		return nil
	}
	if len(args) != 2 {
		return fmt.Errorf("expected args missing: NAME GENERATOR")
	}
	if !flagsAreSet(cmd, "environment") {
		// TODO: make const
		return fmt.Errorf("expected flag missing: --environment")
	}
	return nil
}

// helper functions
func rawHeaderToSlice(header string) ([]string, error) {
	values := strings.SplitN(header, ":", 2)
	if len(values) != 2 {
		// TODO: make const
		return nil, fmt.Errorf("header should be in the format \"key:value\"")
	}
	key := strings.Trim(values[0], " \t")
	value := strings.Trim(values[1], " \t")

	if len(key) == 0 {
		// TODO: make const
		return nil, fmt.Errorf("header should be in the format \"key:value\"")
	}

	if strings.Index(header, "\n") != -1 {
		// TODO: make const
		return nil, fmt.Errorf("header should not contain newline characters")
	}

	return []string{key, value}, nil
}
func createRequestI(cmd *cobra.Command, args []string) {
	// TODO: move this to models
	template := models.Request{
		Name:        "request template",
		Method:      "GET",
		URL:         "http://localhost",
		Environment: models.Environment{Name: "local"},
		Headers: []models.Header{
			{Key: "Content-Type", Value: "application/json"},
		},
	}
	var err error
	data, _ := yaml.Marshal(template)
	data, err = updateData(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to create request: %+v\n", err)
		os.Exit(1)
	}
	if err := yaml.Unmarshal([]byte(data), &template); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to create request: %+v\n", err)
		os.Exit(1)
	}
	if err := template.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to create request: %+v\n", err)
		os.Exit(1)
	}
}
func createEnvironmentI(cmd *cobra.Command, args []string) {
	// TODO: move this to models
	template := models.Environment{
		Name: "environment template",
	}
	var err error
	data, _ := yaml.Marshal(template)
	data, err = updateData(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to create environment: %+v\n", err)
		os.Exit(1)
	}
	if err := yaml.Unmarshal([]byte(data), &template); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to create environment: %+v\n", err)
		os.Exit(1)
	}
	if err := template.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to create environment: %+v\n", err)
		os.Exit(1)
	}
}
func createConstVariableI(cmd *cobra.Command, args []string) {
	// TODO: move this to models
	template := models.Variable{
		Name:        "const variable template",
		Value:       "value",
		Environment: models.Environment{Name: "local"},
		Type:        models.ConstType,
	}
	var err error
	data, _ := yaml.Marshal(template)
	data, err = updateData(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to create variable: %+v\n", err)
		os.Exit(1)
	}
	if err := yaml.Unmarshal([]byte(data), &template); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to create variable: %+v\n", err)
		os.Exit(1)
	}
	if err := template.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to create variable: %+v\n", err)
		os.Exit(1)
	}
}
func createScriptVariableI(cmd *cobra.Command, args []string) {
	// TODO: move this to models
	template := models.Variable{
		Name:        "script variable template",
		Environment: models.Environment{Name: "local"},
		Type:        models.ScriptType,
		Generator:   `date +'%D %T'`,
	}
	var err error
	data, _ := yaml.Marshal(template)
	data, err = updateData(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to create variable: %+v\n", err)
		os.Exit(1)
	}
	if err := yaml.Unmarshal([]byte(data), &template); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to create variable: %+v\n", err)
		os.Exit(1)
	}
	if err := template.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to create variable: %+v\n", err)
		os.Exit(1)
	}
}
func flagsAreSet(cmd *cobra.Command, flagNames ...string) bool {
	if len(flagNames) == 0 {
		return true
	}

	flagsAreSet := true
	flagMap := make(map[string]bool)
	for _, flagName := range flagNames {
		flagMap[flagName] = true
	}

	cmd.Flags().VisitAll(func(pflag *flag.Flag) {
		if _, ok := flagMap[pflag.Name]; ok && !pflag.Changed {
			flagsAreSet = false
		}
	})
	return flagsAreSet
}
