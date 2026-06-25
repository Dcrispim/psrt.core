package visualapp

import "github.com/Dcrispim/psrt.core/psrt"

type undoStack struct {
	items []psrt.Document
	limit int
}

func newUndoStack(limit int) *undoStack {
	return &undoStack{limit: limit}
}

func (s *undoStack) push(doc psrt.Document) {
	s.items = append(s.items, cloneDoc(doc))
	if len(s.items) > s.limit {
		s.items = s.items[1:]
	}
}

func (s *undoStack) pop() (psrt.Document, bool) {
	if len(s.items) == 0 {
		return psrt.Document{}, false
	}
	i := len(s.items) - 1
	doc := s.items[i]
	s.items = s.items[:i]
	return doc, true
}

func cloneDoc(doc psrt.Document) psrt.Document {
	b, err := psrt.ToJSON(doc)
	if err != nil {
		return doc
	}
	out, err := psrt.ParseJSON(b)
	if err != nil {
		return doc
	}
	return out
}
