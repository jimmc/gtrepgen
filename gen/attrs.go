package gen

import (
  "bufio"
  "io/ioutil"
  "encoding/json"
  "strings"
)

func extractAttributeString(templ string) (string, error) {
  magic := "{{/*GT:"      // This is what starts our attribute string.
  terminator := "*/ -}}"  // This is what ends our attribute string.
  scanner := bufio.NewScanner(strings.NewReader(templ))
  if !scanner.Scan() {
    // Empty input.
    return "", scanner.Err()
  }
  line := scanner.Text()
  if !strings.HasPrefix(line, magic) {
    // The content does not start with our magic prefix indicating our attributes.
    return "", nil
  }
  line = strings.TrimPrefix(line, magic)       // Strip the prefix.
  if strings.HasSuffix(line, terminator) {
    // Everything is on one line.
    line = strings.TrimSuffix(line, terminator)
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
    if strings.HasSuffix(line, terminator) {
      line = strings.TrimSuffix(line, terminator)
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
    return "", err
  }
  return strings.TrimSpace(b.String()), nil
}

func ReadTemplateAttributesFromString(templ string) (interface{}, error) {
  b, err := extractAttributeString(templ)
  if err != nil {
    return nil, err
  }
  if b == "" {
    return nil, nil
  }
  var a interface{}
  if err := json.Unmarshal([]byte(b), &a); err != nil {
    return nil, err
  }
  return a, nil
}

func ReadTemplateAttributesFromPath(tplpath string) (interface{}, error) {
  templ, err := ioutil.ReadFile(tplpath)
  if err != nil {
    return nil, err
  }
  return ReadTemplateAttributesFromString(string(templ))
}
