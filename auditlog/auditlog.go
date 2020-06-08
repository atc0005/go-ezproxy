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

package auditlog

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/atc0005/go-ezproxy"
	"github.com/atc0005/go-ezproxy/internal/textutils"
)

const (
	EventLoginSuccess        string = "Login.Success"
	EventLoginSuccessRelogin string = "Login.Success.Relogin"
	EventSessionIPChange     string = "Session.IPChange"
	EventLogout              string = "Logout"

	// EventMinFieldLength is the minimum number of fields required
	// to represent an audit log entry that we will process. The Logout event
	// is 5 fields, Login.Success and Login.Success.Relogin are 6 fields each.
	EventMinFieldLength int = 5
)

// SessionEntry reflects an entry in a audit/YYYYMMDD.txt file that
// contains fields useful (or required) for working with user sessions. Each
// entry in the audit log is tab separated. Not all event types recorded in
// the audit log will have all fields. For example, the `Logout` event type
// does not have an IP Address and the `System` event does not have IP
// Address, Username or Session ID values.
type SessionEntry struct {

	// Datestamp is recorded as a string in an effort to reduce potential
	// friction when ingesting audit log entries. We can convert later
	// when/if needed.
	Datestamp string

	// Event is the event type associated with an entry in the audit file
	Event string

	// IPAddress is an IP Adddress associated with an entry in the audit file
	IPAddress string

	// Username is the username associated with an entry in the audit file
	Username string

	// SessionID is the session ID associated with an entry in the audit file
	SessionID string
}

// SessionEntries is a collection of SessionEntry values.
// Intended for aggregation before bulk processing of some kind.
type SessionEntries []SessionEntry

// auditLogReader represents a file reader specific to EZProxy audit logs
type auditLogReader struct {

	// SearchDelay is the intentional delay between each attempt to open and
	// search the specified filename for the specified username.
	SearchDelay time.Duration

	// SearchRetries is the number of additional search attempts that will be
	// made whenever the initial search attempt returns zero results. Each
	// attempt to read the active file is subject to a race condition; EZproxy
	// does not immediately write session information to disk when creating or
	// terminating sessions, so some amount of delay and a number of retry
	// attempts are used in an effort to work around that write delay.
	SearchRetries int

	// Username is the name of the user account to search for within the
	// specified file.
	Username string

	// Filename is the name of the file which will be parsed/searched for the
	// specified username.
	Filename string
}

// AuditReader is the API for retrieving values from an audit log file
type AuditReader interface {
	ezproxy.SessionsReader

	// SessionEntries uses the previously provided username as a search key, the
	// previously provided filename to search through and returns a slice of
	// SessionEntry values which reflect entries in the specified audit
	// file for that username
	SessionEntries() (SessionEntries, error)

	// AllSessionEntries uses the previously provided filename to search
	// through and return a slice of SessionEntry values which reflect ALL
	// session-related events. The SessionEntry values returned are NOT
	// filtered to a specific username.
	AllSessionEntries() (SessionEntries, error)
}

// Example: "2020-05-24 00:17:37"
const TimeStampLayout string = "2006-01-02 15:04:05"

// AllSessionEntries uses the previously provided filename to search
// through and return a slice of SessionEntry values which reflect ALL
// session-related events. The SessionEntry values returned are NOT
// filtered to a specific username.
func (alr auditLogReader) AllSessionEntries() (SessionEntries, error) {

	// These are events that contain relevant details for our work
	validEvents := []string{
		EventLoginSuccess,
		EventLoginSuccessRelogin,
		EventSessionIPChange,

		// Used to remove any earlier entries since they're no longer relevant
		EventLogout,
	}

	ezproxy.Logger.Printf("Attempting to open %q\n", alr.Filename)

	f, err := os.Open(alr.Filename)
	if err != nil {
		return nil, fmt.Errorf("func AllSessionEntries: error encountered opening file %q: %w", alr.Filename, err)
	}
	defer f.Close()

	ezproxy.Logger.Printf("Searching for: %q\n", alr.Username)

	s := bufio.NewScanner(f)
	var lineno int

	var logoutEvents []SessionEntry

	userSessionIDsIndex := make(map[string]SessionEntry, ezproxy.SessionsLimit)

	// TODO: Does Scan() perform any whitespace manipulation already?
	for s.Scan() {
		lineno++
		currentLine := s.Text()
		ezproxy.Logger.Printf("Scanned line %d from %q: %q\n",
			lineno,
			alr.Filename,
			currentLine,
		)

		currentLine = strings.TrimSpace(currentLine)
		ezproxy.Logger.Printf(
			"Line %d from %q after whitespace removal: %q\n",
			lineno,
			alr.Filename,
			currentLine,
		)

		auditFileEntry := strings.Split(currentLine, "\t")
		if len(auditFileEntry) < EventMinFieldLength {
			continue
		}

		// Event field is the second field, so 1 in a zero-based array/slice
		if !textutils.InList(auditFileEntry[1], validEvents) {
			continue
		}

		// at this point we are dealing with one of these events:
		//
		// Logout
		// Login.Success
		// Login.Success.Relogin

		if strings.EqualFold(auditFileEntry[1], EventLogout) {

			logoutEvents = append(logoutEvents, SessionEntry{
				Datestamp: auditFileEntry[0],
				Event:     auditFileEntry[1],
				IPAddress: "",
				Username:  auditFileEntry[2],
				SessionID: auditFileEntry[3],
			})

			continue
		}

		// at this point we're only dealing with these events:
		//
		// Login.Success
		// Login.Success.Relogin
		userSessionIDsIndex[auditFileEntry[3]] = SessionEntry{
			Datestamp: auditFileEntry[0],
			Event:     auditFileEntry[1],
			IPAddress: auditFileEntry[2],
			Username:  auditFileEntry[3],
			SessionID: auditFileEntry[4],
		}
	}

	ezproxy.Logger.Println("Exited s.Scan() loop")

	// report any errors encountered while scanning the input file
	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("func AllSessionEntries: errors encountered while scanning the input file: %w", err)
	}

	// Loop over logoutEvents, remove matching entries from the
	// userSessionIDsIndex map
	for _, loggedOutSession := range logoutEvents {
		delete(userSessionIDsIndex, loggedOutSession.SessionID)
	}

	// Convert our userSessionIDsIndex map
	userSessions := make(SessionEntries, 0, ezproxy.SessionsLimit)
	for _, entry := range userSessionIDsIndex {
		userSessions = append(userSessions, entry)
	}

	return userSessions, nil

}

// SessionEntries uses the previously provided username as a search key and
// returns a slice of SessionEntry values which reflect entries in the
// specified audit file for that username
func (alr auditLogReader) SessionEntries() (SessionEntries, error) {

	ezproxy.Logger.Printf("Searching for: %q\n", alr.Username)

	searchAttemptsAllowed := alr.SearchRetries + 1

	requestedSessionEntries := make(SessionEntries, 0, ezproxy.SessionsLimit)

	// Perform the search up to X times
	for searchAttempts := 1; searchAttempts <= searchAttemptsAllowed; searchAttempts++ {

		ezproxy.Logger.Printf(
			"Beginning search attempt %d of %d for %q\n",
			searchAttempts,
			searchAttemptsAllowed,
			alr.Username,
		)

		// Intentional delay in an effort to better avoid stale data due to
		// potential race condition with EZproxy write delays.
		ezproxy.Logger.Printf(
			"Intentionally delaying for %v to help avoid race condition due to delayed EZproxy writes\n",
			alr.SearchDelay,
		)
		time.Sleep(alr.SearchDelay)

		allSessionEntries, err := alr.AllSessionEntries()
		if err != nil {
			return nil, fmt.Errorf(
				"func SessionEntries: failed to retrieve all session entries in order to filter to specific username: %w",
				err,
			)
		}

		// Filter ALL session entries in the audit log to the requested username
		for _, entry := range allSessionEntries {
			if strings.EqualFold(entry.Username, alr.Username) {
				requestedSessionEntries = append(requestedSessionEntries, entry)
			}
		}

		// skip further attempts to find entries if we already found some
		if len(requestedSessionEntries) > 0 {
			break
		}

		continue

	}

	return requestedSessionEntries, nil

}

// AllUserSessions returns a list of all session IDs along with their associated
// IP Address in the form of a slice of UserSession values. This list of
// session IDs is intended for further processing such as filtering to a
// specific username or aggregating to check thresholds.
func (alr auditLogReader) AllUserSessions() (ezproxy.UserSessions, error) {

	allUserSessions := make(ezproxy.UserSessions, 0, ezproxy.AllUsersSessionsLimit)

	allSessionEntries, err := alr.AllSessionEntries()
	if err != nil {
		return nil, fmt.Errorf(
			"func AllUserSessions: failed to retrieve all session entries in order to convert to user sessions for all users: %w",
			err,
		)
	}

	for _, entry := range allSessionEntries {
		if strings.EqualFold(entry.Username, alr.Username) {
			allUserSessions = append(allUserSessions, ezproxy.UserSession{
				Username:  entry.Username,
				IPAddress: entry.IPAddress,
				SessionID: entry.SessionID,
			})
		}
	}

	return allUserSessions, nil

}

// UserSessions uses the previously provided username to return a list of all
// matching session IDs along with their associated IP Address in the form of
// a slice of UserSession values.
func (alr auditLogReader) UserSessions() (ezproxy.UserSessions, error) {

	ezproxy.Logger.Printf("Searching for: %q\n", alr.Username)

	searchAttemptsAllowed := alr.SearchRetries + 1

	sessionEntries := make(SessionEntries, 0, ezproxy.SessionsLimit)

	// Perform the search up to X times
	for searchAttempts := 1; searchAttempts <= searchAttemptsAllowed; searchAttempts++ {

		ezproxy.Logger.Printf(
			"Beginning search attempt %d of %d for %q\n",
			searchAttempts,
			searchAttemptsAllowed,
			alr.Username,
		)

		// Intentional delay in an effort to better avoid stale data due to
		// potential race condition with EZproxy write delays.
		ezproxy.Logger.Printf(
			"Intentionally delaying for %v to help avoid race condition due to delayed EZproxy writes\n",
			alr.SearchDelay,
		)
		time.Sleep(alr.SearchDelay)

		var err error
		sessionEntries, err = alr.SessionEntries()
		if err != nil {
			return nil, fmt.Errorf("func UserSessions: unable to convert audit log session entries to user sessions: %w", err)
		}

		// skip further attempts to find entries if we already found some
		if len(sessionEntries) > 0 {
			break
		}

		continue

	}

	return sessionEntries.UserSessions(), nil

}

// UserSession converts an SessionEntry value to UserSession value
func (se SessionEntry) UserSession() ezproxy.UserSession {
	return ezproxy.UserSession{
		SessionID: se.SessionID,
		IPAddress: se.IPAddress,
		Username:  se.Username,
	}
}

// UserSessions converts a collection of SessionEntry values into a collection
// of UserSession values.
func (se SessionEntries) UserSessions() ezproxy.UserSessions {

	userSessions := make(ezproxy.UserSessions, 0, ezproxy.SessionsLimit)

	for idx := range se {
		userSessions = append(userSessions, ezproxy.UserSession{
			SessionID: se[idx].SessionID,
			IPAddress: se[idx].IPAddress,
			Username:  se[idx].Username,
		})
	}

	return userSessions

}

// NewReader creates a new instance of an AuditReader that provides access to
// collections of user sessions and audit log session entries specific to the
// specified username.
func NewReader(username string, filename string) (AuditReader, error) {

	if username == "" {
		return nil, errors.New(
			"func NewReader: missing username",
		)
	}

	if filename == "" {
		return nil, errors.New(
			"func NewReader: missing filename",
		)
	}

	reader := auditLogReader{
		SearchDelay:   ezproxy.DefaultSearchDelay,
		SearchRetries: ezproxy.DefaultSearchRetries,
		Username:      username,
		Filename:      filename,
	}

	return &reader, nil

}

// SetSearchRetries is a helper method for setting the number of additional
// retries allowed when receiving zero search results.
func (alr *auditLogReader) SetSearchRetries(retries int) error {
	if retries < 0 {
		return fmt.Errorf("func SetSearchRetries: %d is not a valid number of search retries", retries)
	}

	alr.SearchRetries = retries

	return nil
}

// SetSearchDelay is a helper method for setting the delay in seconds between
// search attempts.
func (alr *auditLogReader) SetSearchDelay(delay int) error {
	if delay < 0 {
		return fmt.Errorf("func SetSearchDelay: %d is not a valid number of seconds for search delay", delay)
	}

	alr.SearchDelay = time.Duration(delay) * time.Second

	return nil
}
