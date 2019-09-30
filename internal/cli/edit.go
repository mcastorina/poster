package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/mcastorina/poster/internal/models"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var editCmd = &cobra.Command{
	Use:     "edit RESOURCE",
	Aliases: []string{"update", "e"},
	Short:   "Modify a resource",
	Long: `Modify a resource.
`,
}
var editRequestCmd = &cobra.Command{
	Use:     "request REQUEST_NAME",
	Aliases: []string{"req", "r"},
	Short:   "Modify a request resource",
	Long: `Modify a request resource.
`,
	Run:  editRequest,
	Args: cobra.ExactArgs(1),
}
var editEnvironmentCmd = &cobra.Command{
	Use:     "environment ENVIRONMENT_NAME",
	Aliases: []string{"env", "e"},
	Short:   "Modify a environment resource",
	Long: `Modify a environment resource.
`,
	Run:  editEnvironment,
	Args: cobra.ExactArgs(1),
}
var editVariableCmd = &cobra.Command{
	Use:     "variable VARIABLE_NAME",
	Aliases: []string{"var", "v"},
	Short:   "Modify a variable resource",
	Long: `Modify a variable resource.
`,
	Run:  editVariable,
	Args: cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(editCmd)
	editCmd.AddCommand(editRequestCmd)
	editCmd.AddCommand(editEnvironmentCmd)
	editCmd.AddCommand(editVariableCmd)
}

// run functions
func editRequest(cmd *cobra.Command, args []string) {
	request, err := models.GetRequestByName(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to get request: %+v\n", err)
		os.Exit(1)
	}

	data, _ := yaml.Marshal(request)
	data, err = updateData(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to update request: %+v\n", err)
		os.Exit(1)
	}

	newRequest := models.Request{}
	err = yaml.Unmarshal([]byte(data), &newRequest)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to update request: %+v\n", err)
		os.Exit(1)
	}

	request.Delete()
	if err := newRequest.Save(); err != nil {
		request.Save() // Rollback changes
		fmt.Fprintf(os.Stderr, "error: failed to update request: %+v\n", err)
		os.Exit(1)
	}
}
func editEnvironment(cmd *cobra.Command, args []string) {
	environment, err := models.GetEnvironmentByName(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to get environment: %+v\n", err)
		os.Exit(1)
	}

	data, _ := yaml.Marshal(environment)
	data, err = updateData(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to update environment: %+v\n", err)
		os.Exit(1)
	}

	newEnvironment := models.Environment{}
	err = yaml.Unmarshal([]byte(data), &newEnvironment)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to update environment: %+v\n", err)
		os.Exit(1)
	}

	// This can fail due to foreign key constraint
	if err := environment.Delete(); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to update environment: %+v\n", err)
		os.Exit(1)
	}
	if err := newEnvironment.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to update environment: %+v\n", err)
		environment.Save() // Rollback changes
		os.Exit(1)
	}
}
func editVariable(cmd *cobra.Command, args []string) {
	type exportedVariables struct {
		Variables []models.Variable `yaml:"variables"`
	}
	variables := exportedVariables{models.GetVariablesByName(args[0])}
	if len(variables.Variables) == 0 {
		fmt.Fprintf(os.Stderr, "error: failed to get variables: not found\n")
		os.Exit(1)
	}
	var err error
	data, _ := yaml.Marshal(variables)
	data, err = updateData(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to update variables: %+v\n", err)
		os.Exit(1)
	}

	newVariables := exportedVariables{}
	err = yaml.Unmarshal([]byte(data), &newVariables)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to update variables: %+v\n", err)
		os.Exit(1)
	}

	for _, variable := range variables.Variables {
		variable.Delete()
	}
	for i, newVariable := range newVariables.Variables {
		if err := newVariable.Save(); err != nil {
			// Rollback changes
			for j := 0; j < i; j++ {
				newVariables.Variables[j].Delete()
			}
			for _, variable := range variables.Variables {
				variable.Save()
			}
			fmt.Fprintf(os.Stderr, "error: failed to update variables: %+v\n", err)
			os.Exit(1)
		}
	}
}

// argument functions

// helper functions
func updateData(data []byte) ([]byte, error) {
	fileName, err := createTmpFile(data)
	if err != nil {
		return nil, err
	}
	defer os.Remove(fileName)

	if err := invokeEditor(fileName); err != nil {
		return nil, err
	}

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Abort on empty file
	if len(bytes) == 0 {
		// TODO: make const, and should this be an error?
		return nil, fmt.Errorf("empty file")
	}

	// Abort on no changes
	if string(data) == string(bytes) {
		// TODO: make const, and should this be an error?
		return nil, fmt.Errorf("no changes")
	}

	return bytes, nil
}
func createTmpFile(data []byte) (string, error) {
	// Create temporary file
	tmpFile, err := ioutil.TempFile(os.TempDir(), "poster-edit-")
	if err != nil {
		// TODO: log error
		return "", err
	}

	// Write to the file
	if _, err = tmpFile.Write(data); err != nil {
		// TODO: log error
		os.Remove(tmpFile.Name())
		return "", err
	}

	// Close the file
	if err := tmpFile.Close(); err != nil {
		// TODO: log error
		os.Remove(tmpFile.Name())
		return "", err
	}

	return tmpFile.Name(), nil
}
func invokeEditor(path string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return fmt.Errorf("No editor found")
	}
	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
