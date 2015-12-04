package app

import "net/http"

func renderHTML(w http.ResponseWriter, r *http.Request, status int, tpl string, vars tplVars) error {
	if vars == nil {
		vars = tplVars{}
	}
	flashes, err := getFlashes(w, r)
	if err != nil {
		return err
	}
	vars["Flashes"] = flashes
	renderer.HTML(w, status, tpl, vars)
	return nil
}

func addFlash(w http.ResponseWriter, r *http.Request, msg string) error {
	session, err := jar.Get(r, "ses")
	if err != nil {
		return err
	}
	session.AddFlash(msg)
	if err := session.Save(r, w); err != nil {
		return err
	}
	return nil
}

func getFlashes(w http.ResponseWriter, r *http.Request) ([]interface{}, error) {
	session, err := jar.Get(r, "ses")
	if err != nil {
		return nil, err
	}
	flashes := session.Flashes()
	if err := session.Save(r, w); err != nil {
		return nil, err
	}
	return flashes, nil
}

func decodeForm(form interface{}, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	return decoder.Decode(form, r.PostForm)
}
