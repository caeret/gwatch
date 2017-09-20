package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/rjeczalik/notify"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Config struct {
	Dir            string
	Excludes       []string
	Recursive      bool
	Command        string
	ExecuteAtReady bool `toml:"execute_at_ready"`
	Delay          int
}

var (
	config Config
	c      string
)

func init() {
	kingpin.CommandLine.Help = "watch files changing and execute specified command."
	kingpin.Flag("dir", "directory to watch.").Short('d').StringVar(&config.Dir)
	kingpin.Flag("exclude", "exclude diretory or files.").Short('e').StringsVar(&config.Excludes)
	kingpin.Flag("recursive", "recursively watch subdirectories.").Short('r').BoolVar(&config.Recursive)
	kingpin.Flag("init", "execute the command immedialtely while the gwatch is ready.").Short('i').BoolVar(&config.ExecuteAtReady)
	kingpin.Flag("delay", "duration delay to execute the command.").IntVar(&config.Delay)
	kingpin.Flag("config", "config path.").Short('c').StringVar(&c)
	kingpin.Arg("command", "command to execute.").StringVar(&config.Command)
	kingpin.Parse()

	// if config file is not speficied, look the current work directory to find gwatch.toml.
	if len(c) == 0 {
		wdConfig := "gwatch.toml"
		_, err := os.Stat(wdConfig)
		if err == nil {
			c = wdConfig
		}
	}
	// if config file exists, load config data from it.
	if len(c) > 0 {
		_, err := toml.DecodeFile(c, &config)
		if err != nil {
			exit(err, 10)
		}
	}
}

func main() {
	// check command.
	if len(config.Command) == 0 {
		kingpin.Usage()
		return
	}

	// process watching directory.
	dir := config.Dir
	if len(dir) == 0 {
		dir = "./"
	}
	fi, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			exit(fmt.Errorf("directory \"%s\" does not exist", dir), 10)
		}
		exit(err, 10)
	}
	if !fi.IsDir() {
		exit(errors.New("invalid directory"), 11)
	}
	dir, err = filepath.Abs(dir)
	if err != nil {
		exit(err, 12)
	}
	if config.Recursive {
		dir += "/..."
	}

	// watch
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	ch := make(chan notify.EventInfo, 1000)

	if err := notify.Watch(dir, ch, notify.All); err != nil {
		log.Fatal(err)
	}
	defer notify.Stop(ch)

	log.Printf("watch directory \"%s\".\n", dir)
	log.Printf("find command \"%s\".\n", config.Command)

	excludes := resolvePaths(config.Excludes)
	log.Printf("exclude directories \"%s\".", strings.Join(excludes, ", "))

	if config.ExecuteAtReady {
		log.Println("gwatch is ready, execute the command.")
		err := executeCommand(config.Command)
		if err != nil {
			exit(err, 13)
		}
	}

	var (
		sync  = true
		delay = config.Delay
	)

	log.Printf("command executing delay %d seconds.", delay)
	if delay <= 0 {
		delay = 2
		log.Printf("command executing delay is fixed to %d seconds.", delay)
	}

out:
	for {
		select {
		case ei := <-ch: // received file change notification
			var (
				path      = ei.Path()
				inExclude bool
			)
			for _, exclude := range excludes {
				if strings.HasPrefix(path, exclude) {
					inExclude = true
					break
				}
			}
			if !inExclude {
				log.Printf("\"%s\" changed\n", path)
				sync = false
			}
		case <-sigCh: // received system interrupt signal
			break out
		case <-time.After(time.Duration(delay) * time.Second): // check if the command should be executed every interval
			if !sync {
				err := executeCommand(config.Command)
				if err != nil {
					exit(err, 13)
				}
				sync = true
			}
		}
	}
}

// resolve symlink or relative path to real absolute path.
func resolvePaths(paths []string) []string {
	for i, path := range paths {
		real, err := filepath.EvalSymlinks(path)
		if err == nil {
			paths[i] = real
		}
		real, err = filepath.Abs(path)
		if err == nil {
			paths[i] = real
		}
	}
	return paths
}

// execute the command.
func executeCommand(cmd string) error {
	command := exec.Command("sh", "-c", cmd)
	var buff bytes.Buffer
	command.Stdout = &buff
	command.Stderr = &buff
	err := command.Run()
	if err == nil {
		fmt.Println(strings.Repeat("=", 80))
		fmt.Print(buff.String())
		fmt.Println(strings.Repeat("=", 80))
	}
	return err
}

// exit progress with the given code.
func exit(err error, code int) {
	log.Println(err)
	os.Exit(code)
}
