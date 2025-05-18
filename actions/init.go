package actions

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/optique-dev/optique"
	"github.com/optique-dev/cli/manifests"
	"github.com/optique-dev/cli/views"
)

type Initialization struct {
	URL     string
	Name    string
	Version string
}

var URL = "https://github.com/optique-dev/template"
var DEFAULT_MODULE = "github.com/optique-dev"

func NewInitialization(name string) Initialization {
	view, err := views.LaunchInitForm()
	if err != nil {
		optique.Error(fmt.Sprintf("error launching form: %s", err))
		os.Exit(1)
	}
	return Initialization{
		URL:     view.Repository,
		Name:    name,
		Version: view.Version,
	}
}



func Initialize(generation Initialization) {
	err := cloneTemplate(URL, generation.Name)
	if err != nil {
		optique.Error(fmt.Sprintf("Error cloning template: %s", err))
		os.Exit(1)
	}
	err = setupGoModule(&generation)
	if err != nil {
		optique.Error(fmt.Sprintf("Error setting up go module: %s", err))
		os.Exit(1)
	}
	err = goBack()
	if err != nil {
		optique.Error(fmt.Sprintf("Error going back: %s", err))
		os.Exit(1)
	}
}

func createProjectFolder(name string) error {
	err := os.Mkdir(name, 0755)
	if err != nil {
		return err
	}
	return nil
}

func goBack() error {
	return os.Chdir("..")
}

func cloneTemplate(url string, name string) error {
	ExecWithLoading("Cloning template", "git", "clone", url, name)

	return nil
}

var IMPORT_TO_FIX = []string{
	"./main.go",
	"./cycle.go",
}

func genProjectManifest(config *Initialization) error {
	manifest := optique.OptiqueProjectManifest{
		Name: config.Name,
		Module: config.URL,
		Ignore: []string{
			"go.mod",
			"go.sum",
		},
		Repositories: []string{},
		Applications: []string{},
	}
	manifest_json, err := json.Marshal(&manifest)
	if err != nil {
		return err
	}
	f, err := os.Create(optique.PROJECT_MANIFEST)
	defer f.Close()
	if err != nil {
		return err
	}
	_, err = f.Write(manifest_json)
	return err
}

func setupGoModule(config *Initialization) error {
	// go to project folder
	err := os.Chdir(config.Name)
	if err != nil {
		optique.Debug(fmt.Sprintf("Error changing directory: %s", err))
		return err
	}

	_, err = exec.Command("rm", "-rf", ".git").CombinedOutput()
	if err != nil {
		optique.Debug(fmt.Sprintf("Error removing .git: %s", err))
		return err
	}

	if err:= manifests.ClearIgnoredFiles(optique.PROJECT_MANIFEST); err != nil {
		return err
	}
	if err := ReplaceInAllFiles(DEFAULT_MODULE+"/template", config.URL); err != nil {
		return err
	}
	
	err = ExecWithLoading("Initializing module", "go", "mod", "init", config.URL)
	if err != nil {
		return err
	}
	optique.Info(fmt.Sprintf("Module initialized: %s\n", config.URL))

	for _, file := range IMPORT_TO_FIX {
		ExecWithLoading(fmt.Sprintf("Fixing imports for %s\n", file), "gopls", "imports", "-w", file)
	}
	ExecWithLoading("Installing dependencies", "go", "mod", "tidy")

	if err := genProjectManifest(config); err != nil {
		return err
	}

	return nil
}
