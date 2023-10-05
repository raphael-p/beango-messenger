package client

import "html/template"

var skeleton *template.Template

// TODO: use local htmx.min & json-enc instead of having to download from unpkg
func GetSkeleton() (*template.Template, error) {
	if skeleton != nil {
		return skeleton, nil
	}
	tmpl := `<!DOCTYPE html>
		<html>
		<head>
			<script src="https://unpkg.com/htmx.org@1.9.6"></script>
			<script src="https://unpkg.com/htmx.org/dist/ext/json-enc.js"></script>
			<script src="/resources/script.js"></script>
			<link rel="stylesheet" type="text/css" href="/resources/style.css">
			<title>Beango Messenger</title>
		</head>
		<body>
			<div id="header">{{.header}}</div>
			<div id="content">{{.content}}</div>
			<div id="footer">{{.footer}}</div>
		</body>
		</html>`

	t, err := template.New("skeleton").Parse(tmpl)
	if err != nil {
		return nil, err
	}
	skeleton = t
	return skeleton, nil
}
