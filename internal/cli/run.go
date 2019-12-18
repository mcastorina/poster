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
	Run: run,
}

func init() {
	rootCmd.AddCommand(runCmd)

	// run flags
	runCmd.Flags().StringP("env", "e", "", "Run the resources in the specified environment")
	runCmd.Flags().StringArrayP("header", "H", []string{}, "Add or overwrite request headers")
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

		// Add or override values
		if err := resource.UpdateHeaders(headers); err != nil {
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
