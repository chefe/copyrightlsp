package rpc_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/chefe/copyrightlsp/rpc"
)

type EncodingExample struct {
	Testing bool `json:"testing"`
}

func TestEncodeMessage(t *testing.T) {
	t.Parallel()

	expected := "Content-Length: 16\r\n\r\n{\"testing\":true}"
	got := rpc.EncodeMessage(EncodingExample{Testing: true})

	if expected != got {
		t.Fatalf("expected: %s, actual: %s", expected, got)
	}
}

func TestDecodeMessage(t *testing.T) {
	t.Parallel()

	type want struct {
		method        string
		content       string
		contentLength int
	}

	tests := []struct {
		name  string
		lines []string
		want  want
	}{
		{
			name:  "only content length",
			lines: []string{"Content-Length: 17", "", "{\"method\":\"test\"}"},
			want: want{
				method:        "test",
				content:       "{\"method\":\"test\"}",
				contentLength: 17,
			},
		},
		{
			name: "content length first and then content type",
			lines: []string{
				"Content-Length: 16",
				"Content-Type: application/vscode-jsonrpc; charset=utf-8",
				"",
				"{\"method\":\"one\"}",
			},
			want: want{
				method:        "one",
				content:       "{\"method\":\"one\"}",
				contentLength: 16,
			},
		},
		{
			name: "content type first and then content length",
			lines: []string{
				"Content-Type: application/vscode-jsonrpc; charset=utf-8",
				"Content-Length: 22",
				"",
				"{\"method\":\"something\"}",
			},
			want: want{
				method:        "something",
				content:       "{\"method\":\"something\"}",
				contentLength: 22,
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			message := []byte(strings.Join(tt.lines, "\r\n"))

			method, content, err := rpc.DecodeMessage(message)
			if err != nil {
				t.Fatalf("expected no error but got one: %s", err)
			}

			if method != tt.want.method {
				t.Fatalf("method expected: '%s', got: '%s'", tt.want.method, method)
			}

			contentLength := len(content)
			if contentLength != tt.want.contentLength {
				t.Fatalf("content length expected: %d, got: %d", tt.want.contentLength, contentLength)
			}

			if string(content) != tt.want.content {
				t.Fatalf("content expected: '%s', got: '%s'", tt.want.content, string(content))
			}
		})
	}
}

func TestDecodeMessageWithError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		message []byte
	}{
		{
			name:    "nil",
			message: nil,
		},
		{
			name:    "empty slice",
			message: []byte{},
		},
		{
			name:    "no separator",
			message: []byte("some\r\ndata"),
		},
		{
			name:    "no content length header",
			message: []byte("header\r\n\r\ncontent"),
		},
		{
			name:    "invalid content length header",
			message: []byte("Content-Length: abc\r\n\r\ncontent"),
		},
		{
			name:    "content length header and invalid data",
			message: []byte("Content-Length: abc\r\nheader\r\n\r\ncontent"),
		},
		{
			name:    "invalid content type header",
			message: []byte("Content-Length: abc\r\nContent-Type: text/plain\r\n\r\ncontent"),
		},
		{
			name:    "content length header, content type and invalid data",
			message: []byte("Content-Length: abc\r\nContent-Type: application/vscode-jsonrpc; charset=utf-8\r\nheader\r\n\r\ncontent"),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			method, content, err := rpc.DecodeMessage(tt.message)
			if err == nil {
				t.Fatal("expected error but got none")
			}

			if method != "" {
				t.Fatalf("expected empty string for method but got: '%s'", method)
			}

			if content != nil {
				t.Fatalf("expected nil for content but got: '%s'", string(content))
			}
		})
	}
}

func TestSplit(t *testing.T) {
	t.Parallel()

	type want struct {
		advance int
		message []byte
		err     bool
	}

	tests := []struct {
		name  string
		input []byte
		want  want
	}{
		{
			name:  "nil",
			input: nil,
			want:  want{advance: 0, message: nil, err: false},
		},
		{
			name:  "empty slice",
			input: []byte{},
			want:  want{advance: 0, message: nil, err: false},
		},
		{
			name:  "incomplete header",
			input: []byte("Content-Length"),
			want:  want{advance: 0, message: nil, err: false},
		},
		{
			name:  "complete header but no content",
			input: []byte("Content-Length: 17\r\n\r\n"),
			want:  want{advance: 0, message: nil, err: false},
		},
		{
			name:  "complete header but invalid",
			input: []byte("Content-Length: asd\r\n\r\n"),
			want:  want{advance: 0, message: nil, err: true},
		},
		{
			name:  "complete message",
			input: []byte("Content-Length: 17\r\n\r\n{\"method\":\"demo\"}"),
			want: want{
				advance: 39,
				message: []byte("Content-Length: 17\r\n\r\n{\"method\":\"demo\"}"),
				err:     false,
			},
		},
		{
			name:  "complete message and more",
			input: []byte("Content-Length: 16\r\n\r\n{\"method\":\"nix\"}and some more"),
			want: want{
				advance: 38,
				message: []byte("Content-Length: 16\r\n\r\n{\"method\":\"nix\"}"),
				err:     false,
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			advance, message, err := rpc.Split(tt.input, false)
			if tt.want.err && err == nil {
				t.Fatal("expected error but got none")
			}

			if !tt.want.err && err != nil {
				t.Fatalf("expected no error but got '%s'", err)
			}

			if advance != tt.want.advance {
				t.Fatalf("advance expected: %d, got: '%d'", tt.want.advance, advance)
			}

			if !bytes.Equal(message, tt.want.message) {
				t.Fatalf("message expected: '%s', got: '%s'", tt.want.message, message)
			}
		})
	}
}
