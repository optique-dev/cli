package views

import (
	"errors"

	"github.com/charmbracelet/huh"
)


type GenForm struct {
	Type string
	URL string
}

func (i *GenForm) CreateForm() *huh.Form {
	form :=
		huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Golang module name").
					Value(&i.URL).
					Validate(func(str string) error {
						if len(str) == 0 {
							return errors.New("Module name cannot be empty")
						}
						return nil
					}),
				huh.NewSelect[string]().
					Title("Version").
					Options(
						huh.NewOption("Application", "application"),
						huh.NewOption("Infrastructure", "infrastructure"),
					).
					Value(&i.Type),
			))
	return form
}

func LaunchGenForm() (*GenForm, error) {
	genForm := GenForm{}
	form := genForm.CreateForm()
	if err := form.Run(); err != nil {
		return nil, err
	}
	return &genForm, nil
}
