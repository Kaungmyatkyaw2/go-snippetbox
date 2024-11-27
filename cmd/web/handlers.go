package main

import (
	"eillion/snippetbox/pkg/forms"
	"eillion/snippetbox/pkg/models"
	"fmt"
	"net/http"
	"strconv"
)

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	s, err := app.snippets.Latest()

	if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{Snippets: s}

	app.render(w, r, "home.page.tmpl", data)

}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get(":id"))

	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Get(id)

	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{Snippet: s}

	app.render(w, r, "show.page.tmpl", data)

}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)

	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	if !form.IsValid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}

	title := form.Get("title")
	content := form.Get("content")
	expires := form.Get("expires")

	id, err := app.snippets.Insert(title, content, expires)

	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)

}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}
func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)

	form.Required("name", "email", "password")
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)

	if !form.IsValid() {
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	}

	err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))

	if err == models.ErrDuplicateEmail {
		form.Errors.Add("email", "Address is already in use")
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Your signup was successful. Please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {

	app.render(w, r, "login.page.tmpl", &templateData{
		Form: forms.New(nil),
	})

}
func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)

	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))

	if err == models.ErrInvalidCredentials {
		form.Errors.Add("generic", "Email or Password is incorrect")
		app.render(w, r, "login.page.tmpl", &templateData{Form: form})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "userID", id)

	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)

}
func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "userID")

	app.session.Put(r, "flash", "You've been logged out successfully!")

	http.Redirect(w, r, "/", 303)

}
