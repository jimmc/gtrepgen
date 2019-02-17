Data test
Field: {{data "field" .}}
{{with 123}}hello{{end}}
Row: {{with data "row" .}}
  a: {{.a}}
  b: {{.b}}
  c: {{.c}}
{{- end}}
Rows: {{range data "rows" .}}
  a: {{.a}}
  b: {{.b}}
  c: {{.c}}
{{- end}}
