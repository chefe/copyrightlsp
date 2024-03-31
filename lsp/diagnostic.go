package lsp

// Defines if the dianostic is an error, warning, information or a hint.
type DiagnosticSeverity int

const (
	// Reports an error.
	DiagnosticSeverityError DiagnosticSeverity = 1
	// Reports a warning.
	DiagnosticSeverityWarning DiagnosticSeverity = 2
	// Reports an information.
	DiagnosticSeverityInformation DiagnosticSeverity = 3
	// Reports a hint.
	DiagnosticSeverityHint DiagnosticSeverity = 4
)

type Diagnostic struct {
	// The diagnostic's message.
	Message string `json:"message"`
	// A human-readable string describing the source of this diagnostic, e.g.
	// 'typescript' or 'super lint'.
	Source string `json:"source"`
	// The range at which the message applies.
	Range Range `json:"range"`
	// The diagnostic's severity. Can be omitted. If omitted it is up to the
	// client to interpret diagnostics as error, warning, info or hint.
	Severity DiagnosticSeverity `json:"severity"`
}

func NewErrorDiagnostic(message string) Diagnostic {
	return Diagnostic{
		Message:  message,
		Source:   "copyrighlsp",
		Range:    NewRange(0, 0, 0, 0),
		Severity: DiagnosticSeverityError,
	}
}
