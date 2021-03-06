# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0] - 2020-11-01

### Changed
- Refactor `pkg/caller`.  Make configurable timeout
- Refactor factory, to adding new configurable options for caller's timeout
- Retest & recheck supervisor's flows & processes
- Retest & recheck all available `pkg` domains
- Update example
- Update readme

## [1.0.0-alpha.3] - 2020-10-09
### Removed
- Remove watch as unused api function
- Remove primitive on process registry, move it to process builder

### Added
- Adding supervisor 
- Adding auto retry mechanism when calling to plugin 
- Adding new example for supervisor

### Changed
- Change proto payload from string to bytes
- Update examples `hello-grpc` & `hello-rest`

## [1.0.0-alpha.2] - 2020-09-21
### Changed 
- Refactor main factory build to simplify initiation process
- Update `examples`

## [1.0.0-alpha.1] - 2020-09-19
### Added
- Basic example using implementation using grpc, [hello-grpc](https://github.com/quadroops/goplugin/tree/master/examples/basic/hello-grpc)
- Refactor `proto/plugin.proto` adding `go_option` package

## [1.0.0-alpha.0] - 2020-09-17
### Added
- Basic example using implementation using rest, [hello-rest](https://github.com/quadroops/goplugin/tree/master/examples/basic/hello-rest)

[Unreleased]: https://github.com/quadroops/goplugin/compare/v1.0.0-alpha.2...HEAD
[1.0.0]: https://github.com/quadroops/goplugin/compare/v1.0.0-alpha.3...v1.0.0
[1.0.0-alpha.3]: https://github.com/quadroops/goplugin/compare/v1.0.0-alpha.2...v1.0.0-alpha.3
[1.0.0-alpha.2]: https://github.com/quadroops/goplugin/compare/v1.0.0-alpha.1...v1.0.0-alpha.2
[1.0.0-alpha.1]: https://github.com/quadroops/goplugin/compare/v1.0.0-alpha.0...v1.0.0-alpha.1
[1.0.0-alpha.0]: https://github.com/quadroops/goplugin/releases/tag/v1.0.0-alpha.0
