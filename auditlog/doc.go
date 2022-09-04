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

/*
Package auditlog is intended for the processing of EZproxy audit log files.

# Overview

EZproxy can be configured to record a variety of audit/security related log
entries to each daily audit log file. The specific events recorded vary based
on the chosen configuration settings. Based on the author's experience, the
more detail you enable the better your sysadmin and network security team will
be able to quickly detect and recover from active abuse cases. This
necessitates not only recording the information at the time the event occurs,
but also retaining the information long enough to be useful. A balance should
be struck between that need and legitimate expectations of privacy. It is up
to each site to determine how best to comply with local policy and those
expectations of privacy.

The relevance to us (and this library) is that due to how much the EZproxy
configuration may vary for each instance, we're not really able to *reliably*
use the audit log files as a source of active user session details. Instead,
we have to parse the Active Users and Hosts "state" file for that information.

Even so, the audit log files contain event-based data that is of potential use
to security-related applications in the future. Depending on the need, this
library may be further updated to expose those details to client applications
(e.g., number of failed login attempts and whether a set of failed login
attempts eventually resulted in a successful login).

# Field Types

Each audit log file is composed of tab-separated fields, so presumably a CSV
reader could (and perhaps in hindsight *should*) be used with this file. As of
version 6.x of EZproxy, there are six fields in each audit log file. Each of
the field names below are taken directly from a real audit log file.

01) Date/Time
02) Event
03) IP
04) Username
05) Session
06) Other

Currently we only deal with the first 5 fields.

# Race Condition

NOTE: EZproxy does not immediately update the Active Users and Hosts "state"
file with state changes; when a user account logs in/out, there is a race
condition between when that information is updated and when a Reader created
from this package attempts to read the current state and reconstruct User
Sessions. In an effort to workaround this race condition, this package
attempts to retry session read attempts a limited number of times by default
before giving up. This retry behavior (including a delay between retry
attempts) can be modified by the caller as needed.

This race condition is also believed to affect the audit log, so the same
retry/retry delay behavior is provided for the audit log Reader as with the
active file reader.
*/
package auditlog
