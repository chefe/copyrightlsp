# copyrightlsp

This is a language server for handling copyright notices at the top of a file.
Currently only inserting a notice for a given language is supported.

## Configuration

The server expect a `settings` object with the following structure.

```json
{
    "templates": {
        "sh": ["# Copyright {year}"],
        "html": ["<!-- Copyright {year} -->"],
        "javascript": ["/*", " * Copyright {year}", " */"],
    }
}
```
