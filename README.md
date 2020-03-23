# clix

[![godoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge)](https://godoc.org/github.com/krostar/clix)
[![Licence](https://img.shields.io/github/license/krostar/clix.svg?style=for-the-badge)](https://tldrlegal.com/license/mit-license)
![Latest version](https://img.shields.io/github/tag/krostar/clix.svg?style=for-the-badge)

[![Build Status](https://img.shields.io/travis/krostar/clix/master.svg?style=for-the-badge)](https://travis-ci.org/krostar/clix)
[![Code quality](https://img.shields.io/codacy/grade/abf18371c077479fa9f8902f64ce0fba/master.svg?style=for-the-badge)](https://app.codacy.com/project/krostar/clix/dashboard)
[![Code coverage](https://img.shields.io/codacy/coverage/abf18371c077479fa9f8902f64ce0fba.svg?style=for-the-badge)](https://app.codacy.com/project/krostar/clix/dashboard)

Modulable CLI builder that call agnostic go idiomatic isolated cli handler, inspired by net/http.

## Motivation

As of today, there are very few nice library to use to handle command line interface; one of them is [cobra](https://github.com/spf13/cobra). It is not super obvious how to make cobra-agnostic cli handler that have dependencies (like a logger, some usecases, ...) that can be injected (so that the handler can be properly tested). This project aims to define cli command tree in a visual way (through the builder pattern) with an easy way to decorelate handler from cli initialization, and to help injecting dependencies (that requires flags to be build for example) in subcommands (like a
logger, a database, ...).

## Usage and example

// TODO
More doc and examples in the clix's [godoc](https://godoc.org/github.com/krostar/clix).

## License

This project is under the MIT licence, please see the LICENCE file.
