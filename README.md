# EarthRanger CLI 
![Build Status](https://github.com/doneill/er-cli-go/actions/workflows/go.yml/badge.svg)

## Build
`go build -o bin/er`

## Commands

### Authenticate

```bash
er auth [flags]
er auth [command]

Available commands
token Display token

Flags:
-s, --sitename string EarthRanger site name
-u, --username string EarthRanger username
```

**Examples:**

```bash
# Authenticate a user and sitename
gh auth -s {sitename} -u {username}

# Display token once authenticated
gh auth token
```

### User Details

```bash
# Display the currently authenticated user
er user
```


## Contributors
<a href="https://github.com/doneill/er-cli-go/graphs/contributors">
  <img src="https://contributors-img.web.app/image?repo=doneill/er-cli-go" />
</a>

## Licensing
A copy of the license is available in the repository's [LICENSE](LICENSE) file.
