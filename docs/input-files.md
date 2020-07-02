<!-- omit in toc -->
# go-ezproxy: Input file formats

<!-- omit in toc -->
## Table of Contents

- [Project home](#project-home)
- [Overview](#overview)
- [File formats](#file-formats)
  - [Audit log](#audit-log)
    - [Overview](#overview-1)
    - [Field types](#field-types)
    - [Race Condition](#race-condition)
  - [Active Users and Hosts file](#active-users-and-hosts-file)
    - [Overview](#overview-2)
    - [Known Types](#known-types)
    - [Unknown Types](#unknown-types)
    - [Line Ordering](#line-ordering)
    - [Field Numbers](#field-numbers)
      - [Logins](#logins)
      - [Sessions](#sessions)
    - [Race condition](#race-condition-1)
- [Other documentation](#other-documentation)

## Project home

See [our GitHub repo](https://github.com/atc0005/go-ezproxy) for the latest
code, to file an issue or submit improvements for review and potential
inclusion into the project.

## Overview

This document attempts to summarize the supported input file formats, their
fields and any "gotchas" discovered along the way. Corrections or
clarifications are welcome!

## File formats

<!--
  NOTE: This content is mirrored within the auditlog/doc.go and
  activefile/doc.go files. Mirroring the content between both locations
  *seems* like a workable solution (based on the expected low volume of doc
  changes), but in the future it is likely we will collapse the two into one
  source and just add one or more forwarding references as needed.
-->

### Audit log

#### Overview

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

#### Field types

Each audit log file is composed of tab-separated fields, so presumably a CSV
reader could (and perhaps in hindsight *should*) be used with this file. As of
version 6.x of EZproxy, there are six fields in each audit log file. Each of
the field names below are taken directly from a real audit log file.

| Field | Field name | Note               |
| ----- | ---------- | ------------------ |
| 1     | Date/Time  |                    |
| 2     | Event      |                    |
| 3     | IP         |                    |
| 4     | Username   |                    |
| 5     | Session    |                    |
| 6     | Other      | actual column name |

Currently we only deal with the first 5 fields.

#### Race Condition

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

### Active Users and Hosts file

#### Overview

There is only ever one Active Users and Hosts "state" file at a time. While
Host entries are (from what can be observed) consolidated on one line, session
entries are composed of multiple lines in a very specific order, each with
space-separated fields. These order-specific lines and fields are joined in
order to reconstruct a User Session that reflects active user sessions within
EZproxy.

#### Known Types

Known entry types include (but are not limited to):

| Type              | Prefix |
| ----------------- | ------ |
| Host              | `H`    |
| Group             | `g`    |
| Session           | `S`    |
| Username or Login | `L`    |

Currently, only the last two types (`S`, `L`) are relevant to our purposes.

#### Unknown Types

These types have been observed, but not researched sufficiently to identify
their purpose (Pull Requests for this are welcome!):

| Type | Prefix | Note             |
| ---- | ------ | ---------------- |
| ?    | `P`    |                  |
| ?    | `M`    |                  |
| ?    | `s`    | lowercase letter |

#### Line Ordering

For our purposes, we match lines that start with a capital letter `S` and pair
it with the first line following it that begins with a capital letter `L`. We
skip over any line that begins with a lowercase letter `s`; we do not use the
value provided by this line.

When we match a line beginning with a capital `S`, these are the only supported
line orderings:

`S`
`s`
`L`

and:

`S`
`L`

#### Field Numbers

##### Logins

The line for for Logins (`L`) is composed of 2 fields:

| Field | Value    | Note |
| ----- | -------- | ---- |
| 1     | `L`      |      |
| 2     | Username |      |

##### Sessions

The line for Sessions (`S`) is composed of 11 fields:

| Field | Value                                                 | Note                                                         |
| ----- | ----------------------------------------------------- | ------------------------------------------------------------ |
| 1     | `S`                                                   | capital letter                                               |
| 2     | Session ID                                            |                                                              |
| 3     | unknown                                               | appears to be a UNIX timestamp                               |
| 4     | unknown                                               | appears to be two UNIX timestamps separated by a literal dot |
| 5     | unknown                                               | integer; number 1 was common                                 |
| 6     | EZproxy `MaxLifetime` or `User Session` timeout value |                                                              |
| 7     | IP Address                                            |                                                              |
| 8     | unknown                                               | 0 is recorded                                                |
| 9     | unknown                                               | 0 is recorded                                                |
| 10    | unknown                                               | 0 is recorded                                                |
| 11    | unknown                                               | asterisk is recorded                                         |

#### Race condition

NOTE: EZproxy does not immediately update the Active Users and Hosts "state"
file with state changes; when a user account logs in/out, there is a race
condition between when that information is updated and when a Reader created
from this package attempts to read the current state and reconstruct User
Sessions. In an effort to workaround this race condition, this package
attempts to retry session read attempts a limited number of times by default
before giving up. This retry behavior (including a delay between retry
attempts) can be modified by the caller as needed.

## Other documentation

See [the doc index](README.md) and the [main README](../README.md) for
additional information (including official EZproxy reference links).
