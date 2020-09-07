// Author Sean Joyce <seanjoyce.012@gmail.com>
// Data 09/2020

// cls && go build cmd/main.go && main.exe
// cls && del main.exe && go build main.go && main.exe

package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"time"
)


type Page struct {
	Title    string
	Body     string
	NumberOne     int
	NumberTwo     int
}

var templates = template.Must(template.ParseFiles(
	"web/template/view.html",
	))
var validPath = regexp.MustCompile("^/(view)/([a-zA-Z0-9]+)$")

func loadPage(title string) (*Page, error) {
	return &Page{Title: title, NumberOne: 15, NumberTwo: 4}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
			http.Redirect(w, r, "/view/"+title, http.StatusFound)
			return
	}
	renderTemplate(w, "view", p)
}

func makeHandler(fn func(http.ResponseWriter,
												*http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
			m := validPath.FindStringSubmatch(r.URL.Path)
			if m == nil {
					http.NotFound(w, r)
					return
			}
			fn(w, r, m[2])
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

const (
	address     = "localhost:8080"
	defaultName = "world"
)

func main() {
	fmt.Println("Please Login:")
	fmt.Println("\n\nhttp://localhost:8080/view/index")

	http.HandleFunc("/view/", makeHandler(viewHandler))

  fs := http.FileServer(http.Dir("./web/static"))
  http.Handle("/web/static/", http.StripPrefix("/web/static/", fs))

	s := &http.Server{
		Addr:           ":8080",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())

}
