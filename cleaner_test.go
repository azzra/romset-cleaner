package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestCannotCreateDestDir(t *testing.T) {

	dir, err := ioutil.TempDir("", "romset-cleaner")
	if err != nil {
		t.Errorf("Cannot TempDir %v: %v", dir, err)
	}

	flag.Set("dry_run", "false")
	flag.Set("dest_dir", dir+"non_existent")

	origMkdir := Mkdir
	defer func() { Mkdir = origMkdir }()
	Mkdir = func(name string, perm os.FileMode) error {
		return errors.New("cannot create dir")
	}

	origLogFatal := LogFatal
	defer func() { LogFatal = origLogFatal }()
	errors := []string{}
	LogFatal = func(args ...interface{}) {
		errors = append(errors, "error")
	}

	main()

	if count := len(errors); count != 1 {
		t.Errorf("Expected one error, actual %v", count)
	}

	flag.Set("dry_run", "true")
	flag.Set("dest_dir", "")

}

func TestCannotReadRomDir(t *testing.T) {

	origReadDir := ReadDir
	defer func() { ReadDir = origReadDir }()
	ReadDir = func(dirname string) ([]os.FileInfo, error) {
		return nil, errors.New("cannot read dir")
	}

	origLogFatal := LogFatal
	defer func() { LogFatal = origLogFatal }()
	errors := []string{}
	LogFatal = func(args ...interface{}) {
		errors = append(errors, "error")
	}

	main()

	if count := len(errors); count != 1 {
		t.Errorf("Expected one error, actual %v", count)
	}

}

func TestNormalizeFilename(t *testing.T) {

	const in, expectedFilename, expectedBasename = "filename [EUR] (Rev 1) [beta].tst", "filename (EUR) (Rev 1) (beta).tst", "filename"
	cleaned, basename, _ := normalizeFilename(in)

	if cleaned != expectedFilename {
		t.Errorf("normalizeFilename(%v) = %v, want %v", in, cleaned, expectedFilename)
	}

	if basename != expectedBasename {
		t.Errorf("normalizeFilename(%v) = %v, want %v", in, basename, expectedBasename)
	}
}

func TestNormalizeFilenameNil(t *testing.T) {

	names := []string{"FooBar.tst", "FooBar (bar.tst", "FooBar ().tst", "FooBar )ddsq(.tst"}

	for _, file := range names {
		_, _, err := normalizeFilename(file)

		if err == nil {
			t.Errorf("normalizeFilename(\"filename foobar\", %v) should return an error", file)
		}
	}

}

func TestContains(t *testing.T) {

	in := []string{"foo", "bar"}

	if contains("foo", in) == false {
		t.Errorf("contains(\"foo\", %v) should be true", in)
	}

	if contains("baz", in) == true {
		t.Errorf("contains(\"baz\", %v) should be false", in)
	}

}

func TestFindMatchingRomMatched(t *testing.T) {

	in := []Rom{Rom{"FooBar (foo)", []string{"foo"}}, Rom{"FooBar (bar)", []string{"bar"}}}

	if rom := findMatchingRom(in, []string{"bar"}); rom != &in[1] {
		t.Errorf("findMatchingRom(%v, [\"bar\"]) should be %v, %v instead", in, &in[1], rom)
	}

	if rom := findMatchingRom(in, []string{"foz", "foo"}); rom != &in[0] {
		t.Errorf("findMatchingRom(%v, [\"foz\", \"foo\"]) should be %v, %v instead", in, &in[0], rom)
	}

}

func TestFindMatchingRomNotMatched(t *testing.T) {

	in := []Rom{Rom{"FooBar (foo)", []string{"foo"}}, Rom{"FooBar (bar)", []string{"bar"}}}

	if rom := findMatchingRom(in, []string{"baz"}); rom != nil {
		t.Errorf("findMatchingRom(%v, [\"baz\"]) should be %v, %v instead", in, nil, nil)
	}

}

func TestExtractAttributes(t *testing.T) {

	in := "FooBar (foo) (oof,baz 1).tst"
	expected := []string{"foo", "oof", "baz1"}

	if attributes := extractAttributes(in); reflect.DeepEqual(attributes, expected) == false {
		t.Errorf("extractAttributes(%v) should be %v, %v instead", in, expected, attributes)
	}

	in = "FooBar (foo, bar) (  baz  ).tst"
	expected = []string{"foo", "bar", "baz"}

	if attributes := extractAttributes(in); reflect.DeepEqual(attributes, expected) == false {
		t.Errorf("extractAttributes(%v) should be %v, %v instead", in, expected, attributes)
	}

}

func setupMainFunc(t *testing.T) (string, []string) {

	dir, err := ioutil.TempDir("", "romset-cleaner")
	sampleRoms := []string{"FooBar (foo).tst", "FooBar (bar) (oof,baz) (boof 1).tst", "FooBar Foo Edition.tst", "BarFoo (one).tst"}

	if err != nil {
		t.Errorf("Cannot TempDir %v: %v", dir, err)
	}

	err = os.Mkdir(dir+"/roms", 0750)
	if err != nil {
		t.Errorf("Cannot Mkdir %v: %v", dir, err)
	}

	err = os.Mkdir(dir+"/barbaz", 0750)
	if err != nil {
		t.Errorf("Cannot Mkdir %v: %v", dir, err)
	}

	for _, file := range sampleRoms {
		if err := ioutil.WriteFile(dir+"/"+file, []byte("foofoo"), 0666); err != nil {
			t.Errorf("Cannot WriteFile %v: %v", dir, err)
		}
	}

	flag.Set("rom_dir", dir)
	flag.Set("dest_dir", "")
	flag.Set("keeped", "baz")

	return dir, sampleRoms

}

func TestMainFuncDefault(t *testing.T) {

	dir, sampleRoms := setupMainFunc(t)

	main()

	for _, file := range sampleRoms {
		expected := dir + "/" + file
		if _, err := os.Stat(expected); os.IsNotExist(err) {
			t.Errorf("File %v should exist", expected)
		}
	}

	if expected := *romDir + "/moved"; *destDir != expected {
		t.Errorf("Foriger %v should exist", expected)
	}

}

func TestMainFuncNotDryRun(t *testing.T) {

	dir, sampleRoms := setupMainFunc(t)

	flag.Set("dry_run", "false")

	main()

	expected := dir + "/moved/" + sampleRoms[1]
	if _, err := os.Stat(expected); os.IsNotExist(err) {
		t.Errorf("File %v should exist", expected)
	}

	expected = dir + "/" + sampleRoms[0]
	if _, err := os.Stat(expected); os.IsNotExist(err) {
		t.Errorf("File %v should exist", expected)
	}

}

func TestMainFuncDryRun(t *testing.T) {

	dir, sampleRoms := setupMainFunc(t)

	flag.Set("dry_run", "true")

	main()

	for _, file := range sampleRoms {
		expected := dir + "/" + file
		if _, err := os.Stat(expected); os.IsNotExist(err) {
			t.Errorf("File %v should exist", expected)
		}
	}

}

func TestMainFuncNotMatched(t *testing.T) {

	dir, sampleRoms := setupMainFunc(t)

	flag.Set("dry_run", "false")
	flag.Set("keeped", "zab")

	main()

	for _, file := range sampleRoms {
		expected := dir + "/" + file
		if _, err := os.Stat(expected); os.IsNotExist(err) {
			t.Errorf("File %v should exist", expected)
		}
	}

}

func TestMainFuncOnlyOneTrue(t *testing.T) {

	dir, sampleRoms := setupMainFunc(t)

	flag.Set("dry_run", "false")
	flag.Set("keeped", "foo")
	flag.Set("keep_one", "true")

	main()

	expected := dir + "/moved/" + sampleRoms[3]
	if _, err := os.Stat(expected); os.IsNotExist(err) {
		t.Errorf("File %v should exist", expected)
	}

	expected = dir + "/moved/" + sampleRoms[0]
	if _, err := os.Stat(expected); os.IsNotExist(err) {
		t.Errorf("File %v should exist", expected)
	}

}

func TestMainFuncOnlyOneFalse(t *testing.T) {

	dir, sampleRoms := setupMainFunc(t)

	flag.Set("dry_run", "false")
	flag.Set("keeped", "foo")
	flag.Set("keep_one", "false")

	main()

	expected := dir + "/" + sampleRoms[3]
	if _, err := os.Stat(expected); os.IsNotExist(err) {
		t.Errorf("File %v should exist", expected)
	}

	expected = dir + "/moved/" + sampleRoms[0]
	if _, err := os.Stat(expected); os.IsNotExist(err) {
		t.Errorf("File %v should exist", expected)
	}

}

func TestCannotRenameRom(t *testing.T) {

	setupMainFunc(t)

	flag.Set("dry_run", "false")
	flag.Set("keeped", "foo")

	origRename := Rename
	defer func() { Rename = origRename }()
	Rename = func(oldpath string, newpath string) error {
		return errors.New("cannot rename")
	}

	origLogFatal := LogFatal
	defer func() { LogFatal = origLogFatal }()
	errors := []string{}
	LogFatal = func(args ...interface{}) {
		errors = append(errors, "error")
	}

	main()

	if count := len(errors); count != 1 {
		t.Errorf("Expected one error, actual %v", count)
	}

}
