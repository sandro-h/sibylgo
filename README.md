# sibylgo

![CI](https://github.com/sandro-h/sibylgo/workflows/CI/badge.svg) [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=sandro-h_sibylgo&metric=alert_status)](https://sonarcloud.io/dashboard?id=sandro-h_sibylgo)

Text-based TODO application.

It lets you write TODO list as you want (within the bounds of the syntax).
You can add arbitrary additional text and comments to the todos. The application will never reformat your TODO file, so you have full control of it.

**Disclaimer:** This tool is for my own personal use. Feedback and requests are welcome, but not necessarily acted upon.

## Components

### Backend

The backend does all the heavy lifting:

* Parsing the text todo file
* Providing formatting information for the VSCode extension
* Providing folding information for the VSCode extension
* Providing data for HTML calendar
* Providing commands to clean and trash done todos
* Sending mail reminders for upcoming todos
* Creating backups of the todo file
* Inserting todos from external sources (like open Bitbucket PRs)

The backend is a `sibylgo.exe` (Windows) or `sibylgo` (Linux) console application that can be started in the background somewhere. All the other components interact with it via REST calls.

### VSCode extension

The VSCode extension is a thin client that interacts with the backend
and formats the todo file if opened in VSCode.  
Also handles folding of hierarchical todos and provides the clean and
trash commands in the command palette and editor context menu.

The extension is a `sibyl.vsix` file that can be installed manually with:

```shell
code --install-extension sibyl.vsix
```

Or via the "Install from VSIX..." option in the VSCode extension GUI.

### Calendar

The calendar is a simple `sibylcal.html` file that displays the
current month/week/day, using data from the backend.

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

Work states:

```text
[] my open todo
[p] my in progress todo
[w] my waiting todo
[x] my done todo
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
log_level: info
todoFile: path/to/todo.txt
host: localhost

optimized_format: true

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

outlook_events:
  enabled: true
```

## Development

Main Go application:

```shell
make deps
make build
```

VSCode Extension:

```shell
cd VSCode_ext
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

Make sure to update and commit `version.txt`.

```shell
make release
```

Pushes a release tag which triggers a CI build that includes a release job.
