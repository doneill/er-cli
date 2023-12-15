# EarthRanger CLI 
![Build Status](https://github.com/doneill/er-cli/actions/workflows/go.yml/badge.svg)

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
# Returns table data
| USERNAME |        EMAIL         | FIRST NAME | LAST NAME |                  ID                  | PIN  |              SUBJECT ID              |
|----------|----------------------|------------|-----------|--------------------------------------|------|--------------------------------------|
| cccy     |                      | CC         |  CY       | 015945ff-c220-4674-a070-3f1112e445fg |      | 12c245f6-8d77-4e15-a82c-be4a717034df |
```


## Contributors
<a href="https://github.com/doneill/er-cli/graphs/contributors">
  <img src="https://contributors-img.web.app/image?repo=doneill/er-cli" />
</a>

## Licensing
A copy of the license is available in the repository's [LICENSE](LICENSE) file.
