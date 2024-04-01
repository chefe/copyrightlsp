// Provide support to encode and decode messages of the language server
// protocol. This implementation is based on the version 3.17 of the
// specification. See the link bellow for the complete specification.
//
// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/
package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func EncodeMessage(msg any) string {
	content, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content)
}

func parseContentLengthHeader(field []byte) (int, error) {
	if !strings.HasPrefix(string(field), "Content-Length: ") {
		return 0, errors.New("header did not start with content length")
	}

	contentLengthBytes := field[len("Content-Length: "):]

	contentLength, err := strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return 0, fmt.Errorf("faild to parse content length header: %w", err)
	}

	return contentLength, nil
}

func isValidContentTypeHeader(field []byte) bool {
	return string(field) == "Content-Type: application/vscode-jsonrpc; charset=utf-8"
}

func parseMessageHeader(header []byte) (int, error) {
	fieldOne, fieldTwo, _ := bytes.Cut(header, []byte{'\r', '\n'})
	if fieldTwo == nil {
		return parseContentLengthHeader(fieldOne)
	}

	if strings.HasPrefix(string(fieldOne), "Content-Type: ") {
		if !isValidContentTypeHeader(fieldOne) {
			return 0, errors.New("unsupported content type provided")
		}

		return parseContentLengthHeader(fieldTwo)
	}

	if strings.HasPrefix(string(fieldTwo), "Content-Type: ") {
		if !isValidContentTypeHeader(fieldTwo) {
			return 0, errors.New("unsupported content type provided")
		}

		return parseContentLengthHeader(fieldOne)
	}

	return 0, errors.New("invalid header fields provided")
}

func DecodeMessage(msg []byte) (string, []byte, error) {
	header, content, found := bytes.Cut(msg, []byte{'\r', '\n', '\r', '\n'})
	if !found {
		return "", nil, errors.New("did not find separator between header and content")
	}

	contentLength, err := parseMessageHeader(header)
	if err != nil {
		return "", nil, fmt.Errorf("faild to parse header: %w", err)
	}

	contentBytes := content[:contentLength]

	type BaseMessage struct {
		Method string `json:"method"`
	}

	var baseMessage BaseMessage
	if err := json.Unmarshal(contentBytes, &baseMessage); err != nil {
		return "", nil, fmt.Errorf("faild to unmarshal message: %w", err)
	}

	return baseMessage.Method, contentBytes, nil
}

func Split(data []byte, _ bool) (int, []byte, error) {
	header, content, found := bytes.Cut(data, []byte{'\r', '\n', '\r', '\n'})
	if !found {
		return 0, nil, nil
	}

	contentLength, err := parseMessageHeader(header)
	if err != nil {
		return 0, nil, err
	}

	if len(content) < contentLength {
		return 0, nil, nil
	}

	seperatorLength := 4
	totalLength := len(header) + seperatorLength + contentLength

	return totalLength, data[:totalLength], nil
}
