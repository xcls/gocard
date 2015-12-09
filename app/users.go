package app

import (
	"net/http"

	"github.com/mcls/gocard/stores/common"
	"github.com/mcls/gocard/valid"
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

func (f *NewUserForm) Validate() []error {
	vd := valid.NewValidator()
	vd.ValidateMinLength("Email", f.Email, 4)
	vd.ValidateMinLength("Password", f.Password, 6)
	vd.ValidateMinLength("Password Confirmation", f.PasswordConfirmation, 6)
	vd.ValidateConfirmation("Password", f.Password, f.PasswordConfirmation)
	return vd.Errors()
}

func RegisterHandler(rc *RequestContext) error {
	form := new(NewUserForm)
	if rc.Request.Method == "GET" {
		return rc.HTML(http.StatusOK, "users/register", tplVars{"User": form})
	}

	decodeForm(form, rc.Request)
	formErrors := form.Validate()
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
