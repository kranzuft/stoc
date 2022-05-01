# Lofty

## Overview 

Bridges the gap between simple searching on a single string, and more powerful search using regex.
Plain english and simple boolean algebra is used.

## RoadMap

- quotes around strings
- overriding keywords e.g. \& or \and to override keywords when not using quotes
- using golang error more instead of just printing an error, since this is meant to be a library not an app
- some way to make strings case sensitive and case insensitive, both for the full command and part of it. 
  Probably can do using single and double quotes. quotes would be optional and default to case insensitive 
  but would add ability to override globally in that case global override would also override double and single quotes 
  to denote case insensitivity in the command
- the reference app (command line)
- might be good to  include basic regex concepts like \w\d\s	word, digit, whitespace \W\D\S	not word, digit, 
  whitespace, and potentially ^ and $