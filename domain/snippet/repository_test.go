package snippet

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	testhelper "soloterm/shared/testing"
)

func TestRepository_Save(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	defer testhelper.TeardownTestDB(t, db)
	repo := NewRepository(db)

	t.Run("insert new snippet", func(t *testing.T) {
		sr := &Snippet{Name: "Sparks", Content: "@actions, @themes", Position: 0}

		err := repo.Save(sr)
		require.NoError(t, err)

		assert.NotZero(t, sr.ID)
		assert.False(t, sr.CreatedAt.IsZero())
		assert.False(t, sr.UpdatedAt.IsZero())
	})

	t.Run("update existing snippet", func(t *testing.T) {
		sr := &Snippet{Name: "Room", Content: "@descriptors", Position: 0}
		require.NoError(t, repo.Save(sr))

		firstUpdatedAt := sr.UpdatedAt
		time.Sleep(10 * time.Millisecond)

		sr.Content = "@descriptors, @room_contents"
		require.NoError(t, repo.Save(sr))

		assert.Equal(t, "@descriptors, @room_contents", sr.Content)
		assert.NotEqual(t, firstUpdatedAt, sr.UpdatedAt)
	})
}

func TestRepository_GetByID(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	defer testhelper.TeardownTestDB(t, db)
	repo := NewRepository(db)

	id := testhelper.CreateTestSnippet(t, db, "Sparks", "@actions, @themes", 0)

	t.Run("get existing", func(t *testing.T) {
		sr, err := repo.GetByID(id)
		require.NoError(t, err)
		assert.Equal(t, id, sr.ID)
		assert.Equal(t, "Sparks", sr.Name)
	})

	t.Run("get non-existent id", func(t *testing.T) {
		_, err := repo.GetByID(99999)
		assert.Error(t, err)
	})

	t.Run("get zero id", func(t *testing.T) {
		_, err := repo.GetByID(0)
		assert.Error(t, err)
	})
}

func TestRepository_GetAll(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	defer testhelper.TeardownTestDB(t, db)
	repo := NewRepository(db)

	testhelper.CreateTestSnippet(t, db, "Sparks", "@actions, @themes", 0)
	testhelper.CreateTestSnippet(t, db, "Room", "@descriptors", 1)

	snippets, err := repo.GetAll()
	require.NoError(t, err)
	assert.Len(t, snippets, 2)
}

func TestRepository_GetAll_OrderedByPosition(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	defer testhelper.TeardownTestDB(t, db)
	repo := NewRepository(db)

	testhelper.CreateTestSnippet(t, db, "Bravo", "@actions", 1)
	testhelper.CreateTestSnippet(t, db, "Alpha", "@themes", 0)

	snippets, err := repo.GetAll()
	require.NoError(t, err)
	require.Len(t, snippets, 2)
	assert.Equal(t, "Alpha", snippets[0].Name)
	assert.Equal(t, "Bravo", snippets[1].Name)
}

func TestRepository_GetByName(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	defer testhelper.TeardownTestDB(t, db)
	repo := NewRepository(db)

	testhelper.CreateTestSnippet(t, db, "Sparks", "@actions, @themes", 0)

	t.Run("get existing by name", func(t *testing.T) {
		sr, err := repo.GetByName("Sparks")
		require.NoError(t, err)
		assert.Equal(t, "Sparks", sr.Name)
	})

	t.Run("case-insensitive match", func(t *testing.T) {
		sr, err := repo.GetByName("sparks")
		require.NoError(t, err)
		assert.Equal(t, "Sparks", sr.Name)
	})

	t.Run("non-existent name", func(t *testing.T) {
		_, err := repo.GetByName("nonexistent")
		assert.Error(t, err)
	})
}

func TestRepository_Delete(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	defer testhelper.TeardownTestDB(t, db)
	repo := NewRepository(db)

	t.Run("delete existing", func(t *testing.T) {
		id := testhelper.CreateTestSnippet(t, db, "Sparks", "@actions", 0)

		rows, err := repo.Delete(id)
		require.NoError(t, err)
		assert.Equal(t, int64(1), rows)

		_, err = repo.GetByID(id)
		assert.Error(t, err)
	})

	t.Run("delete non-existent id", func(t *testing.T) {
		_, err := repo.Delete(99999)
		assert.Error(t, err)
	})

	t.Run("delete zero id", func(t *testing.T) {
		_, err := repo.Delete(0)
		assert.Error(t, err)
	})
}

func TestRepository_SwapPositions(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	defer testhelper.TeardownTestDB(t, db)
	repo := NewRepository(db)

	idA := testhelper.CreateTestSnippet(t, db, "Alpha", "@actions", 0)
	idB := testhelper.CreateTestSnippet(t, db, "Bravo", "@themes", 1)

	err := repo.SwapPositions(idA, 0, idB, 1)
	require.NoError(t, err)

	a, _ := repo.GetByID(idA)
	b, _ := repo.GetByID(idB)
	assert.Equal(t, 1, a.Position)
	assert.Equal(t, 0, b.Position)
}
