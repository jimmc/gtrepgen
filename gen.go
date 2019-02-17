package gen

import (
  "io"
  "io/ioutil"
  "path"
  "text/template"
)

const templateExtension = ".tpl"

// FromString executes the given literal template with the specified dot value.
func FromString(w io.Writer, name, templ string, dot interface{}) error {
  tpl, err := template.New(name).Parse(templ)
  if err != nil {
    return err
  }
  if err := tpl.Execute(w, dot); err != nil {
    return err
  }
  return nil
}

// FromPath reads a template from the given file path and executes it with the specified dot value.
func FromPath(w io.Writer, name, tplpath string, dot interface{}) error {
  templ, err := ioutil.ReadFile(tplpath)
  if err != nil {
    return err
  }
  return FromString(w, name, string(templ), dot)
}

// FromForm reads a template from a named file within a reference directory
// and executes it with the specified dot value.
func FromForm(w io.Writer, name, refpath string, dot interface{}) error {
  tplpath := path.Join(refpath, name) + templateExtension
  return FromPath(w, name, tplpath, dot)
}
