// Copyright © 2019 Konstantinos Konstantinidis <kkonstan@ianos.co.uk>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/viper"
)

// kubectlCommand
// Run kubectl with specified arguments
func kubectlCommand(cmdArgs []string) []byte {
	var (
		cmdOut []byte
		err    error
	)
	cmdName := viper.GetString("kubectl_binary")
	if cmdOut, err = exec.Command(cmdName, cmdArgs...).CombinedOutput(); err != nil {
		//fmt.Fprintln(os.Stderr, "There was an error running kubectl command: ", err)
	}
	return cmdOut
}

// kubectlClusters
// Detect kubectl configured contexts
func kubectlClusters() map[string]string {

	var clusters map[string]string

	clusters = make(map[string]string)

	cmdName := viper.GetString("kubectl_binary")
	cmdArgs := []string{"config", "get-contexts", "--no-headers", "--output=name"}
	cmd := exec.Command(cmdName, cmdArgs...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			clusters[scanner.Text()] = scanner.Text()
		}
	}()

	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
		os.Exit(1)
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
		os.Exit(1)
	}

	return clusters
}

// kubectlGetClientVersion
// Return kubectl client version
func kubectlGetClientVersion() string {
	var (
		cmdOut []byte
		err    error
	)

	cmdName := viper.GetString("kubectl_binary")
	cmdArgs := []string{"version", "--client", "--short"}

	// Run kubectl version
	if cmdOut, err = exec.Command(cmdName, cmdArgs...).CombinedOutput(); err != nil {
		// Not installed, return empty string
		return ""
	}

	// Determine version from output
	f := func(c rune) bool {
		return c == rune(':')
	}

	// Return trimmed version
	return strings.TrimSpace(strings.FieldsFunc(strings.Replace(string(cmdOut), "\n", ":", -1), f)[1])
}

// kubectlExists
// Check that kubectl exists
func kubectlExists() bool {
	return kubectlGetClientVersion() != ""
}

func init() {
	viper.SetDefault("kubectl_binary", "kubectl")
}
