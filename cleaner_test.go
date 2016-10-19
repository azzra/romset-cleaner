package main

import (
	"flag"
	"io/ioutil"
	"os"	
	"reflect"
	"testing"
)

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

	err = os.Mkdir(dir + "/roms", 0750)
	if err != nil {
		t.Errorf("Cannot Mkdir %v: %v", dir, err)
	}

	err = os.Mkdir(dir + "/roms/barbaz", 0750)
	if err != nil {
		t.Errorf("Cannot Mkdir %v: %v", dir, err)
	}

	for _, file := range sampleRoms {	
		if err := ioutil.WriteFile(dir + "/roms/" + file, []byte("foofoo"), 0666); err != nil {
			t.Errorf("Cannot WriteFile %v: %v", dir, err)
		}
	}


	flag.Set("rom_dir", dir + "/roms")

	return dir, sampleRoms

}

func TestMainFuncNotDryRun(t *testing.T) {

	dir, sampleRoms := setupMainFunc(t)

	flag.Set("keeped", "baz")
	flag.Set("dry_run", "false")

	main()

	expected := dir + "/roms/moved/" + sampleRoms[1]
	if _, err := os.Stat(expected); os.IsNotExist(err) {
		t.Errorf("File %v should exist", expected)
	}


	expected = dir + "/roms/" +  sampleRoms[0]
	if _, err := os.Stat(expected); os.IsNotExist(err) {
		t.Errorf("File %v should exist", expected)
	}

}


func TestMainFuncDryRun(t *testing.T) {

	dir, sampleRoms := setupMainFunc(t)


	flag.Set("keeped", "baz")
	flag.Set("dry_run", "true")

	main()

	for _, file := range sampleRoms {	
		expected := dir + "/roms/" + file
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
		expected := dir + "/roms/" + file
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

	expected := dir + "/roms/moved/" + sampleRoms[3]
	if _, err := os.Stat(expected); os.IsNotExist(err) {
		t.Errorf("File %v should exist", expected)
	}

	expected = dir + "/roms/moved/" +  sampleRoms[0]
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

	expected := dir + "/roms/" + sampleRoms[3]
	if _, err := os.Stat(expected); os.IsNotExist(err) {
		t.Errorf("File %v should exist", expected)
	}

	expected = dir + "/roms/moved/" +  sampleRoms[0]
	if _, err := os.Stat(expected); os.IsNotExist(err) {
		t.Errorf("File %v should exist", expected)
	}

}

