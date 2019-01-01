package envus

import (
	"os"
	"testing"
)

var (
	testProperty1 *string
)

func TestOverrideEnvusFilename(t *testing.T) {
	refresh()
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Happy overridden filename",
			args: args{
				filename: "testfile.toml",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			OverrideEnvusFilename(tt.args.filename)
			if envusFile != tt.args.filename {
				t.Fail()
			}
		})
	}
}

func TestUseExternalEnvironment(t *testing.T) {
	refresh()
	tests := []struct {
		name string
	}{
		{
			name: "Using tier 2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			UseLocalEnvironment()
			if !useTier2 {
				t.Fail()
			}
		})
	}
}

func TestRegister(t *testing.T) {
	refresh()
	type args struct {
		key          string
		defaultValue string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Happy Register() with No Load()",
			args: args {
				key: "TEST_1",
				defaultValue: "Hello World",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testProperty1 = bareboneStringPointer()
			Register(tt.args.key, testProperty1, tt.args.defaultValue)
			if *testProperty1 != tt.args.defaultValue {
				t.Fail()
			}
		})
	}
}

func TestFlow(t *testing.T) {
	refresh()

	// Prep local environment
	os.Setenv("URL", "https://test.com")

	var username string
	var password string
	var url string
	var token string

	Register("USERNAME", &username, "testuser")
	Register("PASSWORD", &password, "testpassword")
	Register("URL", &url, "testurl")
	Register("TOKEN", &token, "xyz123")

	UseLocalEnvironment()
	OverrideEnvusFilename("test.toml")
	Load()

	if username != "jabgibson" {
		t.Log("incorrect username: "+username)
		t.Fail()
	}

	if password != "abc123" {
		t.Log("incorrect password: "+password)
		t.Fail()
	}

	if url != "https://test.com" {
		t.Log("incorrect url: "+url)
		t.Fail()
	}

	if token != "xyz123" {
		t.Log("incorrect token: "+token)
		t.Fail()
	}

	if Fetch("secret") != "Rehab the beast" {
		t.Log("failure to find non registered property")
		t.Fail()
	}
}

func TestWithoutLocalEnvironment(t *testing.T) {
	refresh()

	// Prep local environment
	os.Setenv("USERNAME", "local-username")

	var username string

	Register("USERNAME", &username, "default-username")

	OverrideEnvusFilename("non-existent-file")
	Load()

	if username != "default-username" {
		t.Log("incorrect username: "+username)
		t.Fail()
	}
}

func TestUsingJsonFile(t *testing.T) {
	refresh()

	var username string

	Register("USERNAME", &username, "default-username")

	OverrideEnvusFilename("test.json")
	Load()

	if username != "jabgibson" {
		t.Log("incorrect username: "+username)
		t.Fail()
	}
}

func bareboneStringPointer() *string {
	x := "BAREBONE STRING POINTER"
	return &x
}
