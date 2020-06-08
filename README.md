<!-- omit in toc -->
# go-ezproxy

Go library and tooling for working with EZproxy.

[![Latest Release](https://img.shields.io/github/release/atc0005/go-ezproxy.svg?style=flat-square)][release-latest]
[![GoDoc](https://godoc.org/github.com/atc0005/go-ezproxy?status.svg)][docs-homepage]
![Validate Codebase](https://github.com/atc0005/go-ezproxy/workflows/Validate%20Codebase/badge.svg)
![Validate Docs](https://github.com/atc0005/go-ezproxy/workflows/Validate%20Docs/badge.svg)

<!-- omit in toc -->
## Table of contents

- [Status](#status)
- [Overview](#overview)
- [Project home](#project-home)
- [Features](#features)
  - [Current](#current)
  - [Missing](#missing)
- [Changelog](#changelog)
- [Documentation](#documentation)
- [Examples](#examples)
- [License](#license)
- [References](#references)
  - [Related projects](#related-projects)
  - [Official EZproxy docs](#official-ezproxy-docs)

## Status

Alpha; very much getting a feel for how the project will be structured
long-term and what functionality will be offered.

As of this writing, the existing functionality was added specifically to
support another project in-development named "brick". This project is subject
to change in order to better support that one.

## Overview

This library is intended to provide common EZproxy-related functionality such
as reporting or terminating active login sessions (either for all usernames or
specific usernames), filtering (or not) audit file entries or traffic patterns
(not implemented yet) for specific usernames or domains.

**NOTE**: Just to be perfectly clear, this library is intended to supplement
the provided functionality of the official OCLC-developed/supported `EZproxy`
application, not in any way replace it.

## Project home

See [our GitHub repo][repo-url] for the latest code,
to file an issue or submit improvements for review and potential inclusion
into the project.

## Features

### Current

- generate a list of audit records for session-related events
  - for all usernames
  - for a specific username

- generate a list of active sessions using audit log
  - using entires without a corresponding logout event type

- generate a list of active sessions using active file
  - for all usernames
  - for a specific username

- terminate user sessions
  - single user session
  - bulk user sessions

### Missing

- Anything to do with traffic log entries

## Changelog

See the [`CHANGELOG.md`](CHANGELOG.md) file for the changes associated with
each release of this application. Changes that have been merged to `master`,
but not yet an official release may also be noted in the file under the
`Unreleased` section. A helpful link to the Git commit history since the last
official release is also provided for further review.

## Documentation

Please see our [GoDoc][docs-homepage] coverage. If something doesn't make
sense, please [file an issue][repo-url] and note what is (or was) unclear.

## Examples

Please see our [GoDoc][docs-homepage] coverage for general usage and the
[examples](examples/README.md) doc for a list of applications developed using
this module.

## License

Taken directly from the [`LICENSE`](LICENSE) and [`NOTICE.txt`](NOTICE.txt) files:

```License
Copyright 2020-Present Adam Chalkley

https://github.com/atc0005/go-ezproxy/blob/master/LICENSE

Licensed under the Apache License, Version 2.0 (the "License"); you may not use
this file except in compliance with the License. You may obtain a copy of the
License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied. See the License for the
specific language governing permissions and limitations under the License.
```

## References

### Related projects

- [atc0005/brick](https://github.com/atc0005/brick) project
  - this project uses this library to provides tools (two as of this writing)
    intended to help manage login sessions.

### Official EZproxy docs

- <https://help.oclc.org/Library_Management/EZproxy/EZproxy_configuration/EZproxy_system_elements>
- <https://help.oclc.org/Library_Management/EZproxy/Configure_resources/Audit>
- <https://help.oclc.org/Library_Management/EZproxy/Configure_resources/LogFormat>
- <https://help.oclc.org/Library_Management/EZproxy/Configure_resources/Option_LogSession>
- <https://help.oclc.org/Library_Management/EZproxy/Configure_resources/Option_LogUser>
- <https://help.oclc.org/Library_Management/EZproxy/Get_started/Join_the_EZproxy_listserv_and_Community_Center>

<!-- Footnotes here  -->

[repo-url]: <https://github.com/atc0005/go-ezproxy>  "This project's GitHub repo"

[docs-homepage]: <https://godoc.org/github.com/atc0005/go-ezproxy>  "GoDoc coverage"

[release-latest]: <https://github.com/atc0005/go-ezproxy/releases/latest>  "Latest Release"

<!-- []: PLACEHOLDER "DESCRIPTION_HERE" -->
