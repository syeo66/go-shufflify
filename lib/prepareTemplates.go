package lib

import "html/template"

func PrepareTemplates() map[string]*template.Template {
	tmpl := make(map[string]*template.Template)

	funcMap := template.FuncMap{
		"percent": func(current int, maximum int) float64 {
			if maximum == 0 {
				return 0
			}

			return float64(current) / float64(maximum) * 100.0
		},
	}

	tmpl["index.html"] = template.Must(template.New("index.html").Funcs(funcMap).ParseFiles(
		"templates/base.html",
		"templates/index.html",
		"templates/player.html",
		"templates/queue.html",
		"templates/tools.html",
	))

	tmpl["login.html"] = template.Must(template.New("login.html").Funcs(funcMap).ParseFiles(
		"templates/base.html",
		"templates/login.html",
	))

	tmpl["player.html"] = template.Must(template.New("player.html").Funcs(funcMap).ParseFiles("templates/player.html"))
	tmpl["queue.html"] = template.Must(template.New("queue.html").Funcs(funcMap).ParseFiles("templates/queue.html"))
	tmpl["tools.html"] = template.Must(template.New("tools.html").Funcs(funcMap).ParseFiles("templates/tools.html"))

	return tmpl
}
