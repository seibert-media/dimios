package parser

import (
	"testing"

	"os"

	. "github.com/bborbe/assert"
	"github.com/bborbe/teamvault_utils/connector"
)

func TestParseContentWithoutPlaceholder(t *testing.T) {
	teamvaultConnector := connector.NewDummy()
	teamvaultParser := New(teamvaultConnector)
	contentWithoutPlaceholder := []byte("hello world")
	resultContent, err := teamvaultParser.Parse(contentWithoutPlaceholder)
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(resultContent, Is(contentWithoutPlaceholder)); err != nil {
		t.Fatal(err)
	}
}

func TestParseTeamvaultUsername(t *testing.T) {
	teamvaultConnector := connector.NewDummy()
	teamvaultParser := New(teamvaultConnector)
	resultContent, err := teamvaultParser.Parse([]byte(`{{ "asdf" | teamvaultUser }}`))
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(string(resultContent), Is("asdf")); err != nil {
		t.Fatal(err)
	}
}

func TestParseTeamvaultPassword(t *testing.T) {
	teamvaultConnector := connector.NewDummy()
	teamvaultParser := New(teamvaultConnector)
	resultContent, err := teamvaultParser.Parse([]byte(`{{ "asdf" | teamvaultPassword }}`))
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(string(resultContent), Is("6Jk10in-e7lYHEQubLMEW1MDb0fcFcw8t4aW5HEgvNI=")); err != nil {
		t.Fatal(err)
	}
}

func TestParseTeamvaultUrl(t *testing.T) {
	teamvaultConnector := connector.NewDummy()
	teamvaultParser := New(teamvaultConnector)
	resultContent, err := teamvaultParser.Parse([]byte(`{{ "asdf" | teamvaultUrl}}`))
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(string(resultContent), Is("XsLMyuFYK_HQTI1aoP1u0iX6UdYavwOdQoXINGeG4Ek=")); err != nil {
		t.Fatal(err)
	}
}

func TestParseTeamvaultFile(t *testing.T) {
	teamvaultConnector := connector.NewDummy()
	teamvaultParser := New(teamvaultConnector)
	resultContent, err := teamvaultParser.Parse([]byte(`{{ "asdf" | teamvaultFile}}`))
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(string(resultContent), Is("asdf-file")); err != nil {
		t.Fatal(err)
	}
}

func TestParseTeamvaultFileBase64(t *testing.T) {
	teamvaultConnector := connector.NewDummy()
	teamvaultParser := New(teamvaultConnector)
	resultContent, err := teamvaultParser.Parse([]byte(`{{ "asdf" | teamvaultFileBase64}}`))
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(string(resultContent), Is("YXNkZi1maWxl")); err != nil {
		t.Fatal(err)
	}
}

func TestParseBase64(t *testing.T) {
	teamvaultConnector := connector.NewDummy()
	teamvaultParser := New(teamvaultConnector)
	resultContent, err := teamvaultParser.Parse([]byte(`{{ "abc" | base64}}`))
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(string(resultContent), Is("YWJj")); err != nil {
		t.Fatal(err)
	}
}

func TestParseEnv(t *testing.T) {
	teamvaultConnector := connector.NewDummy()
	teamvaultParser := New(teamvaultConnector)
	os.Setenv("testEnv", "hello")
	resultContent, err := teamvaultParser.Parse([]byte(`{{ "testEnv" | env}}`))
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(string(resultContent), Is("hello")); err != nil {
		t.Fatal(err)
	}
}

func TestParseTeamvaultHtpasswd(t *testing.T) {
	teamvaultConnector := connector.NewDummy()
	teamvaultParser := New(teamvaultConnector)
	resultContent, err := teamvaultParser.Parse([]byte(`{{ "abc" | teamvaultHtpasswd}}`))
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(len(resultContent), Gt(0)); err != nil {
		t.Fatal(err)
	}
}
