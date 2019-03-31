package gen

/* This file contains support for embedding a JSON data structure inside a
 * go-template comment as the first item in the file. If the first line of
 * the file starts with our special comment prefix, and the comments ends
 * with our special suffix, then we assume the contents of that comment
 * are a JSON blob, which we read and parse. The calling application can
 * decide what the fields should be.
 */

import (
  "bufio"
  "fmt"
  "io"
  "io/ioutil"
  "encoding/json"
  "os"
  "path"
  "strings"

  "github.com/google/go-cmp/cmp"
)

/* TemplateAttributes holds the information for one file from a directory. */
type TemplateAttributes struct {
  Name string
  Attributes interface{}
  Err error
}

/* extractAttributeString reads and returns the portion of the input between
 * our special start and end strings.
 */
func extractAttributeString(templ io.Reader) (string, error) {
  magicStart := "{{/*GT:"      // This is what starts our attribute string.
  magicEnd := "*/ -}}"  // This is what ends our attribute string.
  scanner := bufio.NewScanner(templ)
  if !scanner.Scan() {
    // Empty input.
    return "", scanner.Err()
  }
  line := scanner.Text()
  if !strings.HasPrefix(line, magicStart) {
    // The content does not start with our magic prefix indicating our attributes.
    return "", nil
  }
  line = strings.TrimPrefix(line, magicStart)       // Strip the prefix.
  if strings.HasSuffix(line, magicEnd) {
    // Everything is on one line.
    line = strings.TrimSuffix(line, magicEnd)
    return strings.TrimSpace(line), nil
  }
  var b strings.Builder
  if _, err := b.WriteString(line); err != nil {
    return "", err
  }
  if _, err := b.WriteRune('\n'); err != nil {
    return "", err
  }
  for scanner.Scan() {
    line = scanner.Text()
    if strings.HasSuffix(line, magicEnd) {
      line = strings.TrimSuffix(line, magicEnd)
      if _, err := b.WriteString(line); err != nil {
        return "", err
      }
      if _, err := b.WriteRune('\n'); err != nil {
        return "", err
      }
      break     // We have our whole string, no need to scan any more.
    }
    if _, err := b.WriteString(line); err != nil {
      return "", err
    }
    if _, err := b.WriteRune('\n'); err != nil {
      return "", err
    }
  }
  if err := scanner.Err(); err != nil {
    return "", fmt.Errorf("scanning for template attributes: %v", err)
  }
  return strings.TrimSpace(b.String()), nil
}

/* ReadTemplateAttributesFromReader looks for our special start and end strings
 * in the given stream. If found, it parses the string between and returns
 * that value.
 */
func ReadTemplateAttributesFromReader(templ io.Reader) (interface{}, error) {
  var a interface{}
  err := ReadTemplateAttributesFromReaderInto(templ, &a)
  return a, err
}

/* ReadTemplateAttributesFromReaderInto looks for our special start and end strings
 * in the given stream. If found, it parses the string between into the given
 * pointer value.
 */
func ReadTemplateAttributesFromReaderInto(templ io.Reader, into interface{}) error {
  b, err := extractAttributeString(templ)
  if err != nil {
    return fmt.Errorf("extracting attributes from template: %v", err)
  }
  if b == "" {
    return nil
  }
  if err := json.Unmarshal([]byte(b), into); err != nil {
    return fmt.Errorf("unmarshalling json for template attributes: %v", err)
  }
  return nil
}

/* ReadTemplateAttributesFromString looks for our special start and end strings
 * in the given string. If found, it parses the string between and returns
 * that value.
 */
func ReadTemplateAttributesFromString(templ string) (interface{}, error) {
  return ReadTemplateAttributesFromReader(strings.NewReader(templ))
}

/* ReadTemplateAttributesFromStringInto looks for our special start and end strings
 * in the given string. If found, it parses the string between into the given pointer.
 */
func ReadTemplateAttributesFromStringInto(templ string, into interface{}) error {
  return ReadTemplateAttributesFromReaderInto(strings.NewReader(templ), into)
}

/* ReadTemplateAttributesFromPath looks for our special start and end strings
 * in the file at the specified path. If found, it parses the string between
 * and returns that value.
 */
func ReadTemplateAttributesFromPath(tplpath string) (interface{}, error) {
  var a interface{}
  err := ReadTemplateAttributesFromPathInto(tplpath, &a)
  return a, err
}

/* ReadTemplateAttributesFromPathInto looks for our special start and end strings
 * in the file at the specified path. If found, it parses the string between
 * into the given location.
 */
func ReadTemplateAttributesFromPathInto(tplpath string, into interface{}) error {
  f, err := os.Open(tplpath)
  if err != nil {
    return fmt.Errorf("opening template file %s: %v", tplpath, err)
  }
  defer f.Close()
  return ReadTemplateAttributesFromReaderInto(f, into)
}

/* ReadDirFilesAttributes calls ReadDirFilesAttributesAs with a newDestPointer
 * that creates an interface{}. This allows reading generic JSON data from each
 * file in the directory.
 */
func ReadDirFilesAttributes(tpldir string) ([]*TemplateAttributes, error) {
  newDestPointer := func() interface{} {
    var a interface{}
    return &a
  }
  return ReadDirFilesAttributesAs(tpldir, newDestPointer)
}

/* ReadDirFilesAttributesAs scans the given directory looking for files
 * with the right filename extension. For each found, it looks for
 * and parses the contents of our special attribute strings.
 * It returns an array of results holding that information.
 * If there are errors reading individual files, those errors
 * are returned in the array, along with an overall error
 * that tells how many files had errors.
 * newDestPointer is a function that is called just before reading
 * the contents of each file. It should return a new instance of
 * a location where the contents of that file should be stored.
 * Typically it will look like: func()interface{}{return &MyStruct{}}.
 */
func ReadDirFilesAttributesAs(tpldir string, newDestPointer func() interface{}) ([]*TemplateAttributes, error) {
  fileinfos, err := ioutil.ReadDir(tpldir)
  if err != nil {
    return nil, fmt.Errorf("reading templates from %s: %v", tpldir, err)
  }
  errCount := 0
  templateAttributes := []*TemplateAttributes{}
  zeroAttrs := newDestPointer()
  for _, fileinfo := range fileinfos {
    if fileinfo.IsDir() {
      continue  // Ignore directories.
    }
    fname := fileinfo.Name()
    if !strings.HasSuffix(fname, templateExtension) {
      continue  // Ignore files without the correct filename extension.
    }
    name := strings.TrimSuffix(fname, templateExtension)
    filepath := path.Join(tpldir, fname)
    attrs := newDestPointer()
    err := ReadTemplateAttributesFromPathInto(filepath, attrs)
    if !cmp.Equal(attrs, zeroAttrs) || err != nil {
      tplAttrs := &TemplateAttributes{
        Name: name,
        Attributes: attrs,
        Err: err,
      }
      templateAttributes = append(templateAttributes, tplAttrs)
    }
    if err != nil {
      errCount++
    }
  }
  if errCount > 0 {
    s := ""
    if errCount > 1 { s = "s" }
    err = fmt.Errorf("%d error%s reading template attributes, see err fields", errCount, s)
    return templateAttributes, err
  }
  return templateAttributes, nil
}
