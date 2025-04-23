package views

import (
	"errors"

	"github.com/charmbracelet/huh"
)


type InitForm struct {
	Repository string
	Version    string
}

func (i *InitForm) CreateFormInit() *huh.Form {
	form :=
		huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Golang module name").
					Value(&i.Repository).
					Validate(func(str string) error {
						if len(str) == 0 {
							return errors.New("Module name cannot be empty")
						}
						return nil
					}),
				huh.NewSelect[string]().
					Title("Version").
					Options(
						huh.NewOption("latest", "latest"),
					).
					Value(&i.Version),
			))
	return form
}

func LaunchInitForm() (*InitForm, error) {
	initForm := InitForm{}
	form := initForm.CreateFormInit()
	if err := form.Run(); err != nil {
		return nil, err
	}
	return &initForm, nil
}
