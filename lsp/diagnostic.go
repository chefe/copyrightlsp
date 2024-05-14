package lsp

// DiagnosticSeverity defines if the dianostic is an error, warning, information or a hint.
type DiagnosticSeverity int

const (
	// DiagnosticSeverityError represents an error diagnostic.
	DiagnosticSeverityError DiagnosticSeverity = 1
	// DiagnosticSeverityWarning represents a warning diagnostic.
	DiagnosticSeverityWarning DiagnosticSeverity = 2
	// DiagnosticSeverityInformation represents an information diagnostic.
	DiagnosticSeverityInformation DiagnosticSeverity = 3
	// DiagnosticSeverityHint represents a hint diagnostic.
	DiagnosticSeverityHint DiagnosticSeverity = 4
)

// Diagnostic represents a diagnostic object of the lsp protocol.
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

// NewErrorDiagnostic creates an new error diagnostic with the given message.
func NewErrorDiagnostic(message string) Diagnostic {
	return Diagnostic{
		Message:  message,
		Source:   "copyrighlsp",
		Range:    NewRange(0, 0, 0, 0),
		Severity: DiagnosticSeverityError,
	}
}
