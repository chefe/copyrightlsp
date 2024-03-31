package state

type documentInfo struct {
	Language string
	Content  string
}

type State struct {
	// Map file names to document infos
	Documents map[string]documentInfo
	// Map languages to template lines
	Templates map[string][]string
	// Map languages to search ranges
	SearchRanges map[string]uint8
}

func NewState() State {
	return State{
		Documents:    map[string]documentInfo{},
		Templates:    map[string][]string{},
		SearchRanges: map[string]uint8{},
	}
}

func (s *State) OpenDocument(document, text, language string) {
	s.Documents[document] = documentInfo{
		Language: language,
		Content:  text,
	}
}

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

func (s *State) CloseDocument(document string) {
	delete(s.Documents, document)
}

func (s *State) UpdateTemplates(templates map[string][]string) {
	s.Templates = templates
}

func (s *State) UpdateSearchRanges(searchRanges map[string]uint8) {
	s.SearchRanges = searchRanges
}

func (s *State) GetSearchRange(language string) uint8 {
	searchRange, ok := s.SearchRanges[language]
	if !ok {
		return 0
	}

	return searchRange
}
