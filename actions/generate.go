// command to bootstrap a new module
package actions

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Courtcircuits/optique/cli/templates"
	"github.com/Courtcircuits/optique/cli/views"
	"github.com/Courtcircuits/optique/core"
)

type ModuleType string

const (
	APPLICATION    ModuleType = "application"
	INFRASTRUCTURE ModuleType = "infrastructure"
)

func GenerateFromForm(name string) {
	view, err := views.LaunchGenForm()
	if err != nil {
		fmt.Println("Error launching form", err)
		fmt.Println("Error launching form")
		os.Exit(1)
	}
	var rtype ModuleType
	if view.Type == "application" {
		rtype = APPLICATION
	} else if view.Type == "infrastructure" {
		rtype = INFRASTRUCTURE
	} else {
		fmt.Println("Error generating module", err)
		fmt.Println("Error launching form")
		os.Exit(1)
	}
	if err := Generate(name, rtype, view.URL); err != nil {
		fmt.Println("Error generating module", err)
		fmt.Println("Error launching form")
		os.Exit(1)
	}

	if err := GoModInit(view.URL); err != nil {
		fmt.Println("Error generating module", err)
		fmt.Println("Error launching form")
		os.Exit(1)
	}
}

func Generate(name string, rtype ModuleType, url string) error {
	if err := CreateModuleFolder(name); err != nil { // goes into the module folder
		return err
	}

	if err := CreateRepositoryManifestFile(name, string(rtype), url); err != nil {
		return err
	}

	var code_opts *CodeGenerationOptions

	switch rtype {
	case APPLICATION:
		code_opts = CodeGenOpts(name, url, templates.APPLICATION_TPL)
	case INFRASTRUCTURE:
		code_opts = CodeGenOpts(name, url, templates.INFRASTRUCTURE_TPL)
	default:
		return fmt.Errorf("invalid module type: %s", rtype)
	}

	return GenerateCode(code_opts)
}

func CreateModuleFolder(name string) error {
	err := os.Mkdir(name, os.ModePerm)
	if err != nil {
		return err
	}
	return os.Chdir(name)
}


func CreateRepositoryManifestFile(name string, rtype string, url string) error {
	template_content := core.OptiqueModuleManifest{
		Name:   name,
		Type:   rtype,
		URL:    url,
		Ignore: []string{
			"go.mod",
			"go.sum",
		},
	}

	template, err := json.MarshalIndent(&template_content, "", "  ")
	if err != nil {
		return err
	}

	f, err := os.Create(core.MODULE_MANIFEST)
	defer f.Close()
	if err != nil {
		return err
	}

	_, err = f.Write(template)

	return err
}

func GoModInit(url string) error {
	return ExecWithLoading("Initializing go.mod", "go", "mod", "init", url)
}

