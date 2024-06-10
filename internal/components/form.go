package components

import (
	"github.com/rivo/tview"
	"github.com/tagus/crypt/internal/repos"
	"github.com/tagus/mango"
	"strings"
)

type ShowModalFn func(string)

type HandleSubmitFn func(cred *repos.Credential, fn ShowModalFn) (*repos.Credential, error)

type Form struct {
	src  *repos.Credential
	dst  *repos.Credential
	app  *tview.Application
	form *tview.Form

	onSave HandleSubmitFn

	title string
}

type FormOpts struct {
	Source      *repos.Credential
	Destination *repos.Credential
	Title       string
	OnSave      HandleSubmitFn
}

func NewForm(opts FormOpts) *Form {
	var src *repos.Credential
	if opts.Source == nil {
		src = &repos.Credential{}
	} else {
		src = opts.Source
	}

	var dst *repos.Credential
	if opts.Destination == nil {
		dst = &repos.Credential{}
	} else {
		dst = opts.Destination
	}

	return &Form{
		src:    src,
		dst:    dst,
		title:  opts.Title,
		app:    tview.NewApplication(),
		onSave: opts.OnSave,
	}
}

func (f *Form) Show() (*repos.Credential, error) {
	var (
		confirmPassword string
		final           *repos.Credential
		err             error
	)
	f.form = tview.NewForm().
		AddInputField("service", f.src.Service, 20, nil, onChange(&f.dst.Service)).
		AddInputField("email", f.src.Email, 20, nil, onChange(&f.dst.Email)).
		AddInputField("username", f.src.Username, 20, nil, onChange(&f.dst.Username)).
		AddPasswordField("password", f.src.Password, 20, '*', onChange(&f.dst.Password)).
		AddPasswordField("confirm password", "", 20, '*', onChange(&confirmPassword)).
		AddTextArea("description", f.src.Description, 40, 0, 0, onChange(&f.dst.Description)).
		AddInputField("tags", strings.Join(f.src.Tags, ", "), 40, nil, onTagsChange(&f.dst.Tags)).
		AddButton("save", func() {
			if f.src.Password != f.dst.Password && f.dst.Password != confirmPassword {
				f.showModal("passwords do not match")
				return
			}
			f.dst.ID = mango.ShortID()
			final, err = f.onSave(f.dst, f.showModal)
			if err != nil {
				f.showModal(err.Error())
			} else {
				f.app.Stop()
			}
		}).
		AddButton("quit", func() {
			f.app.Stop()
		})

	f.form.SetBorder(true).
		SetTitle(f.title).
		SetTitleAlign(tview.AlignLeft)

	f.app = f.app.SetRoot(f.form, true).
		EnableMouse(true).
		EnablePaste(true)

	if err := f.app.Run(); err != nil {
		return nil, err
	}
	return final, err
}

func (f *Form) showModal(title string) {
	modal := tview.NewModal().
		SetText(title).
		AddButtons([]string{"ok"})

	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "ok" {
			f.app.SetRoot(f.form, true)
		}
	})

	f.app.SetRoot(modal, true)
}

func onChange(val *string) func(text string) {
	return func(text string) {
		*val = text
	}
}

func onTagsChange(val *[]string) func(text string) {
	return func(text string) {
		*val = mango.Map(strings.Split(text, ","), func(val string) string {
			return strings.TrimSpace(val)
		})
	}
}
