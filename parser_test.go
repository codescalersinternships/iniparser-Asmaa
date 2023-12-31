package main

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

var validData = `
;server section
[server]
ip = 127.0.0.1
port = 8080

;database section
[database]
host = localhost
port = 5432
name = mydb`

var invalidSectionNameData = `
server]
ip = 127.0.0.1
port = 8080

[database]
host = localhost
port = 5432
name = mydb`

var invalidKeyData = `
[server]
= 127.0.0.1
port = 8080

[database]
host = localhost
port = 5432
name = mydb`

var sectionNameEmptyData = `
[]
ip = 127.0.0.1
port = 8080

[database]
host = localhost
port = 5432
name = mydb`

var redefiningKeyData = `
[server]
host = 127.0.0.1
host = 127.0.1.1

[database]
host = localhost`

var configData = Data{
	"server": {
		"ip":   "127.0.0.1",
		"port": "8080",
	},
	"database": {
		"host": "localhost",
		"port": "5432",
		"name": "mydb",
	},
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func TestLoadFromFile(t *testing.T) {
	want := configData

	ini := NewINIParser()

	dir := t.TempDir()
	filePath := filepath.Join(dir, "config.ini")

	err := os.WriteFile(filePath, []byte(validData), 0644)

	if err != nil {
		t.Fatalf("error creating temp file: %v", err)
	}

	err = ini.LoadFromFile(filePath)
	if err == ErrorFileExtension {
		t.Fatalf("Error: invalid file name")
	}

	got := ini.sections

	if !reflect.DeepEqual(got, want) {
		t.Errorf("config does not match expected config.\nExpected: %+v\nActual: %+v", want, got)
	}

	want = Data{
		"database": {
			"host": "localhost",
			"port": "5432",
			"name": "mydb",
		},
	}

	if reflect.DeepEqual(got, want) {
		t.Errorf("error in config")
	}

	err = ini.LoadFromFile("testt.ini")
	if !(err == ErrorOpeningFile) {
		t.Errorf("error in file")
	}

	dir = t.TempDir()
	filePath = filepath.Join(dir, "config.txt")

	err = os.WriteFile(filePath, []byte(validData), 0644)

	if err != nil {
		t.Fatalf("error creating temp file: %v", err)
	}

	err = ini.LoadFromFile(filePath)
	if err != ErrorFileExtension {
		t.Fatalf("Error: invalid file name")
	}
}

func TestLoadFromString(t *testing.T) {
	want := configData

	ini := NewINIParser()
	err := ini.LoadFromString(validData)
	if err != nil {
		t.Errorf("Error:%v", err)
	}
	got := ini.sections

	if !reflect.DeepEqual(got, want) {
		t.Errorf("config does not match expected config.\nExpected: %+v\nActual: %+v", want, got)
	}

	err = ini.LoadFromString(invalidSectionNameData)
	if !(err == ErrorInvalidFormat) {
		t.Errorf("error invalid format")
	}

	err = ini.LoadFromString(invalidKeyData)
	if !(err == ErrorInvalidKeyFormat) {
		t.Errorf("error invalid key format")
	}

	err = ini.LoadFromString(sectionNameEmptyData)
	if !(err == ErrorSectionNameEmpty) {
		t.Errorf("error invalid section empty")
	}

	err = ini.LoadFromString(redefiningKeyData)
	if !(err == ErrorRedefiningKey) {
		t.Errorf("error redifinition key")
	}
}

func TestGetSections(t *testing.T) {
	want := configData

	ini := NewINIParser()

	dir := t.TempDir()
	filePath := filepath.Join(dir, "config.ini")

	err := os.WriteFile(filePath, []byte(validData), 0644)

	if err != nil {
		t.Fatalf("error creating temp file: %v", err)
	}

	err = ini.LoadFromFile(filePath)
	if err != nil {
		t.Errorf("Error:%v", err)
	}

	got := ini.GetSections()

	// Compare the parsed config with the expected config
	if !reflect.DeepEqual(got, want) {
		t.Errorf("config does not match expected config.\nExpected: %+v\nActual: %+v", want, got)
	}
}

func TestGetSectionNames(t *testing.T) {

	want := []string{"server", "database"}

	ini := NewINIParser()

	dir := t.TempDir()
	filePath := filepath.Join(dir, "config.ini")

	err := os.WriteFile(filePath, []byte(validData), 0644)

	if err != nil {
		t.Fatalf("error creating temp file: %v", err)
	}

	err = ini.LoadFromFile(filePath)
	if err != nil {
		t.Errorf("Error:%v", err)
	}

	got := ini.GetSectionNames()

	if err != nil {
		t.Fatalf("Error: %v", err)
		return
	}

	if !(contains(want, got[0]) || contains(want, got[1])) {
		t.Errorf("section names value don't match expected value.\nExpected: %+v\nActual: %+v", want, got)
	}
}

func TestGet(t *testing.T) {

	want := "8080"
	ini := NewINIParser()

	dir := t.TempDir()
	filePath := filepath.Join(dir, "config.ini")

	err := os.WriteFile(filePath, []byte(validData), 0644)

	if err != nil {
		t.Fatalf("error creating temp file: %v", err)
	}

	err = ini.LoadFromFile(filePath)
	if err != nil {
		t.Errorf("Error:%v", err)
	}

	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	got, err := ini.Get("server", "port")
	if err != nil {
		t.Fatalf("Error: %v", err)
		return
	}

	if !(got == want) {
		t.Errorf("Reading value does not match expected value.\nExpected: %+v\nActual: %+v", want, got)
	}

	_, err = ini.Get("serve", "port")
	if !(err == ErrorSectionNotFound) {
		t.Errorf("wrong section name")
	}

	_, err = ini.Get("server", "portt")
	if !(err == ErrorKeyName) {
		t.Errorf("wrong key name")
	}
}

func TestSet(t *testing.T) {

	want := "8000"

	ini := NewINIParser()

	dir := t.TempDir()
	filePath := filepath.Join(dir, "config.ini")

	err := os.WriteFile(filePath, []byte(validData), 0644)

	if err != nil {
		t.Fatalf("error creating temp file: %v", err)
	}

	err = ini.LoadFromFile(filePath)
	if err != nil {
		t.Errorf("Error:%v", err)
	}

	ini.Set("database", "port", "8000")

	got := ini.sections["database"]["port"]

	if !(got == want) {
		t.Errorf("setting value does not match expected value.\nExpected: %+v\nActual: %+v", want, got)
	}

	ini.Set("database", "portt", "8000")
	got = ini.sections["database"]["portt"]

	if !(got == want) {
		t.Errorf("setting value does not match expected value.\nExpected: %+v\nActual: %+v", want, got)
	}

	ini.Set("databasee", "port", "8000")
	got = ini.sections["databasee"]["port"]

	if !(got == want) {
		t.Errorf("setting value does not match expected value.\nExpected: %+v\nActual: %+v", want, got)
	}
}

func TestString(t *testing.T) {
	want := configData

	ini := NewINIParser()

	err := ini.LoadFromString(validData)
	if err != nil {
		t.Errorf("Error:%v", err)
	}

	got := ini.String()

	// Compare the parsed config with the expected config
	if !(strings.Contains(got, "[server]") || strings.Contains(got, "port = 8080") || strings.Contains(got, "[database]") || strings.Contains(got, "host = localhost")) {
		t.Errorf("config does not match expected config.\nExpected: %+v\nActual: %+v", want, got)
	}

	// Compare the parsed config with the expected config
	if strings.Contains(got, ";server section") {
		t.Errorf("config does not match expected config.\nExpected: %+v\nActual: %+v", want, got)
	}
}

func TestSaveToFile(t *testing.T) {

	ini := NewINIParser()

	dir := t.TempDir()
	filePath := filepath.Join(dir, "config.ini")

	err := os.WriteFile(filePath, []byte(validData), 0644)

	if err != nil {
		t.Fatalf("error creating temp file: %v", err)
	}

	err = ini.LoadFromFile(filePath)
	if err != nil {
		t.Errorf("Error:%v", err)
	}

	got := ini.SaveToFile("false.txt")

	if !(got == ErrorFileExtension) {
		t.Errorf("wrong file name")
	}

	got = ini.SaveToFile("true.ini")
	if got == ErrorFileExtension {
		t.Errorf("wrong file name")
	}

	os.Remove("true.ini")
}
