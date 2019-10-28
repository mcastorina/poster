package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/mcastorina/poster/internal/models"
	"github.com/spf13/cobra"
)

const (
	wideFormat = "wide"
)

var getCmd = &cobra.Command{
	Use:     "get RESOURCE",
	Aliases: []string{"print", "g", "p"},
	Short:   "Print resources",
	Long: `Print resources.
`,
}
var getRequestCmd = &cobra.Command{
	Use:     "request",
	Aliases: []string{"requests", "req", "reqs", "r"},
	Short:   "Print request resources",
	Long: `Print request resources.
`,
	Run: getRequest,
}
var getEnvironmentCmd = &cobra.Command{
	Use:     "environment",
	Aliases: []string{"environments", "env", "envs", "e"},
	Short:   "Print environment resources",
	Long: `Print environment resources.
`,
	Run: getEnvironment,
}
var getVariableCmd = &cobra.Command{
	Use:     "variable [NAME ...]",
	Aliases: []string{"variables", "var", "vars", "v"},
	Short:   "Print variable resources",
	Long: `Print variable resources.
`,
	Run: getVariable,
}
var tabWriter = tabwriter.NewWriter(os.Stdout, 0, 0, 6, ' ', 0)

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.AddCommand(getRequestCmd)
	getCmd.AddCommand(getEnvironmentCmd)
	getCmd.AddCommand(getVariableCmd)

	// get flags
	getCmd.PersistentFlags().StringP("output", "o", "", "Output format")

	// getRequest flags
	getRequestCmd.Flags().StringP("method", "m", "", "Filter by method")
	getRequestCmd.Flags().StringP("environment", "e", "", "Filter by environment")
	getRequestCmd.Flags().StringArray("with-variable", []string{}, "Filter by request containing variable")
	getRequestCmd.Flags().StringArray("with-header", []string{}, "Filter by request containing header key or value")
	getRequestCmd.Flags().StringArray("with-body", []string{}, "Filter by request containing body")

	// getEnvironment flags
	getEnvironmentCmd.Flags().StringArray("with-variable", []string{}, "Filter by environment containing variable")

	// getVariable flags
	getVariableCmd.Flags().StringP("environment", "e", "", "Filter by environment")
	getVariableCmd.Flags().StringP("type", "t", "", "Filter by type")
}

// run functions
func getRequest(cmd *cobra.Command, args []string) {
	requests := getRequestsFromArguments(cmd, args)
	outputFormat, _ := cmd.Flags().GetString("output")
	header := []interface{}{"NAME", "METHOD", "URL", "DEFAULT ENVIRONMENT"}
	if outputFormat == wideFormat {
		header = append(header, "HEADERS", "BODY")
	}
	printTableRow(header...)
	for _, request := range requests {
		row := []interface{}{request.Name, request.Method, request.URL, request.Environment.Name}
		if outputFormat == wideFormat {
			row = append(row, request.Headers, request.Body)
		}
		printTableRow(row...)
	}
	tabWriter.Flush()
}
func getEnvironment(cmd *cobra.Command, args []string) {
	printTableRow("NAME", "VARIABLES")
	for _, env := range models.GetAllEnvironments() {
		printTableRow(env.Name, env.GetVariableNames())
	}
	tabWriter.Flush()
}
func getVariable(cmd *cobra.Command, args []string) {
	variables := getVariablesFromArguments(cmd, args)
	outputFormat, _ := cmd.Flags().GetString("output")
	header := []interface{}{"NAME", "VALUE", "ENVIRONMENT", "TYPE"}
	if outputFormat == wideFormat {
		header = append(header, "GENERATOR", "TIMEOUT", "LAST GENERATED")
	}
	printTableRow(header...)
	for _, variable := range variables {
		value := variable.Value
		if len(value) > 50 {
			value = value[:48] + ".."
		}
		row := []interface{}{variable.Name, value, variable.Environment.Name, variable.Type}
		if outputFormat == wideFormat {
			generator := ""
			timeout := ""
			lastGenerated := ""
			if varGen := variable.Generator; varGen != nil {
				switch variable.Type {
				case models.ScriptType:
					generator = varGen.Script
				case models.RequestType:
					generator = fmt.Sprintf("%s(%s): %s", varGen.RequestName,
						varGen.RequestEnvironment, varGen.RequestPath)
				case models.ConstType:
				}
				timeout = strconv.FormatInt(varGen.Timeout, 10)
				lastGenerated = varGen.LastGenerated.Format("01/02/06 15:04:05")
			}
			row = append(row, generator, timeout, lastGenerated)
		}
		printTableRow(row...)
	}
	tabWriter.Flush()
}

// argument functions

// helper functions
func printTableRow(cols ...interface{}) {
	formatStr := ""
	for i := 0; i < len(cols); i++ {
		formatStr += "%s\t"
	}
	formatStr += "\n"

	fmt.Fprintf(tabWriter, formatStr, cols...)
}
func getRequestsFromArguments(cmd *cobra.Command, args []string) []models.Request {
	methodFlag, _ := cmd.Flags().GetString("method")
	envFlag, _ := cmd.Flags().GetString("environment")
	methodFlag = strings.ToUpper(methodFlag)
	withVariables, _ := cmd.Flags().GetStringArray("with-variable")
	withHeaders, _ := cmd.Flags().GetStringArray("with-header")
	withBodies, _ := cmd.Flags().GetStringArray("with-body")
	requestArr := []models.Request{}
	if envFlag != "" && methodFlag != "" {
		if len(args) == 0 {
			requestArr = models.GetRequestsByEnvironmentAndMethod(envFlag, methodFlag)
		} else {
			for _, arg := range args {
				if request, err := models.GetRequestByName(arg); err == nil &&
					request.Environment.Name == envFlag && request.Method == methodFlag {
					requestArr = append(requestArr, request)
				}
			}
		}
	} else if envFlag != "" {
		if len(args) == 0 {
			requestArr = models.GetRequestsByEnvironment(envFlag)
		} else {
			for _, arg := range args {
				if request, err := models.GetRequestByName(arg); err == nil &&
					request.Environment.Name == envFlag {
					requestArr = append(requestArr, request)
				}
			}
		}
	} else if methodFlag != "" {
		if len(args) == 0 {
			requestArr = models.GetRequestsByMethod(methodFlag)
		} else {
			for _, arg := range args {
				if request, err := models.GetRequestByName(arg); err == nil &&
					request.Method == methodFlag {
					requestArr = append(requestArr, request)
				}
			}
		}
	} else {
		if len(args) == 0 {
			requestArr = models.GetAllRequests()
		} else {
			for _, arg := range args {
				if request, err := models.GetRequestByName(arg); err == nil {
					requestArr = append(requestArr, request)
				}
			}
		}
	}

	requests := requestArr
	if len(withBodies)+len(withHeaders)+len(withVariables) > 0 {
		requestMap := make(map[string]models.Request)
		for _, request := range requestArr {
			requestMap[request.Name] = request
		}

		if len(withBodies) > 0 {
			// TODO: Remove nested loops
			for _, request := range requestMap {
				for _, withBody := range withBodies {
					if !strings.Contains(request.Body, withBody) {
						delete(requestMap, request.Name)
					}
				}
			}
		}
		if len(withHeaders) > 0 {
			// TODO: Remove nested loops
			for _, request := range requestMap {
				headerMap := make(map[string]bool)
				for _, header := range request.Headers {
					headerMap[header.Key] = true
					headerMap[header.Value] = true
				}
				for _, withHeader := range withHeaders {
					if !headerMap[withHeader] {
						delete(requestMap, request.Name)
					}
				}
			}
		}
		if len(withVariables) > 0 {
			// TODO: Remove nested loops
			for _, request := range requestMap {
				variableMap := make(map[string]bool)
				for _, variable := range request.Environment.GetVariablesInRequest(&request) {
					variableMap[variable.Name] = true
				}
				for _, withVar := range withVariables {
					if !variableMap[withVar] {
						delete(requestMap, request.Name)
					}
				}
			}
		}
		requests = []models.Request{}
		for _, request := range requestMap {
			requests = append(requests, request)
		}
	}

	return requests
}
func getEnvironmentsFromArguments(cmd *cobra.Command, args []string) []models.Environment {
	withVariables, _ := cmd.Flags().GetStringArray("with-variable")

	environments := []models.Environment{}
	if len(args) > 0 {
		for _, arg := range args {
			if environment, err := models.GetEnvironmentByName(arg); err == nil {
				environments = append(environments, environment)
			}
		}
	} else {
		environments = models.GetAllEnvironments()
	}
	if len(withVariables) > 0 {
		envMap := make(map[string]models.Environment)
		for _, env := range environments {
			envMap[env.Name] = env
		}

		// TODO: Remove nested loops
		for _, environment := range envMap {
			variableMap := make(map[string]bool)
			for _, variableName := range environment.GetVariableNames() {
				variableMap[variableName] = true
			}
			for _, withVar := range withVariables {
				if !variableMap[withVar] {
					delete(envMap, environment.Name)
				}
			}
		}

		environments = []models.Environment{}
		for _, environment := range envMap {
			environments = append(environments, environment)
		}
	}
	return environments
}
func getVariablesFromArguments(cmd *cobra.Command, args []string) []models.Variable {
	envFlag, _ := cmd.Flags().GetString("environment")
	typeFlag, _ := cmd.Flags().GetString("type")

	// Most specific to least specific
	if envFlag != "" && typeFlag != "" {
		if len(args) == 0 {
			return models.GetVariablesByEnvironmentAndType(envFlag, typeFlag)
		}
		variables := []models.Variable{}
		for _, arg := range args {
			if variable, err := models.GetVariableByNameAndEnvironment(arg, envFlag); err == nil && variable.Type == typeFlag {
				variables = append(variables, variable)
			}
		}
		return variables
	}
	if envFlag != "" {
		if len(args) == 0 {
			return models.GetVariablesByEnvironment(envFlag)
		}
		variables := []models.Variable{}
		for _, arg := range args {
			if variable, err := models.GetVariableByNameAndEnvironment(arg, envFlag); err == nil {
				variables = append(variables, variable)
			}
		}
		return variables
	}
	if typeFlag != "" {
		if len(args) == 0 {
			return models.GetVariablesByType(typeFlag)
		}
		variables := []models.Variable{}
		for _, arg := range args {
			variables = append(variables, models.GetVariablesByNameAndType(arg, typeFlag)...)
		}
		return variables
	}
	if len(args) == 0 {
		return models.GetAllVariables()
	}
	variables := []models.Variable{}
	for _, arg := range args {
		variables = append(variables, models.GetVariablesByName(arg)...)
	}
	return variables
}
