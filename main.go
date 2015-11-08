package main

import (
	"net/http"
	"os"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessionauth"
	"github.com/martini-contrib/sessions"
)

func main() {
	m := martini.Classic()

	m.Use(render.Renderer(render.Options{
		Layout: "_layout",
	}))

	// TODO - secret should be secret.
	store := sessions.NewCookieStore([]byte("secret123"))
	m.Use(sessions.Sessions("the_session", store))
	m.Use(sessionauth.SessionUser(GenerateAnonymousUser))

	m.Get("/", sessionauth.LoginRequired, func(r render.Render) {
		r.HTML(200, "index", "lol")
	})

	/** Login Handling **/

	m.Get("/login", func(r render.Render) {
		r.HTML(200, "login", nil, render.HTMLOptions{
			Layout: "_login_layout",
		})
	})

	m.Post("/login", binding.Bind(User{}), func(session sessions.Session, postedUser User, r render.Render, req *http.Request) {

		// if not logged in
		if postedUser.Passcode != "" && postedUser.Passcode == os.Getenv("WEB_PASSCODE") {
			user := &User{}
			err := sessionauth.AuthenticateSession(session, user)
			if err != nil {
				r.Text(500, "Error authenticating session")
				return
			}

			params := req.URL.Query()
			redirect := params.Get(sessionauth.RedirectParam)
			r.Redirect(redirect)
			return
		} else {
			r.Redirect(sessionauth.RedirectUrl)
			return

		}
	})

	m.Get("/logout", sessionauth.LoginRequired, func(session sessions.Session, user sessionauth.User, r render.Render) {
		sessionauth.Logout(session, user)
		r.Redirect("/")
	})

	m.Run()
}
