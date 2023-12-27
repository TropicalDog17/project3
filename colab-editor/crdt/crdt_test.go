package crdt

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDocument(t *testing.T) {
	doc := New()

	require.Equal(t, doc.Length(), 2, "empty document should have length of two(characterStart and characterEnd)")
}

func TestInsert(t *testing.T) {
	doc := New()

	position := 1
	value := "a"

	content, err := doc.Insert(position, value)
	require.NoError(t, err)
	require.Equal(t, content, "a")

	expectedDoc := Document{
		Characters: []Character{
			{ID: "start", Visible: false, Value: "", IDPrevious: "", IDNext: "01"},
			{ID: "01", Visible: true, Value: "a", IDPrevious: "start", IDNext: "end"},
			{ID: "end", Visible: false, Value: "", IDPrevious: "01", IDNext: ""},
		},
	}
	require.Equal(t, doc, expectedDoc)
}

func TestIntegrateInsert_SamePosition(t *testing.T) {
	doc := &Document{
		Characters: []Character{
			{ID: "start", Visible: false, Value: "", IDPrevious: "", IDNext: "1"},
			{ID: "1", Visible: false, Value: "t", IDPrevious: "start", IDNext: "2"},
			{ID: "2", Visible: false, Value: "u", IDPrevious: "1", IDNext: "end"},
			{ID: "end", Visible: false, Value: "", IDPrevious: "2", IDNext: ""},
		},
	}
	// Insert a new character at the start. (IDPrevious = start)
	newChar := Character{ID: "3", Visible: false, Value: "b", IDPrevious: "start", IDNext: "1"}

	charPrev := Character{ID: "start", Visible: false, Value: "", IDPrevious: "", IDNext: "1"}
	charNext := Character{ID: "1", Visible: false, Value: "t", IDPrevious: "start", IDNext: "2"}

	content, err := doc.IntegrateInsert(newChar, charPrev, charNext)
	if err != nil {
		t.Errorf("error : %v \n", err)
	}
	expectedDoc := &Document{
		Characters: []Character{
			{ID: "start", Visible: false, Value: "", IDPrevious: "", IDNext: "3"},
			{ID: "3", Visible: false, Value: "b", IDPrevious: "start", IDNext: "1"},
			{ID: "1", Visible: false, Value: "t", IDPrevious: "3", IDNext: "2"},
			{ID: "2", Visible: false, Value: "u", IDPrevious: "1", IDNext: "end"},
			{ID: "end", Visible: false, Value: "", IDPrevious: "2", IDNext: ""},
		},
	}
	require.Equal(t, content, expectedDoc)
}

func TestIntegrateInsertAndDelete_Commutation(t *testing.T) {
	doc := &Document{
		Characters: []Character{
			{ID: "start", Visible: false, Value: "", IDPrevious: "", IDNext: "1"},
			{ID: "1", Visible: true, Value: "t", IDPrevious: "start", IDNext: "2"},
			{ID: "2", Visible: true, Value: "u", IDPrevious: "1", IDNext: "3"},
			{ID: "3", Visible: true, Value: "n", IDPrevious: "2", IDNext: "4"},
			{ID: "4", Visible: true, Value: "1", IDPrevious: "3", IDNext: "end"},
			{ID: "end", Visible: false, Value: "", IDPrevious: "4", IDNext: ""},
		},
	}
	newChar := Character{ID: "5", Visible: true, Value: "a", IDPrevious: "2", IDNext: "3"}
	charPrev := Character{ID: "2", Visible: true, Value: "u", IDPrevious: "1", IDNext: "2"}
	charNext := Character{ID: "3", Visible: true, Value: "n", IDPrevious: "2", IDNext: "4"}

	doc1, _ := doc.IntegrateInsert(newChar, charPrev, charNext)

	delChar := Character{ID: "4", Visible: true, Value: "1", IDPrevious: "3", IDNext: "end"}

	doc1 = doc1.IntegrateDelete(delChar)

	expectedDoc := &Document{
		Characters: []Character{
			{ID: "start", Visible: false, Value: "", IDPrevious: "", IDNext: "1"},
			{ID: "1", Visible: true, Value: "t", IDPrevious: "start", IDNext: "2"},
			{ID: "2", Visible: true, Value: "u", IDPrevious: "1", IDNext: "5"},
			{ID: "5", Visible: true, Value: "a", IDPrevious: "2", IDNext: "3"},
			{ID: "3", Visible: true, Value: "n", IDPrevious: "5", IDNext: "4"},
			{ID: "4", Visible: false, Value: "1", IDPrevious: "3", IDNext: "end"},

			{ID: "end", Visible: false, Value: "", IDPrevious: "4", IDNext: ""},
		},
	}
	require.Equal(t, expectedDoc, doc1)
	doc = &Document{
		Characters: []Character{
			{ID: "start", Visible: false, Value: "", IDPrevious: "", IDNext: "1"},
			{ID: "1", Visible: true, Value: "t", IDPrevious: "start", IDNext: "2"},
			{ID: "2", Visible: true, Value: "u", IDPrevious: "1", IDNext: "3"},
			{ID: "3", Visible: true, Value: "n", IDPrevious: "2", IDNext: "4"},
			{ID: "4", Visible: true, Value: "1", IDPrevious: "3", IDNext: "end"},
			{ID: "end", Visible: false, Value: "", IDPrevious: "4", IDNext: ""},
		},
	}
	require.Equal(t, 6, doc.Length())
	doc2 := doc.IntegrateDelete(delChar)
	expectedDoc2 := &Document{
		Characters: []Character{
			{ID: "start", Visible: false, Value: "", IDPrevious: "", IDNext: "1"},
			{ID: "1", Visible: true, Value: "t", IDPrevious: "start", IDNext: "2"},
			{ID: "2", Visible: true, Value: "u", IDPrevious: "1", IDNext: "3"},
			{ID: "3", Visible: true, Value: "n", IDPrevious: "2", IDNext: "4"},
			{ID: "4", Visible: false, Value: "1", IDPrevious: "3", IDNext: "end"},
			{ID: "end", Visible: false, Value: "", IDPrevious: "4", IDNext: ""},
		},
	}
	require.Equal(t, expectedDoc2, doc2)

	newChar = Character{ID: "5", Visible: true, Value: "a", IDPrevious: "2", IDNext: "3"}
	charPrev = Character{ID: "2", Visible: true, Value: "u", IDPrevious: "1", IDNext: "3"}
	charNext = Character{ID: "3", Visible: true, Value: "n", IDPrevious: "2", IDNext: "4"}
	doc2, _ = doc2.IntegrateInsert(newChar, charPrev, charNext)

	require.Equal(t, expectedDoc, doc2)

}
