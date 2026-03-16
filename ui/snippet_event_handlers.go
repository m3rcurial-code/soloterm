package ui

func (a *App) handleSnippetShow(_ *SnippetShowEvent) {
	a.snippetView.returnFocus = a.GetFocus()
	a.snippetView.Refresh()
	a.pages.ShowPage(SNIPPET_MODAL_ID)
	a.SetFocus(a.snippetView.table)
}

func (a *App) handleSnippetCancel(_ *SnippetCancelEvent) {
	a.pages.HidePage(SNIPPET_MODAL_ID)
	if a.snippetView.returnFocus != nil {
		a.SetFocus(a.snippetView.returnFocus)
	}
}

func (a *App) handleSnippetShowNew(_ *SnippetShowNewEvent) {
	a.snippetView.refreshGames()
	a.snippetView.Form.Reset(a.snippetView.activeGameID())
	a.snippetView.Form.RemoveDeleteButton()
	a.snippetView.formModal.SetTitle(" New Snippet ")
	a.pages.ShowPage(SNIPPET_FORM_MODAL_ID)
	a.SetFocus(a.snippetView.Form)
}

func (a *App) handleSnippetShowEdit(e *SnippetShowEditEvent) {
	a.snippetView.refreshGames()
	a.snippetView.Form.PopulateForEdit(e.Snippet)
	a.snippetView.Form.AddDeleteButton()
	a.snippetView.formModal.SetTitle(" Edit Snippet ")
	a.pages.ShowPage(SNIPPET_FORM_MODAL_ID)
	a.SetFocus(a.snippetView.Form)
}

func (a *App) handleSnippetFormCancel(_ *SnippetFormCancelEvent) {
	a.snippetView.Form.ClearFieldErrors()
	a.pages.HidePage(SNIPPET_FORM_MODAL_ID)
	a.SetFocus(a.snippetView.table)
}

func (a *App) handleSnippetSaved(e *SnippetSavedEvent) {
	a.snippetView.Form.ClearFieldErrors()
	a.pages.HidePage(SNIPPET_FORM_MODAL_ID)
	a.snippetView.Refresh()
	a.snippetView.selectByID(e.Snippet.ID)
	a.SetFocus(a.snippetView.table)
	a.notification.ShowSuccess("Snippet saved successfully")
}

func (a *App) handleSnippetDeleteConfirm(e *SnippetDeleteConfirmEvent) {
	returnFocus := a.GetFocus()
	a.confirmModal.Configure(
		"Are you sure you want to delete this snippet?",
		func() {
			err := a.snippetView.snippetService.Delete(e.SnippetID)
			if err != nil {
				a.HandleEvent(&SnippetDeleteFailedEvent{
					BaseEvent: BaseEvent{action: SNIPPET_DELETE_FAILED},
					Error:     err,
				})
				return
			}
			a.HandleEvent(&SnippetDeletedEvent{
				BaseEvent: BaseEvent{action: SNIPPET_DELETED},
			})
		},
		func() {
			a.pages.HidePage(CONFIRM_MODAL_ID)
			a.SetFocus(returnFocus)
		},
	)
	a.pages.ShowPage(CONFIRM_MODAL_ID)
}

func (a *App) handleSnippetDeleted(_ *SnippetDeletedEvent) {
	a.pages.HidePage(CONFIRM_MODAL_ID)
	a.pages.HidePage(SNIPPET_FORM_MODAL_ID)
	a.snippetView.Refresh()
	a.SetFocus(a.snippetView.table)
	a.notification.ShowSuccess("Snippet deleted successfully")
}

func (a *App) handleSnippetDeleteFailed(e *SnippetDeleteFailedEvent) {
	a.pages.HidePage(CONFIRM_MODAL_ID)
	a.notification.ShowError("Failed to delete snippet: " + e.Error.Error())
}

func (a *App) handleSnippetReorder(e *SnippetReorderEvent) {
	id, err := a.snippetView.snippetService.Reorder(e.SnippetID, e.Direction)
	if err != nil {
		a.notification.ShowError("Failed to reorder snippet: " + err.Error())
		return
	}
	if id == 0 {
		return // boundary, nothing changed
	}
	a.snippetView.Refresh()
	a.snippetView.selectByID(id)
}

func (a *App) handleSnippetUse(e *SnippetUseEvent) {
	a.pages.HidePage(SNIPPET_MODAL_ID)
	_, start, _ := a.diceView.TextArea.GetSelection()
	a.diceView.TextArea.Replace(start, start, e.Content+" ")
	a.SetFocus(a.diceView.TextArea)
}
