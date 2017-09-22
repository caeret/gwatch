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
command = "ls"
execute_at_ready = true
delay = 2
```

## Example usage

Executing `gwatch` with `gwatch.toml` in work dir:

```
gwatch
````

Executing `gwatch` with specified configuration file:

```
gwatch -c /path/to/config.toml
```

