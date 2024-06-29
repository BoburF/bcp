package protocol

import (
	"bytes"
	"log"
	"net"
	"strings"
	"testing"
)

func Test_request_v1ConverTo(t *testing.T) {
	resource := "example"
	additions := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	data := strings.NewReader("sample data")
	request := &Request{
		Resource:  resource,
		Additions: additions,
		Data:      data,
	}

	server, client := net.Pipe()
	err := request.ConverTo(1, &server)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
  log.Println("Request is wrrtten")

	expectedOutput := "v1\x00" +
		"7\x00example\x00" +
		"4\x00key1\x006\x00value1" +
		"4\x00key2\x006\x00value2" +
		"\x00sample data"

	var buf bytes.Buffer
	if client != nil {
		_, err = buf.ReadFrom(client)
		if err != nil {
			t.Fatalf("failed to read from writer: %v", err)
		}
	}

	if buf.String() != expectedOutput {
		t.Errorf("expected: %q,\ngot: %q", expectedOutput, buf.String())
	}
	server.Close()
	client.Close()
}

func Test_request_invalidVersion(t *testing.T) {
	resource := "example"
	additions := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	data := strings.NewReader("sample data")
	request := &Request{
		Resource:  resource,
		Additions: additions,
		Data:      data,
	}

	server, client := net.Pipe()
	err := request.ConverTo(2, &server)
	if err == nil {
		t.Fatal("expected error for unsupported version, got none")
	}
	if err.Error() != "unsupported version" {
		t.Errorf("expected 'unsupported version' error, got %v", err)
	}
	server.Close()
	client.Close()
}
