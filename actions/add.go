package actions

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Courtcircuits/optique/cli/manifests"
	"github.com/Courtcircuits/optique/cli/utils"
	"github.com/Courtcircuits/optique/core"
)

func AddModule(raw_url string) {
	// first go to root of the project
	root, err := core.FindOptiqueJson()
	if err != nil {
		fmt.Println("You are not in an optique project. To create a new one, run `optique init`")
	}
	os.Chdir(root)
	// parse the url
	repo_url, err := utils.ParseGitUrl(raw_url)
	if err != nil {
		panic(err)
	}
	fmt.Println(repo_url)
	data := SetUpSparseModule(repo_url.Repository, repo_url.Path)

	os.Chdir(repo_url.Path)
	if err := manifests.ClearIgnoredFiles(core.MODULE_MANIFEST); err != nil {
		panic(err)
	}
	goBack()
	CleanUpSparseModule()
	goBack()
	goBack()
	MoveModule(".optique/tmp/"+repo_url.Path, data.Type+"/"+data.Name)
	CleanUpOptique()
	ExecWithLoading("Installing dependencies", "go", "mod", "tidy")
}

func SetUpSparseModule(repo_url string, path string) *core.OptiqueModuleManifest {
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

func ParseModuleData() *core.OptiqueModuleManifest {
	// read config.json
	//get current directory
	current_dir, _ := os.Getwd()
	fmt.Println(current_dir)

	fd, err := os.Open(core.MODULE_MANIFEST)
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	var data core.OptiqueModuleManifest
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
