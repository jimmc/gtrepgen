package gen

import (
  "io"
  "io/ioutil"
  "path"
  "text/template"

  "github.com/jimmc/gtrepgen/data"
)

const templateExtension = ".tpl"

// FromString executes the given literal template with the specified dot value.
func FromString(w io.Writer, source data.Source, name, templ string, dot interface{}) error {
  tpl := template.New(name)
  fm := template.FuncMap{
    "data": source.Data,
  }
  tpl = tpl.Funcs(fm)
  tpl, err := tpl.Parse(templ)
  if err != nil {
    return err
  }
  if err := tpl.Execute(w, dot); err != nil {
    return err
  }
  return nil
}

// FromPath reads a template from the given file path and executes it with the specified dot value.
func FromPath(w io.Writer, source data.Source, name, tplpath string, dot interface{}) error {
  templ, err := ioutil.ReadFile(tplpath)
  if err != nil {
    return err
  }
  return FromString(w, source, name, string(templ), dot)
}

// FromForm reads a template from a named file within a reference directory
// and executes it with the specified dot value.
func FromForm(w io.Writer, source data.Source, name, refpath string, dot interface{}) error {
  tplpath := path.Join(refpath, name) + templateExtension
  return FromPath(w, source, name, tplpath, dot)
}
