package actions

import (
	"context"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"os/signal"

	"maps"

	"github.com/charmbracelet/huh/spinner"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func ExecWithLoading(label string, name string, commands ...string) error {
	cmd := context.Background()
	ctx, cancel := context.WithCancel(cmd)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()
	go func() {
		if err := spinner.New().Title(label).Context(ctx).Run(); err != nil {
		}
		select {
		case <-c:
			cancel()
			panic("Err")
		case <-ctx.Done():
		}
	}()
	output, err := exec.CommandContext(ctx, name, commands...).CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		return err
	}
	return nil
}

type CodeGenerationOptions struct {
	Path              string
	Name              string
	URL               string
	Template          string
	AdditionalOptions map[string]any
}

func CodeGenOpts(name, url, template string) *CodeGenerationOptions {
	return &CodeGenerationOptions{
		Name:              name,
		URL:               url,
		Template:          template,
		AdditionalOptions: map[string]any{},
	}
}

func WithPath(path string) *CodeGenerationOptions {
	return &CodeGenerationOptions{
		Path: path,
	}
}

func WithAdditionalOptions(options map[string]any) *CodeGenerationOptions {
	return &CodeGenerationOptions{
		AdditionalOptions: options,
	}
}

func GenerateCode(options *CodeGenerationOptions) error {
	if options.Path == "" {
		options.Path = options.Name + ".go"
	}

	codegen_params := map[string]any{
		"Name":            options.Name,
		"URL":             options.URL,
		"NameCapitalized": cases.Title(language.English).String(options.Name),
	}
	maps.Copy(codegen_params, options.AdditionalOptions)

	fmt.Println("Generating code...", codegen_params)
	fmt.Println("Template:", options.Template)

	tmpl, err := template.New(options.Template).Parse(options.Template)
	if err != nil {
		return err
	}

	f, err := os.Create(options.Path)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, codegen_params)
}
