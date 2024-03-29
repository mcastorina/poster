package cli

import (
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
		log.Errorf("Failed to get request: %+v\n", err)
		os.Exit(1)
	}

	data, _ := yaml.Marshal(request)
	data, err = updateData(data)
	if err != nil {
		log.Errorf("Failed to update request: %+v\n", err)
		os.Exit(1)
	}

	newRequest := models.Request{}
	err = yaml.Unmarshal([]byte(data), &newRequest)
	if err != nil {
		log.Errorf("Failed to update request: %+v\n", err)
		os.Exit(1)
	}

	request.Delete()
	if err := newRequest.Save(); err != nil {
		request.Save() // Rollback changes
		log.Errorf("Failed to update request: %+v\n", err)
		os.Exit(1)
	}
}
func editEnvironment(cmd *cobra.Command, args []string) {
	environment, err := models.GetEnvironmentByName(args[0])
	if err != nil {
		log.Errorf("Failed to get environment: %+v\n", err)
		os.Exit(1)
	}

	data, _ := yaml.Marshal(environment)
	data, err = updateData(data)
	if err != nil {
		log.Errorf("Failed to update environment: %+v\n", err)
		os.Exit(1)
	}

	newEnvironment := models.Environment{}
	err = yaml.Unmarshal([]byte(data), &newEnvironment)
	if err != nil {
		log.Errorf("Failed to update environment: %+v\n", err)
		os.Exit(1)
	}

	// This can fail due to foreign key constraint
	if err := environment.Delete(); err != nil {
		log.Errorf("Failed to update environment: %+v\n", err)
		os.Exit(1)
	}
	if err := newEnvironment.Save(); err != nil {
		log.Errorf("Failed to update environment: %+v\n", err)
		environment.Save() // Rollback changes
		os.Exit(1)
	}
}
func editVariable(cmd *cobra.Command, args []string) {
	variables := models.GetVariablesByName(args[0])
	if len(variables) == 0 {
		log.Errorf("Failed to get variables: not found\n")
		os.Exit(1)
	}
	var err error
	data, _ := yaml.Marshal(variables)
	data, err = updateData(data)
	if err != nil {
		log.Errorf("Failed to update variables: %+v\n", err)
		os.Exit(1)
	}

	newVariables := []models.Variable{}
	err = yaml.Unmarshal([]byte(data), &newVariables)
	if err != nil {
		log.Errorf("Failed to update variables: %+v\n", err)
		os.Exit(1)
	}

	for _, variable := range variables {
		variable.Delete()
	}
	for i, newVariable := range newVariables {
		if err := newVariable.Save(); err != nil {
			// Rollback changes
			for j := 0; j < i; j++ {
				newVariables[j].Delete()
			}
			for _, variable := range variables {
				variable.Save()
			}
			log.Errorf("Failed to update variables: %+v\n", err)
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
		return nil, errorFileEmpty
	}

	// Abort on no changes
	if string(data) == string(bytes) {
		return nil, errorFileUnchanged
	}

	return bytes, nil
}
func createTmpFile(data []byte) (string, error) {
	// Create temporary file
	tmpFile, err := ioutil.TempFile(os.TempDir(), "poster-edit-")
	if err != nil {
		log.Errorf("%+v\n", err)
		return "", err
	}

	// Write to the file
	if _, err = tmpFile.Write(data); err != nil {
		log.Errorf("%+v\n", err)
		os.Remove(tmpFile.Name())
		return "", err
	}

	// Close the file
	if err := tmpFile.Close(); err != nil {
		log.Errorf("%+v\n", err)
		os.Remove(tmpFile.Name())
		return "", err
	}

	return tmpFile.Name(), nil
}
func invokeEditor(path string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return errorNoEditorFound
	}
	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
