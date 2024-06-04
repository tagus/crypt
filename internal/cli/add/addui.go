package add

import (
	"context"
	"github.com/rivo/tview"
	"github.com/tagus/crypt/internal/cli/environment"
	"github.com/tagus/crypt/internal/repos"
	"github.com/tagus/mango"
)

type Form struct {
	cr environment.CryptRepo
}

func (f *Form) Show(ctx context.Context, service string) (*repos.Credential, error) {
	cred := &repos.Credential{
		ID:      mango.ShortID(),
		Service: service,
	}
	app := tview.NewApplication()

	var (
		wasSubmitted    bool
		form            *tview.Form
		confirmPassword string
		err             error
	)
	form = tview.NewForm().
		AddInputField("service", service, 20, nil, onChange(&cred.Service)).
		AddInputField("email", "", 20, nil, onChange(&cred.Email)).
		AddPasswordField("password", "", 20, '*', onChange(&cred.Password)).
		AddPasswordField("confirm password", "", 20, '*', onChange(&confirmPassword)).
		AddTextArea("description", "", 40, 0, 0, onChange(&cred.Description)).
		AddButton("save", func() {
			if cred.Password != confirmPassword {
				modal := f.buildModal("passwords do not match", func() {
					app.SetRoot(form, true)
				})
				app.SetRoot(modal, false)
				return
			}
			cred, err = f.cr.InsertCredential(ctx, cred)
			wasSubmitted = true
			app.Stop()
		}).
		AddButton("quit", func() {
			app.Stop()
		})

	form.SetBorder(true).
		SetTitle("add a new service credential").
		SetTitleAlign(tview.AlignLeft)

	app = app.SetRoot(form, true).
		EnableMouse(true).
		EnablePaste(true)

	if err := app.Run(); err != nil {
		return nil, err
	}
	if wasSubmitted {
		return cred, err
	}
	return nil, nil
}

func onChange(val *string) func(text string) {
	return func(text string) {
		*val = text
	}
}

func (f *Form) buildModal(msg string, onContinue func()) *tview.Modal {
	modal := tview.NewModal().
		SetText(msg).
		AddButtons([]string{"ok"})

	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "ok" {
			onContinue()
		}
	})

	return modal
}
