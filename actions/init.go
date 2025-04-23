package actions

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/optique-dev/cli/manifests"
	"github.com/optique-dev/cli/views"
	"github.com/optique-dev/core"
)

type Initialization struct {
	URL     string
	Name    string
	Version string
}

var URL = "https://github.com/Courtcircuits/optique"
var DEFAULT_MODULE = "github.com/Courtcircuits/optique"

func NewInitialization(name string) Initialization {
	view, err := views.LaunchInitForm()
	if err != nil {
		core.Error(fmt.Sprintf("error launching form: %s", err))
		os.Exit(1)
	}
	return Initialization{
		URL:     view.Repository,
		Name:    name,
		Version: view.Version,
	}
}



func Initialize(generation Initialization) {
	err := createProjectFolder(generation.Name)
	if err != nil {
		core.Error(fmt.Sprintf("Error creating project folder: %s", err))
		os.Exit(1)
	}
	err = cloneTemplate(URL, generation.Name)
	if err != nil {
		core.Error(fmt.Sprintf("Error cloning template: %s", err))
		os.Exit(1)
	}
	err = setupGoModule(&generation)
	if err != nil {
		core.Error(fmt.Sprintf("Error setting up go module: %s", err))
		os.Exit(1)
	}
	err = goBack()
	if err != nil {
		core.Error(fmt.Sprintf("Error going back: %s", err))
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

	// go to project folder
	err := os.Chdir(name)
	if err != nil {
		return err
	}

	current_dir, err := os.Getwd()
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(current_dir)

	if err != nil {
		return err
	}

	folders_to_delete := []string{}
	files_to_delete := []string{}
	for _, entry := range entries {
		if entry.IsDir() {
			if entry.Name() != "template" {
				folders_to_delete = append(folders_to_delete, entry.Name())
			}
		} else {
			files_to_delete = append(files_to_delete, entry.Name())
		}
	}

	for _, entry := range folders_to_delete {
		err = os.RemoveAll(entry)
		if err != nil {
			return err
		}
	}
	for _, entry := range files_to_delete {
		err = os.Remove(entry)
		if err != nil {
			return err
		}
	}

	// go to template folder
	err = os.Chdir("template")
	if err != nil {
		return err
	}

	entries, err = os.ReadDir(".")
	for _, entry := range entries {
		_, err := exec.Command("mv", entry.Name(), current_dir).CombinedOutput()
		if err != nil {
			return err
		}
	}
	// move all to parent folder
	err = goBack()
	if err != nil {
		return err
	}

	// remove template folder
	err = os.RemoveAll("template")
	return nil
}

var IMPORT_TO_FIX = []string{
	"./main.go",
	"./cycle.go",
}

func genProjectManifest(config *Initialization) error {
	manifest := core.OptiqueProjectManifest{
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
	f, err := os.Create(core.PROJECT_MANIFEST)
	defer f.Close()
	if err != nil {
		return err
	}
	_, err = f.Write(manifest_json)
	return err
}

func setupGoModule(config *Initialization) error {
	// go to project folder

	if err:= manifests.ClearIgnoredFiles(core.PROJECT_MANIFEST); err != nil {
		return err
	}
	if err := ReplaceInAllFiles(DEFAULT_MODULE+"/template", config.URL); err != nil {
		return err
	}
	
	err := ExecWithLoading("Initializing module", "go", "mod", "init", config.URL)
	if err != nil {
		return err
	}
	core.Info(fmt.Sprintf("Module initialized: %s\n", config.URL))

	for _, file := range IMPORT_TO_FIX {
		ExecWithLoading(fmt.Sprintf("Fixing imports for %s\n", file), "gopls", "imports", "-w", file)
	}
	ExecWithLoading("Installing dependencies", "go", "mod", "tidy")

	if err := genProjectManifest(config); err != nil {
		return err
	}

	return nil
}
