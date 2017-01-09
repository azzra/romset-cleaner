[![Build Status](https://travis-ci.org/azzra/romset-cleaner.png)](https://travis-ci.org/azzra/romset-cleaner)
[![Coverage Status](https://coveralls.io/repos/github/azzra/romset-cleaner/badge.svg?branch=coverage)](https://coveralls.io/github/azzra/romset-cleaner?branch=coverage)
[![Code Climate](https://codeclimate.com/github/azzra/romset-cleaner/badges/gpa.svg)](https://codeclimate.com/github/azzra/romset-cleaner)

# Romset Cleaner

Move your selected roms from all the romset.

From
```
./
	Double Dragon (UE).zip
	Double Dragon (USA) (Beta 1).zip
	Double Dragon (Prototype).zip
```

To
```
./
	Double Dragon (USA) (Beta 1).zip
	Double Dragon (Prototype).zip
./moved
	Double Dragon (UE).zip
```


## Usage

```sh
-dest_dir string
    The destination directory where the roms will be moved in, "{rom_dir}/moved" if empty.

-dry_run
    Print what will moved. (default true)

-keep_one
    Move the file if it's the only one of its kind.

-keeped string
    The attributes you want to keep, in comma separated format. (default "french,france,fr,europe,eur,eu,english,en,eng,uk,word,usa,us")

-rom_dir string
    The directory containing the roms file to process. (default ".")
```


## Local Build and Test

You can use go get command: 

    go get github.com/azzra/romset-cleaner

Testing:

    go test github.com/azzra/romset-cleaner/...
    # or
    go test -tags noasm -coverprofile=coverage.out  -v ./...  
