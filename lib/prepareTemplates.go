package lib

import "html/template"

func PrepareTemplates() map[string]*template.Template {
	tmpl := make(map[string]*template.Template)

	tmpl["index.html"] = template.Must(template.ParseFiles(
		"templates/base.html",
		"templates/index.html",
		"templates/player.html",
		"templates/queue.html",
		"templates/tools.html",
	))

	tmpl["login.html"] = template.Must(template.ParseFiles(
		"templates/base.html",
		"templates/login.html",
	))

	tmpl["player.html"] = template.Must(template.ParseFiles("templates/player.html"))
	tmpl["queue.html"] = template.Must(template.ParseFiles("templates/queue.html"))
	tmpl["tools.html"] = template.Must(template.ParseFiles("templates/tools.html"))

	return tmpl
}
