package utils

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

type GitUrl struct {
	Repository string
	Path       string
}

// clone a git repo from an https url that can include a subfolder from the repo
// used for example in `optique add github.com/user/repo/subfolder`
func ParseGitUrl(raw_url string) (*GitUrl, error) {
	splitted := strings.Split(raw_url, "/")
	repository := append([]string{"https://"}, splitted[:3]...)
	path := splitted[3:]
	if err := PingRepository(strings.Join(repository, "/")); err != nil {
		return nil, ErrRepositoryNotFound
	}
	return &GitUrl{
		Repository: strings.Join(repository, "/"),
		Path:       strings.Join(path, "/"),
	}, nil
}

func PingRepository(url string) (err error) {
	fmt.Println("pinging", url)
	cmd := exec.Command("git", "ls-remote", url)
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

var ErrRepositoryNotFound = errors.New("repository not found")
