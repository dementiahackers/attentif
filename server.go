package main

import (
	"flag"
	"net/http"

	"github.com/dementiahackers/attentif/internal/auth"
	"github.com/dementiahackers/attentif/internal/db"
	"github.com/dementiahackers/attentif/internal/templates"
	"github.com/rs/xhandler"
)

var (
	domain = flag.String("domain", "http://localhost", "Site domain")
	port   = flag.String("port", "8080", "Server port")
	client = flag.String("facebook-id", "1522378854752474", "Facebook Client ID")
	secret = flag.String("facebook-secret", "ec04d5b98a928fbd02df51574e4d48dd", "Facebook Client Secret")
)

func init() {
	flag.Parse()
}

func main() {
	auth.Config(*domain, *port, *client, *secret)
	tpl := templates.New("templates")

	// chain authenticated middleware
	c := xhandler.Chain{}
	c.UseC(func(next xhandler.HandlerC) xhandler.HandlerC {
		return auth.NewMiddleware(next)
	})

	// server static assets files
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		token, err := auth.GetToken(code)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		} else {
			id, err := db.CreateUser(token)
			if err != nil {
				tpl.Error(w, err)
			} else {
				auth.SaveSession(w, id)
				http.Redirect(w, r, "/entries/new", http.StatusFound)
			}
		}
	})

	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		auth.DestroySession(w)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			if r.URL.Path != "/" {
				tpl.NotFound(w)
				return
			}
			_, err := auth.CurrenUser(r)
			if err == nil {
				http.Redirect(w, r, "/entries/new", 302)
			} else {
				p := struct {
					FacebookURL string
				}{
					auth.RedirectURL(),
				}
				tpl.Render(w, "index", p)
			}
		} else {
			http.Error(w, "", http.StatusMethodNotAllowed)
		}
	})

	http.ListenAndServe(":"+*port, nil)
}
