Report for companyid {{.}}
Company data:
{{with row "select name from company where id = ?" .}}
    Name: {{.name}}
{{end}}
People in that company:
{{range rows "select firstname, lastname from person where companyid = ?" .}}
    First name: {{.firstname}}
    Last name:  {{.lastname}}
{{end}}
