// +build discard

package main

import (
	"html/template"
	"net/http"
)

//Home get /
func Home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("").ParseFiles("./template/index.html")
	err = tmpl.ExecuteTemplate(w, "base", "")
	checkErr(err)
}
