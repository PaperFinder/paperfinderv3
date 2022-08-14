# Paperfinder v3

The third iteration of this godforsaken piece of junk

## Features

- Less and lighter dependencies
- New design
- Fast
- Server side rendering

## Dependencies

1. Poppler (provides the `pdftotext` executable)
2. ImageMagick (provides the `convert` executable)
3. Go 1.18+

## Usage

- To retrieve all past papers run with `--papers` flag
- To reindex all past papers run with `--index` flag
- To start the web server run with no flags

`MANGLE_KEY` environment variable is a secret key to prevent scraping from server

`SEARCH_THRESHOLD` environment variable specifies the % of words that have to match for a search to come as valid

## Environment

Example .env

```bash
PAPER_CONFIG=papers.json
PAPER_FOLDER=_pastpapers
MANGLE_KEY=secret
SEARCH_THRESHOLD=80

HTTP_PORT=8080
```

## Upcoming

- [ ] HTTPS support
- [ ] Better design
- [ ] Support for more subjects, exams and exam boards
