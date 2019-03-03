//go:generate go run data/generate.go
package main

import (
	"fmt"
	"strings"
	"sync"

	shellquote "github.com/kballard/go-shellquote"
)

// Version should be populated with -ldflags on build with the semver version
var Version string

// Commit should be populated with -ldflags on build with the current git commit
var Commit string

func main() {
	config := InitConfig()
	godev := InitGoDev(config)
	godev.Start()
}

func InitGoDev(config *Config) *GoDev {
	return &GoDev{
		config: config,
		logger: InitLogger(&LoggerConfig{
			Name:   "main",
			Format: "production",
			Level:  config.LogLevel,
		}),
	}
}

type GoDev struct {
	config *Config
	logger *Logger
}

func (godev *GoDev) Start() {
	defer godev.logger.Infof("godev has ended")
	godev.logger.Infof("godev has started")
	if godev.config.RunView {
		godev.viewFile()
	} else if godev.config.RunVersion {
		fmt.Printf("godev %s-%s\n", Version, Commit)
	} else {
		godev.startWatching()
	}
}

func (godev *GoDev) startWatching() {
	config := godev.config
	logger := godev.logger
	godev.logUniversalConfigurations()
	godev.logWatchModeConfigurations()
	watcher := InitWatcher(&WatcherConfig{
		FileExtensions: config.FileExtensions,
		IgnoredNames:   config.IgnoredNames,
		RefreshRate:    config.Rate,
		LogLevel:       config.LogLevel,
	})
	watcher.RecursivelyWatch(config.WatchDirectory)
	var pipeline []*ExecutionGroup
	for _, execGroup := range config.ExecGroups {
		executionGroup := &ExecutionGroup{}
		var executionCommands []*Command
		commands := strings.Split(execGroup, config.CommandsDelimiter)
		for _, command := range commands {
			if sections, err := shellquote.Split(command); err != nil {
				panic(err)
			} else {
				executionCommands = append(
					executionCommands,
					InitCommand(&CommandConfig{
						Application: sections[0],
						Arguments:   sections[1:],
						LogLevel:    config.LogLevel,
					}),
				)
			}
		}
		executionGroup.commands = executionCommands
		pipeline = append(pipeline, executionGroup)
	}
	runner := InitRunner(&RunnerConfig{
		Pipeline: pipeline,
		LogLevel: config.LogLevel,
	})

	var wg sync.WaitGroup
	watcher.BeginWatch(&wg, func(events *[]WatcherEvent) bool {
		for _, e := range *events {
			logger.Trace(e)
		}
		runner.Trigger()
		return true
	})

	logger.Infof("started watcher at %s", config.WatchDirectory)

	runner.Trigger()
	wg.Wait()
}

func (godev *GoDev) logUniversalConfigurations() {
	config := godev.config
	logger := godev.logger
	logger.Debugf("flag - init       : %v", config.RunInit)
	logger.Debugf("flag - test       : %v", config.RunTest)
	logger.Debugf("flag - view       : %v", config.RunView)
	logger.Debugf("watch directory   : %s", config.WatchDirectory)
	logger.Debugf("build output      : %s", config.BuildOutput)
}

func (godev *GoDev) logWatchModeConfigurations() {
	config := godev.config
	logger := godev.logger
	logger.Debugf("file extensions   : %v", config.FileExtensions)
	logger.Debugf("ignored names     : %v", config.IgnoredNames)
	logger.Debugf("refresh interval  : %v", config.Rate)
	logger.Debugf("execution delim   : %s", config.CommandsDelimiter)
	logger.Debug("execution groups as follows...")
	for egIndex, execGroup := range config.ExecGroups {
		logger.Debugf("  %v) %s", egIndex+1, execGroup)
		commands := strings.Split(execGroup, config.CommandsDelimiter)
		for cIndex, command := range commands {
			sections, err := shellquote.Split(command)
			if err != nil {
				panic(err)
			}
			app := sections[0]
			args := sections[1:]
			logger.Debugf("    %v > %s %v", cIndex+1, app, args)
		}
	}
}

func (godev *GoDev) viewFile() {
	config := godev.config
	logger := godev.logger
	switch strings.ToLower(config.View) {
	case "dockerfile":
		logger.Info("previewing contents of Dockerfile")
		fmt.Println(DataDockerfile)
		logger.Info("end of preview for contents of Dockerfile")
	case "makefile":
		logger.Info("previewing contents of Makefile")
		fmt.Println(DataMakefile)
		logger.Info("end of preview for contents of Makefile")
	case ".dockerignore":
		logger.Info("previewing contents of .dockerignore")
		fmt.Println(DataDotDockerignore)
		logger.Info("end of preview for contents of .dockerignore")
	case ".gitignore":
		logger.Info("previewing contents of .gitignore")
		fmt.Println(DataDotGitignore)
		logger.Info("end of preview for contents of .gitignore")
	case "main.go":
		logger.Info("previewing contents of main.go")
		fmt.Println(DataMainDotgo)
		logger.Info("end of preview for contents of main.go")
	case "go.mod":
		logger.Info("previewing contents of go.mod")
		fmt.Println(DataGoDotMod)
		logger.Info("end of preview for contents of go.mod")
	}
}