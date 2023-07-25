package client

import (
	"html/template"
)

var container *template.Template

func CreateContainer() {
	tmpl := `<!DOCTYPE html>
		<html>
		<head>
			<script src="https://unpkg.com/htmx.org@1.9.2" integrity="sha384-L6OqL9pRWyyFU3+/bjdSri+iIphTN/bvYyM37tICVyOJkWZLpP2vGn6VUEXgzg6h" crossorigin="anonymous"></script>
			<script src="https://unpkg.com/htmx.org/dist/ext/json-enc.js"></script>
			<link rel="stylesheet" type="text/css" href="/resources/style.css">
			<title>Beango Messenger</title>
		</head>
		<body>
			<div id="container">{{.}}</div>
		</body>
		</html>`

	t := template.Must(template.New("container").Parse(tmpl))
	container = t
}
