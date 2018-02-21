package config

import "html/template"

// Templ hold our template collection, and init only once
var Templ *template.Template

func init() {
	Templ = template.Must(template.ParseGlob("templates/*.html"))
}
