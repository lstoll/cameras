package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/secure"
	"github.com/martini-contrib/sessionauth"
	"github.com/martini-contrib/sessions"
)

func main() {
	m := martini.Classic()

	m.Use(render.Renderer(render.Options{
		Layout: "_layout",
	}))

	// TODO - secret should be secret.
	m.Use(secure.Secure(secure.Options{
		SSLRedirect:     true,
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	}))
	store := sessions.NewCookieStore([]byte(os.Getenv("COOKIE_SECRET")))
	m.Use(sessions.Sessions("the_session", store))
	m.Use(sessionauth.SessionUser(GenerateAnonymousUser))
	m.Run()

	/** Main router **/

	m.Get("/", sessionauth.LoginRequired, cameraList)
	m.Get("/camimage", sessionauth.LoginRequired, cameraImage)

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

func getCameras() map[string]string {
	ret := map[string]string{}
	// env var, comma separates list of http://u:p@host/img? Maybe name;url?
	// FORMAT: CAMERAS="Downstairs;http://username:password@host/image.jpg,Inside;http://..."
	cams := strings.Split(os.Getenv("CAMERAS"), ",")
	for _, i := range cams {
		f := strings.Split(i, ";")
		ret[f[0]] = f[1]
	}
	return ret
}

func cameraImage(res http.ResponseWriter, req *http.Request) {
	// get name from query param
	camName := req.URL.Query().Get("cam")
	if camName == "" {
		res.WriteHeader(400)
		fmt.Fprint(res, "need to specify cam")
		return
	}

	camUrl, ok := getCameras()[camName]
	if !ok {
		res.WriteHeader(400)
		fmt.Fprint(res, "Invalid cam name")
		return
	}

	parsedUrl, err := url.Parse(camUrl)
	if err != nil {
		res.WriteHeader(500)
		fmt.Fprintf(res, "%s", err)
		return
	}

	imgReq, err := http.NewRequest("GET", camUrl, nil)
	password, _ := parsedUrl.User.Password()
	imgReq.SetBasicAuth(parsedUrl.User.Username(), password)
	cli := &http.Client{}
	img, err := cli.Do(imgReq)

	if err != nil {
		res.WriteHeader(500)
		fmt.Fprintf(res, "%s", err)
		return
	}

	defer func() { _ = img.Body.Close() }()

	if img.StatusCode != 200 {
		fmt.Fprintf(res, "Remote camera returned non-200 response %d", img.StatusCode)
		return
	}

	res.Header().Set("Content-Length", fmt.Sprint(img.ContentLength))
	res.Header().Set("Content-Type", img.Header.Get("Content-Type"))
	if _, err = io.Copy(res, img.Body); err != nil {
		return
	}

}

func cameraList(r render.Render) {
	// Get list of camera's out of env var

	// render list of them

	// auto refresh in js.
	r.HTML(200, "index", getCameras())
}
