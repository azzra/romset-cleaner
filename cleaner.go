package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

var keepedAttrs = flag.String("keeped", "french,france,fr,europe,eur,eu,english,en,eng,uk,word,usa,us", "The attributes you want to keep, in comma separated format.")
var romDir = flag.String("rom_dir", ".", "The directory containing the roms file to process.")
var destDir = flag.String("dest_dir", "", "The destination directory where the roms will be moved in, \"{rom_dir}/moved\" if empty.")
var dryRun = flag.Bool("dry_run", true, "Print what will be moved.")
var keepIfOnlyOne = flag.Bool("keep_one", false, "Move the file if it's the only one of its kind.")

// ReadDir exported to be mocked
var ReadDir = ioutil.ReadDir

// LogFatal exported to be mocked
var LogFatal = log.Fatal

// Mkdir exported to be mocked
var Mkdir = os.Mkdir

// Rename exported to be mocked
var Rename = os.Rename

// Rom is a file which has specific attribute(s) (zone/lang/..)
type Rom struct {
	filename   string
	attributes []string // usa / fr / rev 1 / ...
}

func normalizeFilename(filename string) (string, string, error) {
	cleaned := strings.Replace(filename, "[", "(", -1)
	cleaned = strings.Replace(cleaned, "]", ")", -1)

	separatorPos := strings.Index(cleaned, "(")
	if separatorPos == -1 || strings.LastIndex(cleaned, ")") <= separatorPos+1 {
		return "", "", errors.New("filename must contains both \"(\" & \")\", none found in " + cleaned)
	}

	basename := strings.TrimSpace(cleaned[0:separatorPos])

	return cleaned, basename, nil
}

func extractAttributes(filename string) []string {

	// anything that is a word/space/comma between parenthesis
	// it looks like there is no negative lookahead in go regex atm & even \(.*[^\)]\) doesn't match...
	re := regexp.MustCompile("\\((\\w+,)*\\w+\\)")
	matches := re.FindAllString(strings.ToLower(strings.Replace(filename, " ", "", -1)), -1)
	var attributes []string

	for _, group := range matches {
		values := strings.Split(group[1:len(group)-1], ",")
		attributes = append(attributes, values...)
	}

	return attributes
}

func contains(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func findMatchingRom(roms []Rom, attributes []string) *Rom {

	for _, attr := range attributes {

		// reverse walking, to have "USA" before "USA Rev1"
		for i := len(roms) - 1; i >= 0; i-- {
			if contains(attr, roms[i].attributes) {
				return &roms[i]
			}
		}

	}

	return nil
}

func main() {

	flag.Parse()
	if *destDir == "" {
		*destDir = strings.TrimRight(*romDir, "/") + "/moved"
	}

	fmt.Println("ROM DIR: " + *romDir + ", DEST DIR: " + *destDir)
	fmt.Println("KEEPED ATTRIBUTES: " + *keepedAttrs)

	dirFiles, err := ReadDir(*romDir)
	if err != nil {
		LogFatal(err)
	}

	// prepare destination directory
	if *dryRun == false {

		if _, err := os.Stat(*destDir); os.IsNotExist(err) {
			err := Mkdir(*destDir, 0750)
			if err != nil {
				LogFatal(err)
			}
		}

	}

	var roms = make(map[string][]Rom)

	// we construct an hashmap indexed by the base rom name and we put each file info in it
	for _, file := range dirFiles {

		if file.IsDir() {
			continue
		}

		filename := file.Name()
		cleanedFilename, baseFilename, err := normalizeFilename(filename)
		if err != nil {
			continue
		}

		roms[baseFilename] = append(roms[baseFilename], Rom{filename, extractAttributes(cleanedFilename)})

	}

	splittedKeepedAttrs := strings.Split(strings.Replace(*keepedAttrs, " ", "", -1), ",")

	// for each files (struct) for a rom, try to find a good one
	for game, romFiles := range roms {

		rom := findMatchingRom(romFiles[0:], splittedKeepedAttrs[0:])

		if *keepIfOnlyOne && len(romFiles) == 1 {
			rom = &romFiles[0]
		}

		if rom != nil {

			fmt.Println("OK: " + game + " - found: " + rom.filename)

			if *dryRun == false {
				err := Rename(*romDir+"/"+rom.filename, *destDir+"/"+rom.filename)
				if err != nil {
					LogFatal(err)
				}
			}

		} else {
			fmt.Println("KO: " + game)
		}

	}

}
