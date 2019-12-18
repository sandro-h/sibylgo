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

## See vscode extension errors

When packaging, you won't see compile errors, so simply run `npm run compile` to see them.

# TODOS

- command to rerun ext sources
- option to insert at start of category
- add a space before/after insert
- use git as backup mechanism -> create a pre/post commit when auto-modifying the todo file.
