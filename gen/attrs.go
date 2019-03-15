package gen

import (
  "bufio"
  "io"
  "encoding/json"
  "os"
  "strings"
)

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
    return "", err
  }
  return strings.TrimSpace(b.String()), nil
}

func ReadTemplateAttributesFromReader(templ io.Reader) (interface{}, error) {
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

func ReadTemplateAttributesFromString(templ string) (interface{}, error) {
  return ReadTemplateAttributesFromReader(strings.NewReader(templ))
}

func ReadTemplateAttributesFromPath(tplpath string) (interface{}, error) {
  f, err := os.Open(tplpath)
  if err != nil {
    return nil, err
  }
  defer f.Close()
  return ReadTemplateAttributesFromReader(f)
}
