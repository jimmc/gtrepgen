Data test {{.}}
Row: {{with row .}}
  a: {{.a}}
  b: {{.b}}
  c: {{.c}}
{{- end}}
Rows: {{range rows .}}
  a: {{.a}}
  b: {{.b}}
  c: {{.c}}
{{- end}}
