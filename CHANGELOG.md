# Changelog

## Overview

All notable changes to this project will be documented in this file.

The format is based on [Keep a
Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Please [open an issue](https://github.com/atc0005/go-ezproxy/issues) for any
deviations that you spot; I'm still learning!.

## Types of changes

The following types of changes will be recorded in this file:

- `Added` for new features.
- `Changed` for changes in existing functionality.
- `Deprecated` for soon-to-be removed features.
- `Removed` for now removed features.
- `Fixed` for any bug fixes.
- `Security` in case of vulnerabilities.

## [Unreleased]

- placeholder

## [v0.1.0] - 2020-06-xx

Initial release!

This release provides an early release version of a library intended for use
with the processing of EZproxy related files and sessions.

### Added

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

- Go modules support (vs classic `GOPATH` setup)

[Unreleased]: https://github.com/atc0005/go-ezproxy/compare/v0.1.0...HEAD
[v0.1.0]: https://github.com/atc0005/go-ezproxy/releases/tag/v0.1.0
