# xwrd

[![Test status](https://github.com/mlange-42/xwrd/actions/workflows/tests.yml/badge.svg)](https://github.com/mlange-42/xwrd/actions/workflows/tests.yml)
[![GitHub](https://img.shields.io/badge/github-repo-blue?logo=github)](https://github.com/mlange-42/xwrd)
[![MIT license](https://img.shields.io/github/license/mlange-42/xwrd)](https://github.com/mlange-42/xwrd/blob/main/LICENSE)

`xwrd` is a word matching and anagram command line tool.

* Find words by patterns of known letters
* Find anagrams
* Find partial anagrams
* Find (partial) anagrams with unknown letters
* Find multi-word anagrams

## Installation

**Using Go:**

```shell
go install github.com/mlange-42/xwrd@latest
```

**Without Go:**

Download binaries for your OS from the [Releases](https://github.com/mlange-42/xwrd/releases/).


## Getting started

### Install and set a dictionary

Before working with `xwrd`, you need a word list, aka dictionary.
Dictionaries can be installed manually, or downloaded via the `xwrd` CLI.

English dictionary [elasticdog/yawl](https://github.com/elasticdog/yawl) (260k words):

```shell
xwrd dict install en/yawl
xwrd dict set en/yawl
```

German dictionary [enz/german-wordlist](https://github.com/enz/german-wordlist) (680k words):

```shell
xwrd dict install de/enz
xwrd dict set de/enz
```

### Anagrams

Run with words to process:

```shell
xwrd anagram <word1> <word2> ...
```

Run interactively by calling without positional arguments:

```shell
xwrd anagram
```

Partial anagrams:

```shell
xwrd anagram --partial
```

Multi-word anagrams:

```shell
xwrd anagram --multi
```

### Fund words by pattern

Run with words to process:

```shell
xwrd match .a..x
xwrd match .a..x q*a
```

Run interactively by calling without positional arguments:

```shell
xwrd match
```

#### Patterns

`.` (period) strands for one arbitrary letter  
`*` (asterisk) stands for 0 or more arbitrary letters

#### Examples

`a....` - find all 5-letter words starting with 'a'  
`*pf` - find all words ending with 'pf'  
`a....b` - find all words of length 6 stat start with 'a' and end with 'b'  
