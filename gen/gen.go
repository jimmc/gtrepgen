package gen

import (
  "fmt"
  htmltemplate "html/template"
  "io"
  "io/ioutil"
  "os"
  "path"
  texttemplate "text/template"
  "time"

  "github.com/golang/glog"

  "github.com/jimmc/gtrepgen/data"
)

const templateExtension = ".tpl"

// Now returns the current time. You can provide your own Now function for testing.
var Now = time.Now

// Generator provides a type that can generate output from either text or HTML templates
// from a literal string, a specific file path, or a named file from a reference directory.
type Generator struct {
  name string
  w io.Writer
  source data.Source
  isHTML bool
  refpaths []string
  funcs map[string]interface{}
  includeResult interface{}
}

// New creates a Generator.
func New(name string, isHTML bool, w io.Writer, source data.Source) *Generator {
  glog.V(1).Infof("gtrepgen.New(%s)", name)
  return &Generator{
    name: name,
    w: w,
    source: source,
    isHTML: isHTML,
  }
}

// Create a copy of a generator with a changed name.
func (g *Generator) WithName(name string) *Generator {
  glog.V(1).Infof("gtrepgen.WithName(%s) from name %s", name, g.name)
  return &Generator{
    name: name,
    w: g.w,
    source: g.source,
    isHTML: g.isHTML,
    refpaths: g.refpaths,
    funcs: g.funcs,
  }
}

// Create a copy of a generator with a changed refpaths.
func (g *Generator) WithRefpaths(refpaths []string) *Generator {
  glog.V(1).Infof("gtrepgen.WithRefpaths(%v) from name %s", refpaths, g.name)
  return &Generator{
    name: g.name,
    w: g.w,
    source: g.source,
    isHTML: g.isHTML,
    refpaths: refpaths,
    funcs: g.funcs,
  }
}

// Create a copy of a generator with a changed funcs.
func (g *Generator) WithFuncs(funcs map[string]interface{}) *Generator {
  glog.V(1).Infof("gtrepgen.WithFuncs() from name %s", g.name)
  return &Generator{
    name: g.name,
    w: g.w,
    source: g.source,
    isHTML: g.isHTML,
    refpaths: g.refpaths,
    funcs: funcs,
  }
}

// include allows us to include another template from our reference directory.
// Args is either no args or a single arg that sets dot.
func (g *Generator) include(name string, args ...interface{}) (interface{}, error) {
  glog.V(2).Infof("gtrepgen.include(%s)", name)
  tplpath, err := g.FindTemplate(name)
  if err != nil {
    return nil, fmt.Errorf("template %s: %v", name, err)
  }
  var dot interface{}
  if len(args) > 1 {
    return nil, fmt.Errorf("too many args (%d) for Generator.template", len(args))
  } else if len(args) == 0 {
    dot = nil
  } else {
    dot = args[0]
  }
  gInclude := g.WithName(name)
  gInclude.includeResult = ""
  if err := gInclude.FromPath(tplpath, dot); err != nil {
    return nil, err
  }
  return gInclude.includeResult, nil
}

// evalTemplate allows us to evaluate a template from a given string, as if we read it from a file.
// Args is either no args or a single arg that sets dot.
func (g *Generator) evalTemplate(template string, args ...interface{}) (interface{}, error) {
  glog.V(2).Infof("gtrepgen.evalTemplate()")
  var dot interface{}
  if len(args) > 1 {
    return nil, fmt.Errorf("too many args (%d) for Generator.template", len(args))
  } else if len(args) == 0 {
    dot = nil
  } else {
    dot = args[0]
  }
  gInclude := g.WithName("evalTemplate")
  gInclude.includeResult = ""
  if err := gInclude.FromString(template, dot); err != nil {
    return nil, err
  }
  return gInclude.includeResult, nil
}

// includeReturn provides a way for an included template to pass a return value back to
// the including template. If the included template does not invoke return, then the
// return value of the include statement is the empty string.
// (A return value of nil causes the string "<no value>" to appear in the output when the
// include statement is called without assigning the output in the caller.)
// The value of the return expression itself is the empty string.
func (g *Generator) includeReturn(returnVal interface{}) (interface{}, error) {
  g.includeResult = returnVal
  return "", nil
}

// htmlFromString executes the given literal template with the specified dot value
// using html/template.
func (g *Generator) htmlFromString(templ string, dot interface{}, fm map[string]interface{}) error {
  glog.V(2).Infof("gtrepgen.htmlFromString()")
  tpl := htmltemplate.New(g.name)
  tpl = tpl.Funcs(fm)
  if g.funcs != nil {
    tpl = tpl.Funcs(g.funcs)
  }
  tpl, err := tpl.Parse(templ)
  if err != nil {
    return fmt.Errorf("parsing html template %s: %v", g.name, err)
  }
  if err := tpl.Execute(g.w, dot); err != nil {
    return fmt.Errorf("executing html template %s: %v", g.name, err)
  }
  return nil
}

// textFromString executes the given literal template with the specified dot value
// using text/template.
func (g *Generator) textFromString(templ string, dot interface{},fm map[string]interface{}) error {
  glog.V(2).Infof("gtrepgen.textFromString()")
  tpl := texttemplate.New(g.name)
  tpl = tpl.Funcs(fm)
  if g.funcs != nil {
    tpl = tpl.Funcs(g.funcs)
  }
  tpl, err := tpl.Parse(templ)
  if err != nil {
    return fmt.Errorf("parsing text template %s: %v", g.name, err)
  }
  if err := tpl.Execute(g.w, dot); err != nil {
    return fmt.Errorf("executing text template %s: %v", g.name, err)
  }
  return nil
}

// FromString executes the given literal template with the specified dot value.
func (g *Generator) FromString(templ string, dot interface{}) error {
  now := Now()
  startTime := func() time.Time { return now }
  fm := map[string]interface{}{  // fm is a (texttemplate|htmltemplate).FuncMap
    "include": g.include,
    "evalTemplate": g.evalTemplate,
    "evenodd": evenodd,
    "formatTime": formatTime,
    "mkmap": mkmap,
    "reportStartTime": startTime,
    "return": g.includeReturn,
    "row": g.source.Row,
    "rows": g.source.Rows,
  }
  if g.isHTML {
    return g.htmlFromString(templ, dot, fm)
  } else {
    return g.textFromString(templ, dot, fm)
  }
}

// FromPath reads a template from the given file path and executes it with the specified dot value.
func (g *Generator) FromPath(tplpath string, dot interface{}) error {
  templ, err := ioutil.ReadFile(tplpath)
  if err != nil {
    return fmt.Errorf("reading template file %s: %v", tplpath, err)
  }
  return g.FromString(string(templ), dot)
}

// FromTemplate reads a template from a named file within a set of reference directories
// and executes it with the specified dot value.
func (g *Generator) FromTemplate(refpaths []string, dot interface{}) error {
  g = g.WithRefpaths(refpaths)
  tplpath, err := g.FindTemplate(g.name)
  if err != nil {
    return err
  }
  return g.FromPath(tplpath, dot)
}

// FindTemplate finds the first readable template in the list of reference directories.
func (g *Generator) FindTemplate(name string) (string, error) {
  return FindTemplateInDirs(name, g.refpaths)
}

// FindAndReadAttributes finds the template and reads the attributes from it.
func FindAndReadAttributes(name string, dirs []string) (interface{}, error) {
  tplpath, err := FindTemplateInDirs(name, dirs)
  if err != nil {
    return nil, err
  }
  return ReadTemplateAttributesFromPath(tplpath)
}

// FindAndReadAttributesInto finds the template and reads the attributes from it
// into the specified destination.
func FindAndReadAttributesInto(name string, dirs []string, dest interface{}) error {
  tplpath, err := FindTemplateInDirs(name, dirs)
  if err != nil {
    return err
  }
  return ReadTemplateAttributesFromPathInto(tplpath, dest)
}

// FindTemplateInDirs finds the first readable template in the given list of directories.
func FindTemplateInDirs(name string, refpaths []string) (string, error) {
  for _, d := range refpaths {
    tplpath := path.Join(d, name) + templateExtension
    f, err := os.Open(tplpath)
    if err == nil {
      f.Close()
      return tplpath, nil
    }
  }
  return "", fmt.Errorf("template for %q not found", name)
}

// mkmap create a map using pairs of keys and values.
func mkmap(args ...interface{}) (map[interface{}]interface{}, error) {
  if len(args)%2 != 0 {
    return nil, fmt.Errorf("mkmap: args count must be even (count=%d)", len(args))
  }
  m := make(map[interface{}]interface{}, len(args)/2)
  for k := 0; k < len(args); k += 2 {
    key := args[k]
    val := args[k+1]
    m[key] = val
  }
  return m, nil
}

// evenodd returns the second or third arg based on whether the first arg is even or odd.
func evenodd(n int, evenret, oddret interface{}) interface{} {
  if n % 2 == 0 {
    return evenret
  } else {
    return oddret
  }
}

// formatTime passes the args through to time.Format.
func formatTime(format string, t time.Time) string {
  return t.Format(format)
}
