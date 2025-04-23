package templates

const APPLICATION_TPL = `package {{.Name}}

type {{.NameCapitalized}} interface {
	Ignite() error
	Stop() error
	// Add more methods here
}

type {{.Name}} struct {}

func New{{.Name}}() (*{{.Name}}, error) {
  panic("implement me")
}

func (m *{{.Name}}) Ignite() error {
  panic("implement me")
}

func (m *{{.Name}}) Stop() error {
  panic("implement me")
}
`
