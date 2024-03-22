package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/chefe/copyrightlsp/lsp"
	"github.com/chefe/copyrightlsp/rpc"
)

func main() {
	logger := getLogger("/tmp/copyrightlsp.log")
	logger.Println("started copyrightlsp")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	writer := os.Stdout

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, content, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("got an error: %s\n", err)
			continue
		}

		handleMessage(logger, writer, method, content)
	}
}

func handleMessage(logger *log.Logger, writer io.Writer, method string, content []byte) {
	logger.Printf("recived message with method '%s'\n", method)

	switch method {
	case "initialize":
		handleInitializeMessage(logger, writer, content)
		logger.Println("sent the 'initialize' response")
	}
}

func handleInitializeMessage(logger *log.Logger, writer io.Writer, message []byte) {
	var request lsp.IntitializeRequest
	if err := json.Unmarshal(message, &request); err != nil {
		logger.Printf("recived invalid 'initialize' message: %s\n", err)
	}

	clientVersion := "UNKNOWN"
	if request.Params.ClientInfo.Version != nil {
		clientVersion = *request.Params.ClientInfo.Version
	}

	logger.Printf("connected to %s (Version: %s)\n", request.Params.ClientInfo.Name, clientVersion)
	replyMessage(logger, writer, lsp.NewInitializeResponse(request.ID))
}

func replyMessage(logger *log.Logger, writer io.Writer, message any) {
	reply := rpc.EncodeMessage(message)

	if _, err := writer.Write([]byte(reply)); err != nil {
		logger.Printf("failed to write response: %s", err)
	}
}

func getLogger(filename string) *log.Logger {
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("failed to open the log file")
	}

	return log.New(logFile, "[copyrightlsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
