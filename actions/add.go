package actions

import (
	"encoding/json"
	"os"

	"github.com/optique-dev/cli/manifests"
	"github.com/optique-dev/cli/utils"
	"github.com/optique-dev/optique"
)

func AddModule(raw_url string) {
	// first go to root of the project
	root, err := optique.FindOptiqueJson()
	if err != nil {
		optique.Error("You are not in an optique project. To create a new one, run `optique init`")
	}
	os.Chdir(root)
	// parse the url
	repo_url, err := utils.ParseGitUrl(raw_url)
	if err != nil {
		panic(err)
	}
	data := SetUpSparseModule(repo_url.Repository, repo_url.Path)

	os.Chdir(repo_url.Path)
	if err := manifests.ClearIgnoredFiles(optique.MODULE_MANIFEST); err != nil {
		panic(err)
	}
	goBack()
	CleanUpSparseModule()
	goBack()
	goBack()
	MoveModule(".optique/tmp/"+repo_url.Path, data.Type+"/"+data.Name)
	CleanUpOptique()
	ExecWithLoading("Installing dependencies", "go", "mod", "tidy")
	if data.Scripts == nil {
		return
	}
	just_scripts, err := manifests.GenScripts(data.Scripts)
	if err != nil {
		panic(err)
	}
	err = manifests.SaveScripts(just_scripts, "./justfile")
	if err != nil {
		panic(err)
	}
}

func SetUpSparseModule(repo_url string, path string) *optique.OptiqueModuleManifest {
	//create temp folder
	os.Mkdir(".optique", os.ModePerm)
	os.Chdir(".optique")
	os.Mkdir("tmp", os.ModePerm)
	os.Chdir("tmp")
	ExecWithLoading("Initializing module", "git", "init")
	ExecWithLoading("Creating sparse module", "git", "remote", "add", "origin", repo_url)
	ExecWithLoading("Sparse-checking module", "git", "sparse-checkout", "init", "--cone")
	ExecWithLoading("Setting up module", "git", "sparse-checkout", "set", path)
	ExecWithLoading("Pulling module", "git", "pull", "origin", "main")
	current_dir, _ := os.Getwd()
	os.Chdir(path)
	defer os.Chdir(current_dir)
	return ParseModuleData()
}

func CleanUpSparseModule() {
	os.RemoveAll(".git")
}

func ParseModuleData() *optique.OptiqueModuleManifest {
	// read config.json

	fd, err := os.Open(optique.MODULE_MANIFEST)
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	var data optique.OptiqueModuleManifest
	err = json.NewDecoder(fd).Decode(&data)
	if err != nil {
		panic(err)
	}

	return &data
}

// move module from tmp location to destination
func MoveModule(path string, destination string) {
	ExecWithLoading("Moving module", "mv", path, destination)
}

func CleanUpOptique() {
	os.RemoveAll(".optique")
}
