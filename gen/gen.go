package gen

import (
  "fmt"
  htmltemplate "html/template"
  "io"
  "io/ioutil"
  "path"
  texttemplate "text/template"

  "github.com/jimmc/gtrepgen/data"
)

const templateExtension = ".tpl"

// Generator provides a type that can generate output from either text or HTML templates
// from a literal string, a specific file path, or a named file from a reference directory.
type Generator struct {
  name string;
  w io.Writer;
  source data.Source;
  isHTML bool;
  refpath string;
}

// New creates a Generator.
func New(name string, isHTML bool, w io.Writer, source data.Source) *Generator {
  return &Generator{
    name: name,
    w: w,
    source: source,
    isHTML: isHTML,
  }
}

// include allows us to include another template from our reference directory.
// Args is either no args or a single arg that sets dot.
func (g *Generator) include(name string, args ...interface{}) (interface{}, error) {
  tplpath := path.Join(g.refpath, name) + templateExtension
  var dot interface{}
  if len(args) > 1 {
    return nil, fmt.Errorf("Too many args (%d) for Generator.template", len(args))
  } else if len(args) == 0 {
    dot = nil
  } else {
    dot = args[0]
  }
  return "", g.FromPath(tplpath, dot)
}

// htmlFromString executes the given literal template with the specified dot value
// using html/template.
func (g *Generator) htmlFromString(templ string, dot interface{}) error {
  tpl := htmltemplate.New(g.name)
  fm := htmltemplate.FuncMap{
    "row": g.source.Row,
    "rows": g.source.Rows,
    "include": g.include,
  }
  tpl = tpl.Funcs(fm)
  tpl, err := tpl.Parse(templ)
  if err != nil {
    return err
  }
  if err := tpl.Execute(g.w, dot); err != nil {
    return err
  }
  return nil
}

// textFromString executes the given literal template with the specified dot value
// using text/template.
func (g *Generator) textFromString(templ string, dot interface{}) error {
  tpl := texttemplate.New(g.name)
  fm := texttemplate.FuncMap{
    "row": g.source.Row,
    "rows": g.source.Rows,
    "include": g.include,
  }
  tpl = tpl.Funcs(fm)
  tpl, err := tpl.Parse(templ)
  if err != nil {
    return err
  }
  if err := tpl.Execute(g.w, dot); err != nil {
    return err
  }
  return nil
}

// FromString executes the given literal template with the specified dot value.
func (g *Generator) FromString(templ string, dot interface{}) error {
  if g.isHTML {
    return g.htmlFromString(templ, dot)
  } else {
    return g.textFromString(templ, dot)
  }
}

// FromPath reads a template from the given file path and executes it with the specified dot value.
func (g *Generator) FromPath(tplpath string, dot interface{}) error {
  templ, err := ioutil.ReadFile(tplpath)
  if err != nil {
    return err
  }
  return g.FromString(string(templ), dot)
}

// FromForm reads a template from a named file within a reference directory
// and executes it with the specified dot value.
func (g *Generator) FromForm(refpath string, dot interface{}) error {
  g.refpath = refpath
  tplpath := path.Join(refpath, g.name) + templateExtension
  return g.FromPath(tplpath, dot)
}
