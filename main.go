// Copyright 2013 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package contains the main entry point of gd.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rakyll/command"
	"github.com/rakyll/drive/commands"
	"github.com/rakyll/drive/config"
)

var context *config.Context

func main() {
	command.On("init", &initCmd{})       // inits a directory as gd directory
	command.On("pull", &pullCmd{})       // pulls from Google Drive
	command.On("push", &pushCmd{})       // pushes to Google Drive
	command.On("diff", &diffCmd{})       // diff a file
	command.On("publish", &publishCmd{}) // publish a file
	command.ParseAndRun()
}

type initCmd struct{}

func (cmd *initCmd) Flags(fs *flag.FlagSet) *flag.FlagSet {
	return fs
}

func (cmd *initCmd) Run(args []string) {
	exitWithError(commands.New(initContext(args), nil).Init())
}

type pullCmd struct {
	isRecursive *bool
	isNoPrompt  *bool
}

func (cmd *pullCmd) Flags(fs *flag.FlagSet) *flag.FlagSet {
	cmd.isRecursive = fs.Bool("r", true, "perform the pull action recursively")
	cmd.isNoPrompt = fs.Bool("no-prompt", false, "no prompt before applying the pull action")
	return fs
}

func (cmd *pullCmd) Run(args []string) {
	context, path := discoverContext(args)
	exitWithError(commands.New(context, &commands.Options{
		Path:        path,
		IsRecursive: *cmd.isRecursive,
		IsNoPrompt:  *cmd.isNoPrompt,
	}).Pull())
}

type pushCmd struct {
	isRecursive *bool
	isNoPrompt  *bool
}

func (cmd *pushCmd) Flags(fs *flag.FlagSet) *flag.FlagSet {
	cmd.isRecursive = fs.Bool("r", true, "perform the push action recursively")
	cmd.isNoPrompt = fs.Bool("no-prompt", false, "no prompt before applying the push action")
	return fs
}

func (cmd *pushCmd) Run(args []string) {
	context, path := discoverContext(args)
	exitWithError(commands.New(context, &commands.Options{
		Path:        path,
		IsRecursive: *cmd.isRecursive,
		IsNoPrompt:  *cmd.isNoPrompt,
	}).Push())
}

type diffCmd struct{}

func (cmd *diffCmd) Flags(fs *flag.FlagSet) *flag.FlagSet {
	return fs
}

func (cmd *diffCmd) Run(args []string) {
	context, path := discoverContext(args)
	exitWithError(commands.New(context, &commands.Options{
		Path: path,
	}).Diff())
}

type publishCmd struct{}

func (cmd *publishCmd) Flags(fs *flag.FlagSet) *flag.FlagSet {
	return fs
}

func (cmd *publishCmd) Run(args []string) {
	context, path := discoverContext(args)
	exitWithError(commands.New(context, &commands.Options{
		Path: path,
	}).Publish())
}

func initContext(args []string) *config.Context {
	var err error
	context, err = config.Initialize(getContextPath(args))
	exitWithError(err)
	return context
}

func discoverContext(args []string) (*config.Context, string) {
	var err error
	context, err = config.Discover(getContextPath(args))
	exitWithError(err)
	relPath := ""
	if len(args) > 0 {
		relPath, err = filepath.Rel(context.AbsPath, args[0])
	}
	exitWithError(err)
	return context, relPath
}

func getContextPath(args []string) (contextPath string) {
	if len(args) > 0 {
		contextPath = args[0]
	}
	if contextPath == "" {
		contextPath, _ = os.Getwd()
	}
	return
}

func exitWithError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}