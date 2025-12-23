A TUI in Golang to modify Trails From Zero save file's difficulty.

# Build Steps

Download the repo with `git clone https://github.com/Shuddown/ChangeDifficultyTrailsFromZero.git` and run 
```
go mod tidy
go build .
```

# Usage

## Windows and Linux

Run the executable by providing the filepath of your savefile as the sole argument. Make sure to back up your save before trying this!
Note that this only works with NIS's official PC release of Trails from Zero. The Geofront Translation is NOT supported.

## macOS

I don't pay for an Apple Developer Account, so you have to run 

```
xattr -d com.apple.quarantine change-difficulty-macos-universal
```

before using the binary
