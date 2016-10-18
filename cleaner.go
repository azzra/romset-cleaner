package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

var keepedAttrs = flag.String("attrs", "french,fr,europe,eur,eu,english,en,eng,uk,word,usa,us", "The attributes you want to keep.")
var splittedKeepedAttrs []string
var romDir = flag.String("rom_dir", ".", "The directory containing the roms file to process.")
var dryRun = flag.Bool("dry_run", true, "Print what will moved.")

type Rom struct {
	filename   string
	attributes []string // usa / fr / rev 1 / ...
}

func normalizeFilename(filename string) (string, string) {
	cleaned := strings.Replace(filename, "[", "(", -1)
	cleaned = strings.Replace(cleaned, "]", ")", -1)

	separatorPos := strings.Index(cleaned, "(")
	basename := strings.TrimSpace(cleaned[0:separatorPos])

	return cleaned, basename
}

func extractAttributes(filename string) []string {

	// anything that is a word/space/comma between parenthesis, looks like there is no negative lookahead in go atm.
	re := regexp.MustCompile("\\((\\w+\\s?,?)+\\)")
	matches := re.FindAllString(strings.ToLower(filename), -1)
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

func findMatchingRom(roms []Rom) *Rom {

	for _, attr := range splittedKeepedAttrs {

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
	var movedDir string

	dirFiles, err := ioutil.ReadDir(*romDir)
	if err != nil {
		log.Fatal(err)
	}

	// prepare destination directory
	if *dryRun == false {

		movedDir = *romDir + "/moved"

		if _, err := os.Stat(movedDir); os.IsNotExist(err) {
			err := os.Mkdir(movedDir, 0750)
			if err != nil {
				log.Fatal(err)
			}
		}

	}

	fmt.Println("ROM DIR: " + *romDir)
	fmt.Println("KEEPED ATTRIBUTES: " + *keepedAttrs)

	var roms = make(map[string][]Rom)
	splittedKeepedAttrs = strings.Split(*keepedAttrs, ",")

	// we construct an hashmap indexed by the base rom name and we put each file info in it
	for _, file := range dirFiles {

		if file.IsDir() {
			continue
		}

		filename := file.Name()
		cleanedFilename, baseFilename := normalizeFilename(filename)

		roms[baseFilename] = append(roms[baseFilename], Rom{filename, extractAttributes(cleanedFilename)})

	}

	// for each files (struct) for a rom, try to find a good one
	for game, files := range roms {

		rom := findMatchingRom(files[0:])

		if rom != nil {

			fmt.Println("OK: " + game + " - found " + rom.filename)

			if *dryRun == false {
				err := os.Rename(*romDir+"/"+rom.filename, movedDir+"/"+rom.filename)
				if err != nil {
					log.Fatal(err)
				}
			}

		} else {
			fmt.Println("KO: " + game)
		}

	}

}
