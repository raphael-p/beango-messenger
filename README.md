# Beango Messenger

![banner for BeanGo messanger](banner.png)

Beango is a messaging app built using Golang and HTMX.

## Approach

### Minimal dependencies

This project uses the bare minimum dependencies, so a lot of things are custom-built from scratch. This does add quite a lot of effort to development, but should provide the benefit of making things more transparent, flexible, and lightweight.

### Adhere to HTMX principles

HTMX was chosen to build to web front end for this app in order to try a different approach to web development. We want state to be managed in a single place, the server. Therefore, we use as little javascript code as possible, and leverage HTMX's features to keep most front end logic in the DOM itself.

## How to run

1. Clone `config/default.docker.env` and rename the copy to `docker.env`.
2. Set BG_DB_PASSWORD to a value of your choice. This will be the password used for postgres users.

### With Docker (recommended)

Prequisite: Docker must be installed on your machine

3. Run `docker-compose up`.

### Without Docker

Prerequisites: Golang, Postgres

3. Run `set -a; source config/docker.env; set +a;` to make the envars available to the current shell.
4. Create a postgres database, along with a username+password with full priviledges on it. Make sure to use the same values as the BG_DB_NAME, BG_DB_USERNAME, and BG_DB_PASSWORD envars.
5. Run `go run main.go`