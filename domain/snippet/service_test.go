package snippet

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	testhelper "soloterm/shared/testing"
)

func TestService_Save(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	defer testhelper.TeardownTestDB(t, db)
	svc := NewService(NewRepository(db))

	t.Run("valid save", func(t *testing.T) {
		sr := &Snippet{Name: "Sparks", Content: "@actions, @themes"}
		result, err := svc.Save(sr)
		require.NoError(t, err)
		assert.NotZero(t, result.ID)
	})

	t.Run("name is required", func(t *testing.T) {
		sr := &Snippet{Name: "", Content: "@actions"}
		_, err := svc.Save(sr)
		assert.Error(t, err)
	})

	t.Run("roll is required", func(t *testing.T) {
		sr := &Snippet{Name: "Sparks", Content: ""}
		_, err := svc.Save(sr)
		assert.Error(t, err)
	})

	t.Run("name too long", func(t *testing.T) {
		sr := &Snippet{Name: strings.Repeat("a", MaxNameLength+1), Content: "@actions"}
		_, err := svc.Save(sr)
		assert.Error(t, err)
	})

	t.Run("roll too long", func(t *testing.T) {
		sr := &Snippet{Name: "Sparks", Content: strings.Repeat("a", MaxContentLength+1)}
		_, err := svc.Save(sr)
		assert.Error(t, err)
	})
}

func TestService_Delete(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	defer testhelper.TeardownTestDB(t, db)
	svc := NewService(NewRepository(db))

	t.Run("delete existing", func(t *testing.T) {
		id := testhelper.CreateTestSnippet(t, db, "Sparks", "@actions", 0)
		err := svc.Delete(id)
		require.NoError(t, err)

		_, err = svc.GetByID(id)
		assert.Error(t, err)
	})

	t.Run("delete non-existent id", func(t *testing.T) {
		err := svc.Delete(99999)
		assert.Error(t, err)
	})
}

func TestService_Reorder(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	defer testhelper.TeardownTestDB(t, db)
	svc := NewService(NewRepository(db))

	idA := testhelper.CreateTestSnippet(t, db, "Alpha", "@actions", 0)
	idB := testhelper.CreateTestSnippet(t, db, "Bravo", "@themes", 1)
	idC := testhelper.CreateTestSnippet(t, db, "Charlie", "@moods", 2)

	t.Run("move down", func(t *testing.T) {
		_, err := svc.Reorder(idA, 1)
		require.NoError(t, err)

		all, _ := svc.GetAll()
		assert.Equal(t, idB, all[0].ID)
		assert.Equal(t, idA, all[1].ID)
		assert.Equal(t, idC, all[2].ID)

		// Restore original order
		_, _ = svc.Reorder(idA, -1)
	})

	t.Run("move up", func(t *testing.T) {
		_, err := svc.Reorder(idC, -1)
		require.NoError(t, err)

		all, _ := svc.GetAll()
		assert.Equal(t, idA, all[0].ID)
		assert.Equal(t, idC, all[1].ID)
		assert.Equal(t, idB, all[2].ID)

		// Restore original order
		_, _ = svc.Reorder(idC, 1)
	})

	t.Run("move first item up does nothing", func(t *testing.T) {
		returnedID, err := svc.Reorder(idA, -1)
		require.NoError(t, err)
		assert.Zero(t, returnedID)

		all, _ := svc.GetAll()
		assert.Equal(t, idA, all[0].ID)
	})

	t.Run("move last item down does nothing", func(t *testing.T) {
		returnedID, err := svc.Reorder(idC, 1)
		require.NoError(t, err)
		assert.Zero(t, returnedID)

		all, _ := svc.GetAll()
		assert.Equal(t, idC, all[2].ID)
	})

	t.Run("non-existent id returns error", func(t *testing.T) {
		_, err := svc.Reorder(99999, 1)
		assert.Error(t, err)
	})
}
