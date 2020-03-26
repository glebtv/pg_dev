## PG_DEV

Tool to optimize various things during app development with postgresql (mostly Ruby On Rails)

### WARNING

DO NOT USE ON PRODUCTION SERVERS IN ANY WAY

This is a tool for developers, who don't have anything valuable in their DB.

This tool can easily delete ALL YOUR DATA, there is NO PROMPTS OR CONFIRMS.

## Installation

Assuming you have a working Go environment and `GOPATH/bin` is in your
`PATH`, `pg_dev` is a breeze to install:

```shell
GO111MODULE=on go get github.com/glebtv/pg_dev
```

Then verify that `pg_dev` was installed correctly:

```shell
pg_dev -h
```

## Changelog

- 0.1.0 Initial Version
- 0.2.0 Reset now uses same syntax as create (user name only)
- 0.2.1 support hstore flag
- 0.3.0 support env toggle and RAILS_ENV env variable
- 0.3.1 dont allow postgres as db name

## Usage

### Create

Create user and db for development

```
pg_dev c app
```

Creates user app with password app, and app_development database for him.

### Create with env

Create user and db for development

```
pg_dev --env test c app
```

Creates user app with password app, and app_test database for him.

### Reset

Drop schema public, create schema public owned by correct user:

**(deletes all data in this database)**

```
pg_dev --migrate --seed r app_development --owner app
```

## Help / Options

```
pg_dev --help
```

```
NAME:
   pg_dev - PostgreSQL dev tool 

USAGE:
   main [global options] command [command options] [arguments...]

VERSION:
   0.2.0

COMMANDS:
   create, c  Create user with password {user}, create database {user}_development, and grant him full privileges
   reset, r   Reset public schema for {user}_development database
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --auth_db value   authentication database name, default postgres (default: "postgres")
   --host value      postgresql host (default: "localhost") [$PGHOST]
   --migrate         Run rails migrations
   --password value  postgresql password [$PGPASSWORD]
   --port value      postgresql port (default: 5432) [$PGPORT]
   --seed            Run rails seeds
   --user value      postgresql user (default: "postgres") [$PGUSER]
   --help, -h        show help
   --version, -v     print only the version
```

```
> pg_dev r --help
```

```
NAME:
   main reset - Reset public schema for {user}_development database

USAGE:
   main reset [command options] [arguments...]

OPTIONS:
   --schema value, -s value    Schema name (default: "public")
   --dbname value, --db value  Database name, default {user}_development
   --no_drop                   Don't drop, just create
   --no_create                 Don't create, just drop
```

```
> pg_dev c --help
NAME:
   main create - Create user with password {user}, create database {user}_development, and grant him full privileges

USAGE:
   main create [command options] [arguments...]

OPTIONS:
   --set_password value        Set new user password, default {user}
   --dbname value, --db value  Database name, default {user}_development
 ```
