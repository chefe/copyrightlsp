package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"

	"github.com/chefe/copyrightlsp/codeactions"
	"github.com/chefe/copyrightlsp/diagnostics"
	"github.com/chefe/copyrightlsp/lsp"
	"github.com/chefe/copyrightlsp/rpc"
	"github.com/chefe/copyrightlsp/state"
)

func main() {
	logFile := flag.String("logFile", "", "log message to the given file")
	flag.Parse()

	logWriter := io.Discard
	if len(*logFile) > 0 {
		logWriter = createLogFileWriter(*logFile)
	}

	logger := log.New(logWriter, "[copyrightlsp]", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Println("started copyrightlsp")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	writer := os.Stdout

	state := state.NewState()

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, content, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("got an error: %s\n", err)
			continue
		}

		if method == "exit" {
			logger.Println("received the 'exit' request")
			return
		}

		handleMessage(logger, writer, &state, method, content)
	}
}

func handleMessage(logger *log.Logger, writer io.Writer, state *state.State, method string, content []byte) {
	logger.Printf("received message with method '%s'\n", method)

	switch method {
	case "initialize":
		handleInitializeMessage(logger, writer, content)
	case "shutdown":
		handleShutdownMessage(logger, writer, content)
	case "textDocument/didOpen":
		handleTextDocumentDidOpenMessage(logger, state, writer, content)
	case "textDocument/didChange":
		handleTextDocumentDidChangeMessage(logger, state, writer, content)
	case "textDocument/didClose":
		handleTextDocumentDidCloseMessage(logger, state, content)
	case "textDocument/codeAction":
		handleTextDocumentCodeActionMessage(logger, writer, state, content)
	case "workspace/didChangeConfiguration":
		handleWorkspaceDidChangeConfigurationMessage(logger, state, content)
	}
}

func handleInitializeMessage(logger *log.Logger, writer io.Writer, message []byte) {
	var request lsp.IntitializeRequest
	if err := json.Unmarshal(message, &request); err != nil {
		logger.Printf("received invalid 'initialize' message: %s\n", err)
	}

	clientVersion := "UNKNOWN"
	if request.Params.ClientInfo.Version != nil {
		clientVersion = *request.Params.ClientInfo.Version
	}

	logger.Printf("connected to %s (Version: %s)\n", request.Params.ClientInfo.Name, clientVersion)
	replyMessage(logger, writer, lsp.NewInitializeResponse(request.ID))
	logger.Println("sent the 'initialize' response")
}

func handleShutdownMessage(logger *log.Logger, writer io.Writer, message []byte) {
	var request lsp.ShutdownRequest
	if err := json.Unmarshal(message, &request); err != nil {
		logger.Printf("received invalid 'shutdown' message: %s\n", err)
	}

	logger.Println("shudown")
	replyMessage(logger, writer, lsp.NewShudownResponse(request.ID))
	logger.Println("sent the 'shutdown' response")
}

func handleTextDocumentDidOpenMessage(logger *log.Logger, state *state.State, writer io.Writer, message []byte) {
	var request lsp.DidOpenTextDocumentNotification
	if err := json.Unmarshal(message, &request); err != nil {
		logger.Printf("received invalid 'textDocument/didOpen' message: %s\n", err)
	}

	state.OpenDocument(request.Params.TextDocument.URI, request.Params.TextDocument.Text, request.Params.TextDocument.LanguageID)
	logger.Printf("opend document %s [%s]\n", request.Params.TextDocument.URI, request.Params.TextDocument.LanguageID)

	diag := diagnostics.CalculateDiagnostics(state, request.Params.TextDocument.URI)
	replyMessage(logger, writer, lsp.NewPublishDiagnosticsNotification(request.Params.TextDocument.URI, diag))
	logger.Printf("calculated document diagnostics %s [%d]\n", request.Params.TextDocument.URI, len(diag))
}

func handleTextDocumentDidChangeMessage(logger *log.Logger, state *state.State, writer io.Writer, message []byte) {
	var request lsp.DidChangeTextDocumentNotification
	if err := json.Unmarshal(message, &request); err != nil {
		logger.Printf("received invalid 'textDocument/didChange' message: %s\n", err)
	}

	for _, change := range request.Params.ContentChanges {
		state.UpdateDocument(request.Params.TextDocument.URI, change.Text)
	}

	logger.Printf("changed document %s\n", request.Params.TextDocument.URI)

	diag := diagnostics.CalculateDiagnostics(state, request.Params.TextDocument.URI)
	replyMessage(logger, writer, lsp.NewPublishDiagnosticsNotification(request.Params.TextDocument.URI, diag))
	logger.Printf("calculated document diagnostics %s [%d]\n", request.Params.TextDocument.URI, len(diag))
}

func handleTextDocumentDidCloseMessage(logger *log.Logger, state *state.State, message []byte) {
	var request lsp.DidCloseTextDocumentNotification
	if err := json.Unmarshal(message, &request); err != nil {
		logger.Printf("received invalid 'textDocument/didClose' message: %s\n", err)
	}

	state.CloseDocument(request.Params.TextDocument.URI)
	logger.Printf("closed document %s\n", request.Params.TextDocument.URI)
}

func handleTextDocumentCodeActionMessage(logger *log.Logger, writer io.Writer, state *state.State, message []byte) {
	var request lsp.CodeActionRequest
	if err := json.Unmarshal(message, &request); err != nil {
		logger.Printf("received invalid 'textDocument/codeAction' message: %s\n", err)
	}

	actions := codeactions.CalculateCodeActions(state, request.Params.TextDocument.URI, request.Params.Range.Start, request.Params.Range.End)
	logger.Printf("calculated %d code actions for %s\n", len(actions), request.Params.TextDocument.URI)
	replyMessage(logger, writer, lsp.NewCodeActionResponse(request.ID, actions))
}

func handleWorkspaceDidChangeConfigurationMessage(logger *log.Logger, state *state.State, message []byte) {
	var request lsp.DidChangeConfigurationNotification
	if err := json.Unmarshal(message, &request); err != nil {
		logger.Printf("received invalid 'workspace/didChangeConfiguration' message: %s\n", err)
	}

	state.UpdateTemplates(request.Params.Settings.Templates)
	state.UpdateSearchRanges(request.Params.Settings.SearchRanges)
	logger.Println("updated settings")
}

func replyMessage(logger *log.Logger, writer io.Writer, message any) {
	reply := rpc.EncodeMessage(message)

	if _, err := writer.Write([]byte(reply)); err != nil {
		logger.Printf("failed to write response: %s", err)
	}
}

func createLogFileWriter(filename string) io.Writer {
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o666)
	if err != nil {
		panic("failed to open or create the log file")
	}

	return logFile
}
