# gwatch

A tool for watching the filesystem changes and executing the specified command.

## Installation

```
go get -u github.com/gaemma/gwatch
cd $GOPATH/src/github.com/gaemma/gwatch
```

`make` will generate `gwatch` binary file under the work dir.

`make install` will install `gwatch` to the `$GOPATH/bin`.

## Configuration

The configuration file should be in `toml` format, for example `gwatch.toml`:

```toml
dir = "./demo"
excludes = ["./demo/.git", "./demo/.idea"]
recursive = true
command = "rsync src dest"
execute_at_ready = true
init = true
delay = 2
```

Command line flags:

```
Flags:
      --help                 Show context-sensitive help (also try --help-long and --help-man).
  -d, --dir=DIR              directory to watch.
  -e, --exclude=EXCLUDE ...  exclude diretory or files.
  -r, --recursive            recursively watch subdirectories.
  -i, --init                 execute the command immedialtely while the gwatch is ready.
      --delay=DELAY          duration delay to execute the command.
  -c, --config=CONFIG        config path.
  ```

## Example usage

Executing `gwatch` with custom configuration file:

```
gwatch -c /path/to/config.toml
```

If `-c` option is not specified, `gwatch` will check if `gwatch.toml` exists in work dir and use it.

