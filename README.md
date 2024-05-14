# copyrightlsp

This is a language server for handling copyright notices at the top of a file.

## Features

- Show an error on the first line of a document if no copyright header is found
- Allow to insert a copyright header with a code action if no header is found

## Configuration

The server expect a `settings` object with the following structure:

```json
{
    "templates": {
        "sh": ["# Copyright {year}"],
        "html": ["<!-- Copyright {year} -->"],
        "javascript": ["/*", " * Copyright {year}", " */"]
    },
    "searchRanges": {
        "sh": 1
    }
}
```

By default no `templates` and no `searchRanges` are configured!

### Templates

A template is an array of strings, which represent the different lines of a
copyright comment. The placeholder `{year}` can be used to insert the current
year.

### Search ranges

The search for the copyright comment starts always at the top of a document and
includes the same amount of lines as the template. With this option the search
range can be increased to includes `n` additional lines, which allow copyright
comments on other lines then the first. This can be useful for files which
start with a shebang (`#!`) line. If no value is specified for a language then
`0` is used as default value.
