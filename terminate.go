// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/brick
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ezproxy

import (
	"bytes"
	"os/exec"
	"strings"
)

type Terminator interface {
	Terminate() []TerminateUserSessionResult
}

// TermTerminateUserSessionResult reflects the result of calling the `kill`
// subcommand of the ezproxy binary to terminate a specific user session.
type TerminateUserSessionResult struct {
	// SessionID is the specific ID associated with the session that we
	// attempt to terminate
	SessionID string

	// ExitCode is what the command called by this application returns
	ExitCode int

	// StdOut is the output (if any) sent to stdout by the command called from
	// this application
	StdOut string

	// StdErr is the output (if any) sent to stderr by the command called from
	// this application
	StdErr string

	// Error is the error (if any) from the attempt to run the specified
	// command
	Error error
}

// TerminateUserSession receives the path to an executable and one or many
// UserSession values, calling the `kill` subcommand of that (presumably
// ezproxy) binary. The result code, stdout, stderr output is captured for
// each subcommand call and returned (along with other details) as a slice of
// `TerminateUserSessionResult`
func TerminateUserSession(executable string, sessions ...UserSession) []TerminateUserSessionResult {

	results := make([]TerminateUserSessionResult, 0, SessionsLimit)

	for _, session := range sessions {

		result := TerminateUserSessionResult{
			SessionID: session.SessionID,
		}

		Logger.Printf(
			"Terminating session %q for username %q ... ",
			session.SessionID,
			session.Username,
		)

		// cmd := exec.Command(
		// 	"echo",
		// 	"hello",
		// )
		cmd := exec.Command(
			executable,
			SubCmdNameSessionTerminate,
			session.SessionID,
		)

		printCmdStr := func(cmd *exec.Cmd) string {
			return strings.Join(cmd.Args, " ")
		}

		Logger.Printf("Executing: %s\n", printCmdStr(cmd))

		// setup buffer to capture stdout
		var cmdStdOut bytes.Buffer
		cmd.Stdout = &cmdStdOut

		//setup buffer to capture stderr
		var cmdStdErr bytes.Buffer
		cmd.Stderr = &cmdStdErr

		cmdErr := cmd.Run()
		if cmdErr != nil {

			switch v := cmdErr.(type) {

			// returned by LookPath when it fails to classify a file as an
			// executable.
			case *exec.Error:

				Logger.Printf(
					"An error occurred attempting to run %q: %v\n",
					printCmdStr(cmd),
					v.Error(),
				)

			// command fail; non-zero (unsuccessful) exit code
			case *exec.ExitError:

				if cmd.ProcessState.ExitCode() == -1 {
					Logger.Println("-1 returned from ExitCode() method")

					if cmd.ProcessState.Exited() {
						Logger.Println("cmd has exited per Exited() method")
					} else {
						Logger.Println("cmd has NOT exited per Exited() method")
					}
				}

			default:

				Logger.Printf(
					"An unexpected error occurred attempting to run %q: [Type: %T Text: %q]\n",
					printCmdStr(cmd),
					cmdErr,
					cmdErr.Error(),
				)

			}

		}

		Logger.Printf("Exit Code: %d\n", cmd.ProcessState.ExitCode())
		Logger.Printf("Captured stdout: %s\n", cmdStdOut.String())
		Logger.Printf("Captured stderr: %s\n", cmdStdErr.String())

		result.ExitCode = cmd.ProcessState.ExitCode()
		result.StdOut = strings.TrimSpace(cmdStdOut.String())
		result.StdErr = strings.TrimSpace(cmdStdErr.String())
		result.Error = cmdErr

		results = append(results, result)

	}

	return results

}

// Terminate attempts to process each UserSession using the provided
// executable, returning the result code, stdout, stderr output as captured
// for each subcommand call (along with other details) as a slice of
// `TerminateUserSessionResult`
func (us UserSessions) Terminate(executable string) []TerminateUserSessionResult {
	return TerminateUserSession(executable, us...)
}
