package state

type documentInfo struct {
	Language string
	Content  string
}

// State holds all the state of the lsp server.
type State struct {
	// Map file names to document infos
	Documents map[string]documentInfo
	// Map languages to template lines
	Templates map[string][]string
	// Map languages to search ranges
	SearchRanges map[string]uint8
}

// NewState creates a new instace of the `State` struct.
func NewState() State {
	return State{
		Documents:    map[string]documentInfo{},
		Templates:    map[string][]string{},
		SearchRanges: map[string]uint8{},
	}
}

// OpenDocument adds the given document to the state.
func (s *State) OpenDocument(document, text, language string) {
	s.Documents[document] = documentInfo{
		Language: language,
		Content:  text,
	}
}

// UpdateDocument updates the given document in the state.
func (s *State) UpdateDocument(document, text string) {
	doc, ok := s.Documents[document]
	if !ok {
		return
	}

	s.Documents[document] = documentInfo{
		Language: doc.Language,
		Content:  text,
	}
}

// CloseDocument removes the given document from the state.
func (s *State) CloseDocument(document string) {
	delete(s.Documents, document)
}

// UpdateTemplates updates the mapping table of language to templates.
func (s *State) UpdateTemplates(templates map[string][]string) {
	s.Templates = templates
}

// UpdateSearchRanges updates the mapping table of language to search ranges.
func (s *State) UpdateSearchRanges(searchRanges map[string]uint8) {
	s.SearchRanges = searchRanges
}

// GetSearchRange returns the configured search range for the given language.
func (s *State) GetSearchRange(language string) uint8 {
	searchRange, ok := s.SearchRanges[language]
	if !ok {
		return 0
	}

	return searchRange
}
