/*
Copyright 2020 sh0rez

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// adapted from: https://github.com/sh0rez/kubectl-neat-diff
package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"

	neat "github.com/addreas/kubectl-neat/pkg"
)

var ignoreLinesRegexes *[]string

var cmd = cobra.Command{
	Use:   "kubectl-neat-diff <file1> <file2>",
	Short: "Remove fields from kubectl diff that carry low / no information",
	Long: `

	De-clutter your kubectl diff output using kubectl-neat (looking at you, managedFields)

	To use, set it as your KUBECTL_EXTERNAL_DIFF tool:

	# append to ~/.bashrc or similar:
	export KUBECTL_EXTERNAL_DIFF=kubectl-neat-diff


`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := neatifyDir(args[0]); err != nil {
			return err
		}
		if err := neatifyDir(args[1]); err != nil {
			return err
		}

		c := exec.Command("diff", "-uN", args[0], args[1])
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Run()
		os.Exit(c.ProcessState.ExitCode())

		return nil
	},
}

func init() {
	ignoreLinesRegexes = cmd.Flags().StringSliceP("ignore-matching-lines", "I", []string{},
		"Ignore changes whose lines all match RegExp.")
}

func Execute() {
	log.SetFlags(0)

	if err := cmd.Execute(); err != nil {
		log.Fatalln("Error:", err)
	}
}

func neatifyDir(dir string) error {
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, fi := range fis {
		filename := filepath.Join(dir, fi.Name())
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}

		n, err := neat.NeatYAMLOrJSON(data, "same")
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(filename, []byte(n), fi.Mode()); err != nil {
			return err
		}
	}

	return nil
}
