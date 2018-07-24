package main

import (
	"html/template"
	"net/http"
)

//Home get /
func Home(w http.ResponseWriter, r *http.Request) {
	if DisallowMethod(w, "GET", r.Method) {
		return
	}
	tmpl, err := template.New("").ParseFiles("./template/index.html")
	err = tmpl.ExecuteTemplate(w, "base", "")
	checkErr(err)
}
