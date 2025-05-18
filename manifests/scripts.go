package manifests

import (
	"bytes"
	"os"
	"text/template"
)

const SCRIPT_TPL = `{{.Key}}:
	{{.Script}}
`

func GenJustScript(key string, value string) (string, error) {
	script := ""
	type ScriptToPass struct {
		Key    string
		Script string
	}
	tmpl, err := template.New("script").Parse(SCRIPT_TPL)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, ScriptToPass{key, value})
	if err != nil {
		return "", err
	}
	script = buf.String()

	return script, nil
}

func GenScripts(scripts map[string]string) (string, error) {
	just_script := ""
	for key, script := range scripts {
		s, err := GenJustScript(key, script)
		if err != nil {
			return "", err
		}
		just_script += s + "\n"
	}
	return just_script, nil
}

func SaveScripts(script string, just_file_path string) error {
	script_as_bytes := []byte(script)
	f, err := os.OpenFile(just_file_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(script_as_bytes); err != nil {
		return err
	}
	return nil
}
