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
pg_dev r app_development --user app
```

## Help / Options

```
pg_dev --help
```

```
NAME:
   pg_dev - PostgreSQL dev tool

USAGE:
   pg_dev [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
     reset_schema, r  Drop schema, create schema
     help, h          Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --host value      postgresql host (default: "localhost") [$PGHOST]
   --password value  postgresql password (default: "postgres") [$PGPASSWORD]
   --port value      postgresql port (default: 5432) [$PGPORT]
   --user value      postgresql user (default: "postgres") [$PGUSER]
   --help, -h        show help
   --version, -v     print only the version
```

```
> pg_dev r --help
```

```
NAME:
   pg_dev reset_schema - Drop schema, create schema

USAGE:
   pg_dev reset_schema [command options] [arguments...]

OPTIONS:
   --schema value, -s value  Owner name (default: "public")
   --user value, -u value    Owner name
   --no_drop                 Don't drop, just create
   --no_create               Don't create, just drop
```

