package ui

import (
	"soloterm/domain/snippet"
	sharedui "soloterm/shared/ui"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// GameOption is a lightweight game reference used to populate the snippet form dropdown.
type GameOption struct {
	ID   *int64
	Name string
}

// SnippetForm represents a form for creating/editing snippets
type SnippetForm struct {
	*sharedui.DataForm
	snippetID    *int64
	snippetPos   int
	nameField    *tview.InputField
	contentField *tview.TextArea
	gameDropdown *tview.DropDown
	gameOptions  []GameOption // parallel to dropdown; gameOptions[0] is always Global (ID nil)
}

// NewSnippetForm creates a new snippet form
func NewSnippetForm() *SnippetForm {
	f := &SnippetForm{
		DataForm: sharedui.NewDataForm(),
	}

	f.nameField = tview.NewInputField().
		SetLabel("Name").
		SetFieldBackgroundColor(tcell.ColorDefault).
		SetFieldWidth(0)

	f.contentField = tview.NewTextArea().
		SetLabel("Content").
		SetMaxLength(snippet.MaxContentLength).
		SetSize(3, 0)

	f.gameDropdown = tview.NewDropDown().
		SetLabel("Game").
		SetFieldBackgroundColor(tcell.ColorDefault)

	f.Clear(true)
	f.AddFormItem(f.nameField)
	f.AddFormItem(f.contentField)
	f.AddFormItem(f.gameDropdown)
	f.SetBorder(false)
	f.SetButtonsAlign(tview.AlignCenter)
	f.SetItemPadding(1)

	return f
}

// SetGames populates the game dropdown. The first option is always "Global" (nil ID).
func (f *SnippetForm) SetGames(options []GameOption) {
	f.gameOptions = options
	names := make([]string, len(options))
	for i, o := range options {
		names[i] = o.Name
	}
	f.gameDropdown.SetOptions(names, nil)
	f.gameDropdown.SetCurrentOption(0)
}

// Reset clears the form for a new entry. If activeGameID is non-nil, the
// dropdown is pre-selected to that game; otherwise it defaults to Global.
func (f *SnippetForm) Reset(activeGameID *int64) {
	f.snippetID = nil
	f.snippetPos = 0
	f.nameField.SetText("")
	f.contentField.SetText("", false)
	f.gameDropdown.SetCurrentOption(0)
	if activeGameID != nil {
		for i, o := range f.gameOptions {
			if o.ID != nil && *o.ID == *activeGameID {
				f.gameDropdown.SetCurrentOption(i)
				break
			}
		}
	}
	f.ClearFieldErrors()
	f.SetFocus(0)
}

// Populate fills the form with an existing snippet's data
func (f *SnippetForm) Populate(sr *snippet.Snippet) {
	f.snippetID = &sr.ID
	f.snippetPos = sr.Position
	f.nameField.SetText(sr.Name)
	f.contentField.SetText(sr.Content, false)
	f.gameDropdown.SetCurrentOption(0)
	for i, o := range f.gameOptions {
		if o.ID != nil && sr.GameID != nil && *o.ID == *sr.GameID {
			f.gameDropdown.SetCurrentOption(i)
			break
		}
	}
	f.ClearFieldErrors()
	f.SetFocus(0)
}

// BuildDomain constructs a Snippet from the current form values
func (f *SnippetForm) BuildDomain() *snippet.Snippet {
	sr := &snippet.Snippet{
		Name:     f.nameField.GetText(),
		Content:  f.contentField.GetText(),
		Position: f.snippetPos,
	}
	if f.snippetID != nil {
		sr.ID = *f.snippetID
	}
	if idx, _ := f.gameDropdown.GetCurrentOption(); idx >= 0 && idx < len(f.gameOptions) {
		sr.GameID = f.gameOptions[idx].ID
	}
	return sr
}

// SetFieldErrors sets errors and updates field labels
func (f *SnippetForm) SetFieldErrors(errors map[string]string) {
	f.DataForm.SetFieldErrors(errors)
	f.updateFieldLabels()
}

// ClearFieldErrors removes all error highlights
func (f *SnippetForm) ClearFieldErrors() {
	f.DataForm.ClearFieldErrors()
	f.updateFieldLabels()
}

func (f *SnippetForm) updateFieldLabels() {
	if f.HasFieldError("name") {
		f.nameField.SetLabel("[" + Style.ErrorTextColor + "]Name[" + Style.NormalTextColor + "]")
	} else {
		f.nameField.SetLabel("Name")
	}

	if f.HasFieldError("content") {
		f.contentField.SetLabel("[" + Style.ErrorTextColor + "]Content[" + Style.NormalTextColor + "]")
	} else {
		f.contentField.SetLabel("Content")
	}

	if f.HasFieldError("game_id") {
		f.gameDropdown.SetLabel("[" + Style.ErrorTextColor + "]Game[" + Style.NormalTextColor + "]")
	} else {
		f.gameDropdown.SetLabel("Game")
	}

}
