package app

import (
	"net/http"

	"github.com/gorilla/mux"
)

type tplVars map[string]interface{}

type RequestContext struct {
	Writer  http.ResponseWriter
	Request *http.Request
}

func (c *RequestContext) HTML(status int, tpl string, vars tplVars) error {
	if vars == nil {
		vars = tplVars{}
	}
	flashes, err := c.flashes()
	if err != nil {
		return err
	}
	vars["Flashes"] = flashes
	return renderer.HTML(c.Writer, status, tpl, vars)
}

func (c *RequestContext) RenderInternalServerErrorHTML(err error) {
	http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
}

func (c *RequestContext) flashes() ([]interface{}, error) {
	session, err := jar.Get(c.Request, "ses")
	if err != nil {
		return nil, err
	}
	flashes := session.Flashes()
	if err := session.Save(c.Request, c.Writer); err != nil {
		return nil, err
	}
	return flashes, nil
}

func (rc *RequestContext) AddFlash(msg string) error {
	session, err := jar.Get(rc.Request, "ses")
	if err != nil {
		return err
	}
	session.AddFlash(msg)
	if err := session.Save(rc.Request, rc.Writer); err != nil {
		return err
	}
	return nil
}

func (c *RequestContext) Vars() map[string]string {
	return mux.Vars(c.Request)
}

func decodeForm(form interface{}, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	return decoder.Decode(form, r.PostForm)
}
