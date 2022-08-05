# STOC (Search Text On Condition)

## Overview

Bridges the gap between simple searching on a single string, and more powerful search using regex.
Plain english and simple boolean algebra is used.

## Dependencies

Only dependencies are for unit testing using ginkgo and gomega

## Current features

- and, or, not conditions
- quotes around strings
- custom types (see tokens_definition.go ```prepareDefaultTokensDefinition()``` function)

## Examples

- not ("foo" or 'baa') and baz // brackets work with unary and binary conditionals
- foo or "baa or baz" // keywords can be in expressions using quotes
- foo and "foo" // this is equivalent to the next example
- foo // you can do simple search as well

## RoadMap

- escaping quotes in expressions
- using golang error more instead of just printing an error, since this is meant to be a library not an app
- some way to make strings case-sensitive and case-insensitive, both for the full command and part of it.
  Probably can do using single and double quotes. quotes would be optional and default to case-insensitive
  but would add ability to override globally in that case global override would also override double and single quotes
  to denote case insensitivity in the command
- the reference app (command line)
- context support

## Potential additions (still considering)

- overriding keywords e.g. \& or \and to override keywords when not using quotes (might be overkill)
- might be good to include basic regex concepts like \w\d\s word, digit, whitespace \W\D\S not word, digit,
  whitespace, and potentially ^ and $

## Frontends

https://github.com/kranzuft/stoc-cli (reference implementation)
