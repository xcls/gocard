package app

import (
	"errors"
	"net/http"

	"github.com/xcls/gocard/stores/common"
	"github.com/xcls/gocard/valid"
)

type NewUserForm struct {
	Email                string
	Password             string
	PasswordConfirmation string
}

func (f *NewUserForm) ToModel() *common.User {
	user := new(common.User)
	user.Email = f.Email
	user.SetPassword(f.Password)
	return user
}

func (f *NewUserForm) Validate(users common.UserStore) []error {
	vd := valid.NewValidator()
	vd.ValidateMinLength("Email", f.Email, 4)
	vd.ValidateMinLength("Password", f.Password, 6)
	vd.ValidateMinLength("Password Confirmation", f.PasswordConfirmation, 6)
	vd.ValidateConfirmation("Password", f.Password, f.PasswordConfirmation)
	u, err := users.FindByEmail(f.Email)
	if err == nil || u != nil {
		vd.AddError(errors.New("User with that email already exists"))
	}
	return vd.Errors()
}

func RegisterHandler(rc *RequestContext) error {
	if rc.CurrentUser != nil {
		return rc.RedirectWithFlash("/", "Already logged in")
	}

	form := new(NewUserForm)
	if rc.Request.Method == "GET" {
		return rc.HTML(http.StatusOK, "users/register", tplVars{"User": form})
	}

	decodeForm(form, rc.Request)
	formErrors := form.Validate(rc.Store.Users)
	if len(formErrors) > 0 {
		return rc.HTML(http.StatusOK, "users/register", tplVars{
			"User":       form,
			"UserErrors": formErrors,
		})
	} else {
		user := form.ToModel()
		if err := rc.Store.Users.Insert(user); err != nil {
			return err
		}
		if err := rc.AddFlash("Registered " + user.Email); err != nil {
			return err
		}
		http.Redirect(rc.Writer, rc.Request, "/", http.StatusFound)
		return nil
	}
}

type LoginForm struct {
	Email    string
	Password string
}

func LoginHandler(rc *RequestContext) error {
	if rc.CurrentUser != nil {
		return rc.RedirectWithFlash("/", "Already logged in")
	}

	form := new(LoginForm)
	if rc.Request.Method == "GET" {
		return rc.HTML(http.StatusOK, "users/login", tplVars{"User": form})
	}

	decodeForm(form, rc.Request)
	user, err := rc.Store.Users.Authenticate(form.Email, form.Password)
	if err != nil {
		return rc.HTML(http.StatusOK, "users/login", tplVars{
			"User":       form,
			"UserErrors": []error{err},
		})
	}

	// Set user session
	userSession := common.NewUserSession(user.ID)
	err = rc.Store.UserSessions.Insert(userSession)
	if err != nil {
		return err
	}
	session, err := jar.Get(rc.Request, "uid")
	if err != nil {
		return err
	}
	session.Values["uid"] = userSession.UID
	if err := session.Save(rc.Request, rc.Writer); err != nil {
		return err
	}

	if err := rc.AddFlash("Welcome, " + user.Email); err != nil {
		return err
	}
	http.Redirect(rc.Writer, rc.Request, "/", http.StatusFound)
	return nil
}

func LogoutHandler(rc *RequestContext) error {
	session, err := jar.Get(rc.Request, "uid")
	if err != nil {
		return err
	}
	delete(session.Values, "uid")
	if err := session.Save(rc.Request, rc.Writer); err != nil {
		return err
	}
	if err := rc.AddFlash("Logged out"); err != nil {
		return err
	}
	http.Redirect(rc.Writer, rc.Request, "/", http.StatusFound)
	return nil
}
