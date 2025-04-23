package templates

const INFRASTRUCTURE_TPL = `package {{.Name}}

type {{.NameCapitalized}} interface {
	Setup() error
	Shutdown() error
	// Add more methods here
}

type {{.Name}} struct {}

func New{{.Name}}() (*{{.Name}}, error) {
  panic("implement me")
}

func (m *{{.Name}}) Setup() error {
  panic("implement me")
}

func (m *{{.Name}}) Shutdown() error {
  panic("implement me")
}
`
