Report for companyid {{.}}
Company data:
{{range data "select name from company where id = ?" .}}
    Name: {{.name}}
{{end}}
People in that company:
{{range data "select firstname, lastname from person where companyid = ?" .}}
    First name: {{.firstname}}
    Last name:  {{.lastname}}
{{end}}
