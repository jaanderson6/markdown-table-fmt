# markdown-table-fmt

Markdown tables drift out of alignment fast: a cell gets edited, the pipes
no longer line up, someone pastes in a row from a spreadsheet with no
padding at all. The table still renders fine once it hits a markdown
processor, but it's unreadable as plain text, and diffs on it turn into
noise because every column shifts by a few spaces.

This reformats markdown tables in place: aligns the pipes, pads every cell
to its column's width, and normalizes the separator row's dashes and
alignment colons. It leaves everything outside of a table (headings,
paragraphs, code blocks) untouched.

Before:

```
| Name |Age| City|
|---|---|---|
| Alice | 30 |New York|
|Bob|4|LA|
```

After:

```
| Name  | Age | City     |
| ----- | --- | -------- |
| Alice | 30  | New York |
| Bob   | 4   | LA       |
```

## Usage

As a CLI, reading from stdin or a file argument, writing the result to
stdout:

```sh
go run ./cmd/mdfmt < messy.md > clean.md
go run ./cmd/mdfmt messy.md > clean.md
go run ./cmd/mdfmt -w messy.md   # rewrite messy.md in place
```

As a library:

```go
import mdtable "github.com/jaanderson6/markdown-table-fmt"

clean := mdtable.Format(messy)
```

`Format` takes a whole document, finds each markdown table in it, and
rewrites just those lines. It handles tables missing their outer pipes,
rows with fewer or more cells than the header (padded/expanded to match),
alignment markers (`:---`, `:---:`, `---:`), and escaped pipes (`\|`)
inside a cell.

## Status

Early. The parser targets the common GitHub-flavored-markdown table shape
and does not yet handle every edge case a hand-written table can produce.
See the test suite in `mdtable_test.go` for what's covered so far.

## License

MIT, see [LICENSE](LICENSE).
