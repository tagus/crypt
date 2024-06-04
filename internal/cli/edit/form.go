package edit

import (
	"context"
	"fmt"
	"github.com/rivo/tview"
	"github.com/tagus/crypt/internal/cli/environment"
	"github.com/tagus/crypt/internal/repos"
)

type Form struct {
	cr environment.CryptRepo
}

func (f *Form) Show(ctx context.Context, cred *repos.Credential) (*repos.Credential, error) {
	app := tview.NewApplication()

	updated := cred.Clone()

	var (
		wasModified     bool
		form            *tview.Form
		confirmPassword string
		err             error
	)
	form = tview.NewForm().
		AddInputField("service", cred.Service, 20, nil, onChange(&updated.Service)).
		AddInputField("email", cred.Email, 20, nil, onChange(&updated.Email)).
		AddPasswordField("password", cred.Password, 20, '*', onChange(&updated.Password)).
		AddPasswordField("confirm password", "", 20, '*', onChange(&confirmPassword)).
		AddTextArea("description", cred.Description, 40, 0, 0, onChange(&updated.Description)).
		AddButton("save", func() {
			if cred.Service != updated.Service ||
				cred.Email != updated.Email ||
				cred.Password != updated.Password ||
				cred.Description != updated.Description {
				wasModified = true
			}
			if wasModified {
				if updated.Password != confirmPassword {
					app.SetRoot(f.buildModal("passwords do not match", func() {
						app.SetRoot(form, true)
					}), true)
					return
				}

				updated, err = f.cr.UpdateCredential(ctx, updated)
				app.Stop()
			} else {
				app.Stop()
			}
		}).
		AddButton("quit", func() {
			app.Stop()
		})

	form.SetBorder(true).
		SetTitle(fmt.Sprintf("editing %v", cred)).
		SetTitleAlign(tview.AlignLeft)

	app = app.SetRoot(form, true).
		EnableMouse(true).
		EnablePaste(true)

	if err := app.Run(); err != nil {
		return nil, err
	}
	if wasModified {
		return updated, err
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
