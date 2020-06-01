// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/go-ezproxy
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
	"io/ioutil"
	"log"
	"os"
	"time"
)

// Logger is a package logger that can be enabled from client code to allow
// logging output from this package when desired/needed for troubleshooting.
// This variable is exported in order to allow subpackages to use it without
// defining their own. The intent is to make it easier for consumers of the
// package to have one set of methods for enabling or disabling logging output
// for this package and subpackages.
var Logger *log.Logger

func init() {

	// Disable logging output by default unless client code explicitly
	// requests it
	Logger = log.New(os.Stderr, "[ezproxy] ", 0)
	Logger.SetOutput(ioutil.Discard)

}

// EnableLogging enables logging output from this package. Output is muted by
// default unless explicitly requested (by calling this function).
func EnableLogging() {
	Logger.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	Logger.SetOutput(os.Stderr)
}

// DisableLogging reapplies default package-level logging settings of muting
// all logging output.
func DisableLogging() {
	Logger.SetFlags(0)
	Logger.SetOutput(ioutil.Discard)
}

const (

	// This is intended to approximate the `::Limit=X` (where X is a positive
	// whole number) value set within the user.txt EZproxy config file. This
	// package uses this value as a preallocation capacity value for maps and
	// slices.
	SessionsLimit int = 4

	// This is simply a guess to use as a baseline for preallocating slice/map
	// capacity in regards to ALL user sessions
	AllUsersSessionsLimit int = SessionsLimit * 10
)

// This is as of the 6.x series
const (
	SessionIDLength int    = 15
	SessionIDRegex  string = "[a-zA-Z0-9]{15}"
)

const (
	KillSubCmdExitCodeSessionTerminated         int    = 0
	KillSubCmdExitTextTemplateSessionTerminated string = "Session %s terminated"

	KillSubCmdExitCodeSessionNotSpecified int    = 1
	KillSubCmdExitTextSessionNotSpecified string = "Session must be specified"

	KillSubCmdExitCodeSessionDoesNotExist         int    = 3
	KillSubCmdExitTextTemplateSessionDoesNotExist string = "Session %s does not exist"
)

const (

	// TODO: Add default path to ezproxy binary as CLI/config flag
	BinaryName                 string = "ezproxy"
	SubCmdNameSessionTerminate string = "kill"
)

const (

	// DefaultSearchDelay is the delay applied before attempting to read
	// either of the Audit File or Active Users File. This intentional delay
	// is applied in an effort to account for time between EZproxy noting an
	// event and recording it to the file we are reading.
	DefaultSearchDelay time.Duration = 1 * time.Second

	// DefaultSearchRetries is the number of retries beyond the first attempt
	// that will be made after the first attempt at finding active sessions
	// for a specified username yields no results.
	DefaultSearchRetries int = 7
)

// FileEntry reflects a line of text found in a file and the line number
// associated with it
type FileEntry struct {
	Text   string
	Number int
}

// A UserSession represents a session for a specific user account. These
// values are returned after processing either an audit file or the active
// file.
type UserSession struct {
	// SessionID SessionID
	SessionID string
	IPAddress string
	Username  string
}

// UserSessions is a collection of UserSession values. Intended for
// aggregation before bulk processing of some kind.
type UserSessions []UserSession

// SessionsReader is the API for retrieving user sessions from one of the
// audit log or active users and hosts files.
type SessionsReader interface {

	// AllUserSessions returns a list of all session IDs along with their associated
	// IP Address in the form of a slice of UserSession values. This list of
	// session IDs is intended for further processing such as filtering to a
	// specific username or aggregating to check thresholds.
	AllUserSessions() (UserSessions, error)

	// UserSessions uses the previously provided username to return a list of
	// all matching session IDs along with their associated IP Address in the
	// form of a slice of UserSession values.
	UserSessions() (UserSessions, error)

	// SetSearchRetries is a helper method for setting the number of additional
	// retries allowed when receiving zero search results.
	SetSearchRetries(retries int) error

	// SetSearchDelay is a helper method for setting the delay in seconds between
	// search attempts.
	SetSearchDelay(delay int) error
}
