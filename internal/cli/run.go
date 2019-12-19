package cli

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/mcastorina/poster/internal/models"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:     "run RESOURCE_NAME [RESOURCE_NAME ...]",
	Aliases: []string{"execute", "exec", "r"},
	Short:   "Execute the named resource",
	Long: `Run the resource.

For request types, this will send the HTTP request in the default environment,
unless overridden with the --env flag.  For suite types, it will send all HTTP
requests in the suite.

All parts of the resource will be parsed for variables and replaced with their
current value.
`,
	Run:  run,
	Args: runArgs,
}

func init() {
	rootCmd.AddCommand(runCmd)

	// run flags
	runCmd.Flags().StringP("env", "e", "", "Run the resources in the specified environment")
	runCmd.Flags().StringArrayP("header", "H", []string{}, "Add or overwrite request headers")
	runCmd.Flags().StringP("data", "d", "", "Add or overwrite the request body")
	runCmd.Flags().StringArrayP("variable", "V", []string{}, "Add or overwrite request variables")
}

func run(cmd *cobra.Command, args []string) {
	// Override environment if set
	env := models.Environment{}
	e, err := cmd.Flags().GetString("env")
	if e != "" {
		env, err = models.GetEnvironmentByName(e)
		if err != nil {
			log.Errorf("%+v\n", err)
			os.Exit(1)
		}
	}

	verboseFlag, _ := cmd.Flags().GetBool("verbose")

	for _, arg := range args {
		resource, err := models.GetRunnableResourceByName(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "eror: could not run %s: %+v\n", arg, err)
			os.Exit(1)
		}
		// Get header flags
		rawHeaders, _ := cmd.Flags().GetStringArray("header")
		headers := []models.Header{}
		for _, rawHeader := range rawHeaders {
			header, _ := rawHeaderToSlice(rawHeader)
			headers = append(headers, models.Header{
				Key:   header[0],
				Value: header[1],
			})
		}
		// Get body flag
		data, _ := cmd.Flags().GetString("data")
		// Get variable flags
		rawVariables, _ := cmd.Flags().GetStringArray("variable")
		variables := []models.Variable{}
		for _, rawVariable := range rawVariables {
			variable, _ := rawVariableToSlice(rawVariable)
			variables = append(variables, models.Variable{
				Name:  variable[0],
				Value: variable[1],
				Type:  models.ConstType,
			})
		}

		// Add or override values
		if err := resource.UpdateHeaders(headers); err != nil {
			log.Errorf("Could not update headers for %s: %+v\n", arg, err)
			os.Exit(1)
		}
		if data != "" {
			if err := resource.UpdateBody(data); err != nil {
				log.Errorf("Could not update the body for %s: %+v\n", arg, err)
				os.Exit(1)
			}
		}
		if err := resource.UpdateVariables(variables); err != nil {
			log.Errorf("Could not update headers for %s: %+v\n", arg, err)
			os.Exit(1)
		}

		var resp *http.Response
		if env.Name == "" {
			// Use default environment
			resp, err = resource.Run()
		} else {
			// Override environment
			resp, err = resource.RunEnv(env)
		}
		if err != nil {
			log.Errorf("Could not run %s: %+v\n", arg, err)
			os.Exit(1)
		}

		printResponse(resp, verboseFlag)
	}
}

// argument functions
func runArgs(cmd *cobra.Command, args []string) error {
	// check headers are valid (key:value)
	headers, _ := cmd.Flags().GetStringArray("header")
	for _, header := range headers {
		if _, err := rawHeaderToSlice(header); err != nil {
			return err
		}
	}
	// check variables are valid (key=value)
	variables, _ := cmd.Flags().GetStringArray("variable")
	for _, variable := range variables {
		if _, err := rawVariableToSlice(variable); err != nil {
			return err
		}
	}
	return nil
}

// helper functions
func printResponse(resp *http.Response, verbose bool) error {
	// print response code
	fmt.Fprintf(os.Stderr, "< %s %s\n", resp.Proto, resp.Status)

	// print headers
	if verbose {
		for header, values := range resp.Header {
			fmt.Fprintf(os.Stderr, "< %s: %s\n", header, strings.Join(values, ","))
		}
	}

	// print body
	if verbose {
		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return err
		}
		fmt.Printf("%s", body)
	}

	return nil
}
func rawVariableToSlice(variable string) ([]string, error) {
	values := strings.SplitN(variable, "=", 2)
	if len(values) != 2 {
		return nil, errorInvalidVariableFormat
	}
	key := strings.Trim(values[0], " \t")
	value := strings.Trim(values[1], " \t")

	if len(key) == 0 {
		return nil, errorInvalidVariableFormat
	}

	return []string{key, value}, nil
}
