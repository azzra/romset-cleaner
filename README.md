[![Build Status](https://travis-ci.org/azzra/romset-cleaner.png)](https://travis-ci.org/azzra/romset-cleaner)
[![Coverage Status](https://coveralls.io/repos/github/azzra/romset-cleaner/badge.svg?branch=coverage)](https://coveralls.io/github/azzra/romset-cleaner?branch=coverage)

# Romset Cleaner

Move your selected roms from all the romset.

## Usage

```sh
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

