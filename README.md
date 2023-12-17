# EarthRanger CLI 
![Build Status](https://github.com/doneill/er-cli/actions/workflows/go.yml/badge.svg)

## Build
`go build -o bin/er`

## Help

```bash
Work with EarthRanger platform from command line

Usage:
  er [command]

Available Commands:
  auth        Authentication with EarthRanger
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  open        Open a SQLite database file
  user        Current authenticated user data

Flags:
  -h, --help      help for er
  -t, --toggle    Help message for toggle
  -v, --version   version for er

Use "er [command] --help" for more information about a command.
```

## Remote Site Server Commands
These commands allow you to work through an authenticated user account with your EarthRanger site server.

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
er auth -s {sitename} -u {username}

# Display token once authenticated
er auth token
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

## Local Database Commands
These commands allow you to work with an exported EarthRanger mobile databse

```bash
This tool is intended to be used specficially with EarthRanger mobile databases

Usage:
  er open [sqlite db file] [flags]

Flags:
  -h, --help     help for open
  -t, --tables   Display all database tables
  -u, --user     Display database account user
```

### Open
Open an EarthRanger database file. This command is the entry to working with the database

```bash
er open earthranger.db
earthranger.db successfully opened!
```

### Tables
Display all tables and record count

```bash
er open earthranger.db -t

+------------------+-------+
|       NAME       | COUNT |
+------------------+-------+
| android_metadata |     1 |
| accounts_user    |     1 |
| sqlite_sequence  |     9 |
| user_profiles    |     2 |
| user_subjects    |     3 |
| event_type       |    51 |
| event_category   |     5 |
| events           |    13 |
| attachments      |    17 |
| sync_states      |     5 |
| patrol_types     |     4 |
| patrols          |     0 |
| patrol_segments  |     0 |
+------------------+-------+
```

### Database user
Query the database for the username of the mobile user

```bash
er open earthranger.db -u
username
```

## Contributors
<a href="https://github.com/doneill/er-cli/graphs/contributors">
  <img src="https://contributors-img.web.app/image?repo=doneill/er-cli" />
</a>

## Licensing
A copy of the license is available in the repository's [LICENSE](LICENSE) file.
