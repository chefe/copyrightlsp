package state

type State struct {
	// Map file names to contents
	Documents map[string]string
}

func NewState() State {
	return State{
		Documents: map[string]string{},
	}
}

func (s *State) OpenDocument(document, text string) {
	s.Documents[document] = text
}

func (s *State) UpdateDocument(document, text string) {
	s.Documents[document] = text
}

func (s *State) CloseDocument(document string) {
	delete(s.Documents, document)
}
