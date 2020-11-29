# sibylgo

![CI](https://github.com/sandro-h/sibylgo/workflows/CI/badge.svg)

Text-based TODO application

## Text syntax

### General notes

* Indentation is currently restricted to tab characters.

### Category

```text
------------------
 My category
------------------
```

With HTML color (to differentiate calendar entries):

```text
------------------
 My category [Coral]
------------------
```

### Todo

```text
[] my open todo
[x] my done todo
```

Hierarchical todos:

```text
[] my top-level todo
    [] my child todo
    [x] my done child todo
        [] more stuff
```

Comments:

```text
[] my top-level todo
    some random comment.
    it needs to be indented.

    [] child todo
        a comment for the child todo
```

On a specific day:

```text
[] get groceries (15.11.20)
```

At a specific date and time:

```text
[] get groceries (15.11.20 08:00)
```

Time range:

```text
[] vacation (10.11.20-15.11.20)
[] new house (10.11.20-)
[] study for example (-15.11.20)
```

Recurring:

```text
[] get groceries (every day)
[] get groceries (today)

[] gym (every Tuesday)
[] team meeting (every 2nd Tuesday)
[] project meeting (every 3rd Tuesday)
[] company meeting (every 4th Tuesday)

[] vacuum (every 10.)

[] John's birthday (every 5.10)
```

Important (!):

```text
[] pick up from airport!
[] even more important!!
```

## Configuration

### sibylgo.yml

```yaml
todoFile: path/to/todo.txt
host: localhost

mailHost: smtp.example.com
mailPort: 3025
mailFrom: foo@example.com 
mailTo: bar@example.com
mailUser: foo
mailPassword: lepass

external_sources:
  prepend: true
  bitbucket_prs:
    bb_url: http://bitbucket.example.com
    bb_user: myuser
    bb_token: aba1234
    category: Today
  dummies:
    dummy_moments:
      - id1:name1
      - id2:name2
```

## Development

Main Go application:

```shell
make deps
make build
```

VSCode Extension:

```shell
cd vscode_ext
make deps
make build
```

### Testing

Some tests rely on input/output testdata files. The output files for all tests can be updated to what
the test is actually outputting with: `go test ./... -update-golden`

For easier comparison between actual and expected output, the test output
can also be written to a temporary file instead of the real golden file:
`go test ./... -update-golden -dry-golden`

### Releasing

```shell
make release
```

Pushes a release tag which triggers a CI build that includes a release job.
