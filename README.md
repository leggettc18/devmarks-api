# Devmarks API

![Devmarks Logo](https://raw.githubusercontent.com/leggettc18/devmarks-frontend-web/main/src/assets/logo.svg)

This is the README for the backend API of Devmarks.

Devmarks will eventually be a Web App to allow developers to organize
Bookmarks amongst their team, organizing them with Folders, Organizations,
Tags, and Colors. Currently it only does Bookmarks and Users, but implementing
the rest will mostly be a repetition of existing patterns.

## Requirements

The following are the requirements specifically for the backend. The frontend
may have its own set of requirements.

- PostgresQL Database
- Golang

## Setup

1. Install Postgresql database onto your host system and configure a database
and users for the app. The exact names do not matter as long as it matches
the configuration in step 5. Such configuration is out of the scope
of this documentation.
2. Clone the repository.

    ```bash
    git clone https://github.com/leggettc18/devmarks-api
    ```

3. Rename `config.example.yaml` to `config.yaml`
4. Supply a randomly generated Secret Key.
5. Supply the necessary database information according to the example format.
6. Build the project. Feel free to supply a different executable name after the
`-o` flag if desired.

    ```bash
    go build -i main.go -o devmarks
    ```

7. Run the migrations to set up database tables.

    ```bash
    ./devmarks migrate up
    ```

8. Run the `serve` command. Optionally provide the `--config` flag if the
config file is either named differently and/or not in the same folder as the
executable.

    ```bash
    ./devmarks --config <path to config.yaml> serve
    OR
    ./devmarks serve
    ```