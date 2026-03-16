package ui

import (
	"strings"

	"soloterm/domain/snippet"
	sharedui "soloterm/shared/ui"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// SnippetView provides snippet management and selection UI
type SnippetView struct {
	app            *App
	snippetService *snippet.Service

	Modal          *tview.Flex
	frame          *tview.Frame
	table          *tview.Table
	Form           *SnippetForm
	errorView      *tview.TextView
	buttonRow      *tview.Flex
	buttons        []*tview.Button
	formArea       *tview.Flex
	innerContent   *tview.Flex
	formBaseHeight int
	formVisible    bool

	focused     tview.Primitive // tracks focused element for Tab cycling
	returnFocus tview.Primitive
}

// NewSnippetView creates a new snippet view
func NewSnippetView(app *App, snippetService *snippet.Service) *SnippetView {
	sv := &SnippetView{app: app, snippetService: snippetService}
	sv.setup()
	return sv
}

func (sv *SnippetView) setup() {
	sv.setupTable()
	sv.Form = NewSnippetForm()
	sv.setupButtons()
	sv.setupLayout()
	sv.setupKeyBindings()
}

func (sv *SnippetView) setupTable() {
	sv.table = tview.NewTable().
		SetSelectable(true, false).
		SetSelectedStyle(tcell.Style{}.Background(tcell.ColorAqua).Foreground(tcell.ColorBlack))
	sv.table.SetBorder(false)

	sv.table.SetSelectionChangedFunc(func(row, _ int) {
		sv.onRowSelected(row)
	})
	sv.table.SetSelectedFunc(func(_, _ int) {
		sv.handleUse()
	})
}

func (sv *SnippetView) setupButtons() {
	sv.buttonRow = tview.NewFlex()
	sv.rebuildButtons()
}

func (sv *SnippetView) addButton(label string, width int, selected func()) {
	btn := tview.NewButton(label).SetSelectedFunc(selected)
	btn.SetFocusFunc(func() {
		sv.focused = btn
		sv.frame.SetBorderColor(Style.BorderFocusColor)
	})
	btn.SetBlurFunc(func() { sv.frame.SetBorderColor(Style.BorderColor) })
	sv.buttons = append(sv.buttons, btn)
	sv.buttonRow.AddItem(btn, width, 0, false)
}

func (sv *SnippetView) rebuildButtons() {
	sv.buttonRow.Clear()
	sv.buttons = nil
	sv.buttonRow.AddItem(nil, 0, 1, false)

	sv.addButton("Save", 8, func() { sv.handleSave() })
	sv.buttonRow.AddItem(nil, 1, 0, false)
	sv.addButton("Delete", 8, func() { sv.handleDelete() })

	sv.buttonRow.AddItem(nil, 0, 1, false)
}

func (sv *SnippetView) setupLayout() {
	tableContainer := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(sv.table, 0, 1, true)
	tableContainer.SetBorder(true)

	sv.errorView = tview.NewTextView().SetDynamicColors(true).SetWrap(false)

	sv.formArea = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(sv.Form, 0, 1, false).
		AddItem(sv.errorView, 0, 0, false)
	sv.formArea.SetBorder(true)

	sv.formBaseHeight = 11
	sv.Form.SetErrorChangeHandler(func(errors map[string]string) {
		text := sharedui.FormatErrors(errors)
		sv.errorView.SetText(text)
		errorLines := 0
		if text != "" {
			errorLines = len(errors)
		}
		sv.formArea.ResizeItem(sv.errorView, errorLines, 0)
		if sv.formVisible {
			sv.innerContent.ResizeItem(sv.formArea, sv.formBaseHeight+errorLines, 0)
		}
	})

	sv.innerContent = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tableContainer, 0, 1, true).
		AddItem(sv.formArea, 0, 0, false).
		AddItem(sv.buttonRow, 0, 0, false)

	sv.frame = tview.NewFrame(sv.innerContent).
		SetBorders(1, 0, 0, 0, 1, 1)
	sv.frame.SetBorder(true).
		SetTitleAlign(tview.AlignLeft).
		SetTitle("[::b] Snippets ([" + Style.HelpKeyTextColor + "]Esc[" + Style.NormalTextColor + "] Close) [-::-]")

	sv.Modal = tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(
			tview.NewFlex().
				SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(sv.frame, 30, 0, true).
				AddItem(nil, 0, 1, false),
			70, 1, true,
		).
		AddItem(nil, 0, 1, false)

	sv.table.SetFocusFunc(func() {
		sv.focused = sv.table
		sv.app.updateFooterHelp(helpBar("Snippets", []helpEntry{
			{"↑/↓", "Navigate"},
			{"Enter", "Use"},
			{"Ctrl+E", "Edit"},
			{"Ctrl+U/D", "Reorder"},
			{"Ctrl+N", "New"},
			{"F12", "Help"},
			{"Esc", "Close"},
		}))
		sv.frame.SetBorderColor(Style.BorderFocusColor)
		tableContainer.SetBorderColor(Style.BorderFocusColor)
	})
	sv.table.SetBlurFunc(func() {
		sv.frame.SetBorderColor(Style.BorderColor)
		tableContainer.SetBorderColor(Style.BorderColor)
	})

	sv.Form.nameField.SetFocusFunc(sv.formFieldFocusFunc(sv.Form.nameField))
	sv.Form.contentField.SetFocusFunc(sv.formFieldFocusFunc(sv.Form.contentField))
	sv.Form.gameDropdown.SetFocusFunc(sv.formFieldFocusFunc(sv.Form.gameDropdown))
	sv.Form.SetFocusFunc(func() {
		sv.app.updateFooterHelp(helpBar("Snippets", []helpEntry{
			{"Ctrl+S", "Save"},
			{"Ctrl+N", "New"},
			{"Esc", "Close"},
		}))
	})
	sv.Form.SetBlurFunc(func() {
		sv.frame.SetBorderColor(Style.BorderColor)
		sv.formArea.SetBorderColor(Style.BorderColor)
	})
}

func (sv *SnippetView) formFieldFocusFunc(p tview.Primitive) func() {
	return func() {
		sv.focused = p
		sv.frame.SetBorderColor(Style.BorderFocusColor)
		sv.formArea.SetBorderColor(Style.BorderFocusColor)
	}
}

func (sv *SnippetView) showForm() {
	if sv.formVisible {
		return
	}
	sv.innerContent.ResizeItem(sv.formArea, sv.formBaseHeight, 0)
	sv.innerContent.ResizeItem(sv.buttonRow, 1, 0)
	sv.formVisible = true
}

func (sv *SnippetView) hideForm() {
	if !sv.formVisible {
		return
	}
	sv.innerContent.ResizeItem(sv.formArea, 0, 0)
	sv.innerContent.ResizeItem(sv.buttonRow, 0, 0)
	sv.formVisible = false
}

func (sv *SnippetView) setupKeyBindings() {
	sv.table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlU:
			sv.handleReorder(-1)
			return nil
		case tcell.KeyCtrlD:
			sv.handleReorder(1)
			return nil
		case tcell.KeyCtrlE:
			row, _ := sv.table.GetSelection()
			sv.onRowSelected(row)
			sv.showDeleteBtn()
			sv.showForm()
			sv.app.SetFocus(sv.Form)
			return nil
		case tcell.KeyUp, tcell.KeyDown, tcell.KeyPgUp, tcell.KeyPgDn, tcell.KeyHome, tcell.KeyEnd:
			sv.hideForm()
			return event
		}
		return event
	})

	// Form boundary: intercept Tab/BackTab only when leaving the form.
	sv.Form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			if sv.focused == sv.Form.gameDropdown {
				if len(sv.buttons) > 0 {
					sv.app.SetFocus(sv.buttons[0])
				} else {
					sv.app.SetFocus(sv.table)
				}
				return nil
			}
		case tcell.KeyBacktab:
			if sv.focused == sv.Form.nameField {
				sv.app.SetFocus(sv.table)
				return nil
			}
		case tcell.KeyCtrlS:
			sv.handleSave()
			return nil
		}
		return event
	})

	sv.Modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			switch sv.focused {
			case sv.table:
				if sv.formVisible {
					sv.app.SetFocus(sv.Form)
				}
				return nil
			case sv.Form.nameField, sv.Form.contentField, sv.Form.gameDropdown:
				return event // let Form handle internal Tab cycling
			default:
				sv.shiftButton(1)
				return nil
			}
		case tcell.KeyBacktab:
			switch sv.focused {
			case sv.table:
				if sv.formVisible && len(sv.buttons) > 0 {
					sv.app.SetFocus(sv.buttons[len(sv.buttons)-1])
				}
				return nil
			case sv.Form.nameField, sv.Form.contentField, sv.Form.gameDropdown:
				return event // let Form.SetInputCapture handle boundary
			default:
				sv.shiftButton(-1)
				return nil
			}
		case tcell.KeyCtrlN:
			sv.Form.Reset(sv.activeGameID())
			sv.hideDeleteBtn()
			sv.showForm()
			sv.app.SetFocus(sv.Form)
			return nil
		case tcell.KeyEsc:
			sv.app.HandleEvent(&SnippetCancelEvent{
				BaseEvent: BaseEvent{action: SNIPPET_CANCEL},
			})
			return nil
		case tcell.KeyF12:
			sv.app.HandleEvent(&ShowHelpEvent{
				BaseEvent:   BaseEvent{action: SHOW_HELP},
				Title:       "Snippets Help",
				ReturnFocus: sv.Modal,
				Text:        sv.buildHelpText(),
			})
			return nil
		}
		return event
	})
}

func (sv *SnippetView) hideDeleteBtn() {
	sv.buttonRow.ResizeItem(sv.buttons[1], 0, 0)
}

func (sv *SnippetView) showDeleteBtn() {
	sv.buttonRow.ResizeItem(sv.buttons[1], 10, 0)
}

// shiftButton moves focus to the next (+1) or previous (-1) visible button.
// Wraps from the last button to the table, and from the first button to the form.
func (sv *SnippetView) shiftButton(dir int) {
	var visible []*tview.Button
	for _, b := range sv.buttons {
		_, _, w, _ := b.GetRect()
		if w > 0 {
			visible = append(visible, b)
		}
	}
	for i, b := range visible {
		if sv.focused == b {
			next := i + dir
			if next < 0 {
				sv.app.SetFocus(sv.Form)
			} else if next >= len(visible) {
				sv.app.SetFocus(sv.table)
			} else {
				sv.app.SetFocus(visible[next])
			}
			return
		}
	}
}

func (sv *SnippetView) onRowSelected(row int) {
	ref := sv.table.GetCell(row, 0).GetReference()
	if ref == nil {
		return
	}
	sr, err := sv.snippetService.GetByID(ref.(int64))
	if err != nil {
		return
	}
	sv.Form.Populate(sr)
}

// refreshGames repopulates the game dropdown on the form.
func (sv *SnippetView) refreshGames() {
	options := []GameOption{{ID: nil, Name: "-- Global --"}}
	if g := sv.app.CurrentGame(); g != nil {
		options = append(options, GameOption{ID: &g.ID, Name: g.Name})
	}
	sv.Form.SetGames(options)
}

// Refresh reloads snippets from the database into the table
func (sv *SnippetView) Refresh() {
	sv.refreshGames()
	sv.table.Clear()

	activeGameID := sv.activeGameID()

	var gameSnippets, globalSnippets []*snippet.Snippet
	var err error

	const loadErrMsg = "Error loading snippets"

	if activeGameID != nil {
		gameSnippets, err = sv.snippetService.GetByGameID(*activeGameID)
		if err != nil {
			sv.app.notification.ShowError(loadErrMsg)
			return
		}
	}
	globalSnippets, err = sv.snippetService.GetGlobal()
	if err != nil {
		sv.app.notification.ShowError(loadErrMsg)
		return
	}

	if len(gameSnippets)+len(globalSnippets) == 0 {
		sv.table.SetCell(0, 0, tview.NewTableCell("No snippets yet. Press Ctrl+N to add one.").
			SetTextColor(Style.EmptyStateMessageColor).
			SetSelectable(false))
		sv.Form.Reset(activeGameID)
		return
	}

	row := 0
	for _, s := range gameSnippets {
		sv.addSnippetRow(row, s)
		row++
	}
	if len(gameSnippets) > 0 && len(globalSnippets) > 0 {
		sv.addSectionDivider(row, "─── Global ───")
		row++
	}
	for _, s := range globalSnippets {
		sv.addSnippetRow(row, s)
		row++
	}

	sv.table.Select(0, 0)
}

func (sv *SnippetView) addSnippetRow(row int, s *snippet.Snippet) {
	sv.table.SetCell(row, 0, tview.NewTableCell(s.Name).SetReference(s.ID).SetExpansion(1))
	sv.table.SetCell(row, 1, tview.NewTableCell(s.Content).SetExpansion(2))
}

func (sv *SnippetView) addSectionDivider(row int, label string) {
	sv.table.SetCell(row, 0, tview.NewTableCell(label).
		SetTextColor(tcell.ColorYellow).
		SetSelectable(false).
		SetExpansion(1))
	sv.table.SetCell(row, 1, tview.NewTableCell("").
		SetSelectable(false).
		SetExpansion(2))
}

func (sv *SnippetView) handleSave() {
	sr := sv.Form.BuildDomain()
	isNew := sr.IsNew()
	if isNew {
		all, err := sv.snippetService.GetAll()
		if err == nil {
			sr.Position = len(all)
		}
	}
	saved, err := sv.snippetService.Save(sr)
	if err != nil {
		if sharedui.HandleValidationError(err, sv.Form) {
			return
		}
		sv.app.notification.ShowError("Failed to save snippet: " + err.Error())
		return
	}
	sv.app.HandleEvent(&SnippetSavedEvent{
		BaseEvent: BaseEvent{action: SNIPPET_SAVED},
		Snippet:   saved,
		IsNew:     isNew,
	})
}

// selectedID returns the snippet ID for the currently selected table row,
// and false if the row has no associated snippet (e.g. a divider row).
func (sv *SnippetView) selectedID() (int64, bool) {
	row, _ := sv.table.GetSelection()
	ref := sv.table.GetCell(row, 0).GetReference()
	if ref == nil {
		return 0, false
	}
	return ref.(int64), true
}

func (sv *SnippetView) handleReorder(direction int) {
	id, ok := sv.selectedID()
	if !ok {
		return
	}
	sv.app.HandleEvent(&SnippetReorderEvent{
		BaseEvent: BaseEvent{action: SNIPPET_REORDER},
		SnippetID: id,
		Direction: direction,
	})
}

func (sv *SnippetView) handleUse() {
	id, ok := sv.selectedID()
	if !ok {
		sv.app.notification.ShowWarning("Select a snippet to use")
		return
	}
	sr, err := sv.snippetService.GetByID(id)
	if err != nil || sr.Content == "" {
		sv.app.notification.ShowWarning("Select a snippet to use")
		return
	}
	sv.app.HandleEvent(&SnippetUseEvent{
		BaseEvent: BaseEvent{action: SNIPPET_USE},
		Content:   sr.Content,
	})
}

func (sv *SnippetView) handleDelete() {
	id, ok := sv.selectedID()
	if !ok {
		sv.app.notification.ShowWarning("Select a snippet to delete")
		return
	}
	sv.app.HandleEvent(&SnippetDeleteConfirmEvent{
		BaseEvent: BaseEvent{action: SNIPPET_DELETE_CONFIRM},
		SnippetID: id,
	})
}

func (sv *SnippetView) activeGameID() *int64 {
	if g := sv.app.CurrentGame(); g != nil {
		return &g.ID
	}
	return nil
}

func (sv *SnippetView) buildHelpText() string {
	return strings.NewReplacer(
		"[yellow]", "["+Style.HelpKeyTextColor+"]",
		"[white]", "["+Style.NormalTextColor+"]",
		"[green]", "["+Style.HelpSectionColor+"]",
	).Replace(`[green]What are Snippets?[white]

Snippets are frequently used rolls for quick reuse. Any expression you can type in the dice roller can be saved as a snippet: dice, lists, table references, or combinations of all three.

  [yellow]Characters[white]       {Frank; Bill; Joe}
  [yellow]Sparks[white]           @actions, @themes
  [yellow]Encounter[white]        @creatures, 2d4+1

[green]Using Snippets[white]

  [yellow]Enter[white]       Insert the selected snippet content into the dice input
  [yellow]Ctrl+E[white]      Open the edit form for the selected snippet
  [yellow]Ctrl+N[white]      Clear the form to create a new snippet
  [yellow]Ctrl+S[white]      Save the current form (when form is open)
  [yellow]Ctrl+U/D[white]    Reorder the selected snippet up or down

[green]Content Format[white]

The Content field accepts plain text labels, dice expressions, or both. Multiple expressions can be written on separate lines.

  [yellow]2d6[white]                   Roll 2 six-sided dice
  [yellow]Attack: 1d20+5[white]        Labelled roll
  [yellow]{A; B; C}[white]             Pick randomly from a list
  [yellow]@tablename[white]            Roll on a table
  [yellow]@category/tablename[white]   Roll on a scoped table

[green]Game Scope[white]

Snippets can be global (available everywhere) or tied to the active game. When a game is loaded, game snippets appear first in the list followed by global snippets.
`)
}

func (sv *SnippetView) selectByID(id int64) {
	for row := 0; row < sv.table.GetRowCount(); row++ {
		ref := sv.table.GetCell(row, 0).GetReference()
		if ref != nil && ref.(int64) == id {
			sv.table.Select(row, 0)
			return
		}
	}
}
