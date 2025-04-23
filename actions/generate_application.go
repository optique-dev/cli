package actions

/**
* Generate application module
*
* This module will generate a new application module
*
*
* @param name name of the module
* @param url url of the module
* @param type type of the module
*
* The generated code MUST implement the Application interface
* - Ignite() error
* - Stop() error
 */

const APPLICATION_TPL = `
package {{.URL}}

type {{.NameCapitalized}} struct {}

func New{{.NameCapitalized}}() (*{{.NameCapitalized}}, error) {
  panic("implement me")
}

func (m *{{.NameCapitalized}}) Ignite() error {
  panic("implement me")
}

func (m *{{.NameCapitalized}}) Stop() error {
  panic("implement me")
}
`
