package cli

import (
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
	Short:   "Create a request resource",
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
var createRequestVariableCmd = &cobra.Command{
	Use:     "request-variable NAME REQUEST",
	Aliases: []string{"request-var", "req-variable", "req-var", "rv"},
	Short:   "Create a request variable resource",
	Long: `Create request-variable will create and save a request variable resource.
Variables in a request are denoted by prefixing the name with a colon
(e.g. :variable-name).

A variable resource contains the following attributes:

    name                Name of the variable
    value               Current value of the variable
    type                Type of variable (const, request, script)
    environment         Environment this variable belongs to
    generator           How to generate the value
`,
	Run:  createRequestVariable,
	Args: createRequestVariableArgs,
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.AddCommand(createRequestCmd)
	createCmd.AddCommand(createEnvironmentCmd)
	createCmd.AddCommand(createConstVariableCmd)
	createCmd.AddCommand(createScriptVariableCmd)
	createCmd.AddCommand(createRequestVariableCmd)

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

	// create request-variable flags
	createRequestVariableCmd.Flags().StringP("environment", "e", "", "Environment to store variable in")
	createRequestVariableCmd.Flags().StringP("jsonpath", "j", "", "JSONPath to extract the value from the result body")
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
		log.Errorf("Could not save request: %+v\n", err)
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
		log.Errorf("Could not save environment: %+v\n", err)
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
		log.Errorf("Could not save variable: %+v\n", err)
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
		Generator: &models.VariableGenerator{
			Script: args[1],
		},
	}
	if err := variable.Save(); err != nil {
		log.Errorf("Could not save variable: %+v\n", err)
		os.Exit(1)
	}
}
func createRequestVariable(cmd *cobra.Command, args []string) {
	if interactive, _ := cmd.Flags().GetBool("interactive"); interactive {
		createRequestVariableI(cmd, args)
		return
	}
	environment, _ := cmd.Flags().GetString("environment")
	jPath, _ := cmd.Flags().GetString("jsonpath")
	variable := &models.Variable{
		Name:        args[0],
		Type:        models.RequestType,
		Environment: models.Environment{Name: environment},
		Generator: &models.VariableGenerator{
			RequestName: args[1],
			RequestPath: jPath,
		},
	}
	if err := variable.Save(); err != nil {
		log.Errorf("Could not save variable: %+v\n", err)
		os.Exit(1)
	}
}

// argument functions
func createRequestArgs(cmd *cobra.Command, args []string) error {
	if interactive, _ := cmd.Flags().GetBool("interactive"); interactive {
		return nil
	}
	if len(args) != 2 {
		return errorMissingArgs("METHOD URL")
	}
	if !flagsAreSet(cmd, "name", "environment") {
		return errorMissingFlags("--name, --environment")
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
		return errorInvalidMethod
	}
	args[0] = strings.ToUpper(args[0])
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
		return errorMissingArg("NAME")
	}
	return nil
}
func createConstVariableArgs(cmd *cobra.Command, args []string) error {
	if interactive, _ := cmd.Flags().GetBool("interactive"); interactive {
		return nil
	}
	if len(args) != 2 {
		return errorMissingArgs("NAME VALUE")
	}
	if !flagsAreSet(cmd, "environment") {
		return errorMissingFlag("--environment")
	}
	return nil
}
func createScriptVariableArgs(cmd *cobra.Command, args []string) error {
	if interactive, _ := cmd.Flags().GetBool("interactive"); interactive {
		return nil
	}
	if len(args) != 2 {
		return errorMissingArgs("NAME GENERATOR")
	}
	if !flagsAreSet(cmd, "environment") {
		return errorMissingFlag("--environment")
	}
	return nil
}
func createRequestVariableArgs(cmd *cobra.Command, args []string) error {
	if interactive, _ := cmd.Flags().GetBool("interactive"); interactive {
		return nil
	}
	if len(args) != 2 {
		return errorMissingArgs("NAME REQUEST")
	}
	if !flagsAreSet(cmd, "environment") {
		return errorMissingFlag("--environment")
	}
	return nil
}

// helper functions
func rawHeaderToSlice(header string) ([]string, error) {
	values := strings.SplitN(header, ":", 2)
	if len(values) != 2 {
		return nil, errorInvalidHeaderFormat
	}
	key := strings.Trim(values[0], " \t")
	value := strings.Trim(values[1], " \t")

	if len(key) == 0 {
		return nil, errorInvalidHeaderFormat
	}

	if strings.Index(header, "\n") != -1 {
		return nil, errorHeaderContainsNewlineChars
	}

	return []string{key, value}, nil
}
func createRequestI(cmd *cobra.Command, args []string) {
	template := requestTemplate()
	var err error
	data, _ := yaml.Marshal(template)
	data, err = updateData(data)
	if err != nil {
		log.Errorf("Failed to create request: %+v\n", err)
		os.Exit(1)
	}
	resource := Request{}
	if err := yaml.Unmarshal([]byte(data), &resource); err != nil {
		log.Errorf("Failed to create request: %+v\n", err)
		os.Exit(1)
	}
	if err := resource.Save(); err != nil {
		log.Errorf("Failed to create request: %+v\n", err)
		os.Exit(1)
	}
}
func createEnvironmentI(cmd *cobra.Command, args []string) {
	template := environmentTemplate()
	var err error
	data, _ := yaml.Marshal(template)
	data, err = updateData(data)
	if err != nil {
		log.Errorf("Failed to create environment: %+v\n", err)
		os.Exit(1)
	}
	resource := Environment{}
	if err := yaml.Unmarshal([]byte(data), &resource); err != nil {
		log.Errorf("Failed to create environment: %+v\n", err)
		os.Exit(1)
	}
	if err := resource.Save(); err != nil {
		log.Errorf("Failed to create environment: %+v\n", err)
		os.Exit(1)
	}
}
func createConstVariableI(cmd *cobra.Command, args []string) {
	template := constVariableTemplate()
	var err error
	data, _ := yaml.Marshal(template)
	data, err = updateData(data)
	if err != nil {
		log.Errorf("Failed to create variable: %+v\n", err)
		os.Exit(1)
	}
	resource := Variable{}
	if err := yaml.Unmarshal([]byte(data), &resource); err != nil {
		log.Errorf("Failed to create variable: %+v\n", err)
		os.Exit(1)
	}
	if err := resource.Save(); err != nil {
		log.Errorf("Failed to create variable: %+v\n", err)
		os.Exit(1)
	}
}
func createScriptVariableI(cmd *cobra.Command, args []string) {
	template := scriptVariableTemplate()
	var err error
	data, _ := yaml.Marshal(template)
	data, err = updateData(data)
	if err != nil {
		log.Errorf("Failed to create variable: %+v\n", err)
		os.Exit(1)
	}
	resource := Variable{}
	if err := yaml.Unmarshal([]byte(data), &resource); err != nil {
		log.Errorf("Failed to create variable: %+v\n", err)
		os.Exit(1)
	}
	if err := resource.Save(); err != nil {
		log.Errorf("Failed to create variable: %+v\n", err)
		os.Exit(1)
	}
}
func createRequestVariableI(cmd *cobra.Command, args []string) {
	template := requestVariableTemplate()
	var err error
	data, _ := yaml.Marshal(template)
	data, err = updateData(data)
	if err != nil {
		log.Errorf("Failed to create variable: %+v\n", err)
		os.Exit(1)
	}
	resource := Variable{}
	if err := yaml.Unmarshal([]byte(data), &resource); err != nil {
		log.Errorf("Failed to create variable: %+v\n", err)
		os.Exit(1)
	}
	if err := resource.Save(); err != nil {
		log.Errorf("Failed to create variable: %+v\n", err)
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
