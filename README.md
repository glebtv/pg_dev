## PG_DEV

Tool to optimize various things during app development with postgresql

### WARNING

DO NOT USE ON PRODUCTION SERVERS IN ANY WAY

This is a tool for developers, who don't have anything valuable in their DB.

This tool can easily delete ALL YOUR DATA, there is NO PROMPTS OR CONFIRMS.

## Installation

Assuming you have a working Go environment and `GOPATH/bin` is in your
`PATH`, `pg_dev` is a breeze to install:

```shell
go get github.com/glebtv/pg_dev
```

Then verify that `pg_dev` was installed correctly:

```shell
pg_dev -h
```

## Usage

```
pg_dev reset_schema
```

## Help

```
pg_dev --help
```
