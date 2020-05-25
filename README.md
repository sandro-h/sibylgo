# sibylgo

[![CircleCI](https://circleci.com/gh/sandro-h/sibylgo.svg?style=svg&circle-token=9e65f022c014e5685c7fbd76148892f711d58bed)](https://circleci.com/gh/sandro-h/sibylgo)

Text-based TODO application

# Configuration

**sibylgo.yml**
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

# Development

Main Go application:
```
make deps-go
make build-go
```

VSCode Extension:
```
make deps-vscode
make build-vscode
```

Testing:  
Some tests rely on input/output testdata files. The output files for all tests can be updated to what
the test is actually outputting with: `go test ./... -update-golden`

For easier comparison between actual and expected output, the test output
can also be written to a temporary file instead of the real golden file:
`go test ./... -update-golden -dry-golden`


## See vscode extension errors

When packaging, you won't see compile errors, so simply run `npm run compile` to see them.

# TODOS

- add a space before/after insert
