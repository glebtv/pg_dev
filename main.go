package main

import (
	"database/sql"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/lib/pq"
	"github.com/urfave/cli"
)

var App *cli.App

var host = ""
var port = 5432
var user = "postgres"
var password = ""
var dbname = ""

func getConnStr() string {
	connstr := make([]string, 0)
	connstr = append(connstr, "host="+host)
	connstr = append(connstr, "port="+strconv.Itoa(port))
	connstr = append(connstr, "user="+user)
	connstr = append(connstr, "dbname="+dbname)
	connstr = append(connstr, "password="+password)
	connstr = append(connstr, "sslmode=disable")
	return strings.Join(connstr, " ")
}

func connect() (*sql.DB, error) {
	return sql.Open("postgres", getConnStr())
}

func init() {
	App = cli.NewApp()
	App.EnableBashCompletion = true
	App.Name = "pg_dev"
	App.Usage = "PostgreSQL dev tool "
	App.Version = "0.1.0"

	cli.VersionFlag = cli.BoolFlag{
		Name:  "version, v",
		Usage: "print only the version",
	}

	App.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "host",
			Usage:       "postgresql host",
			Value:       "localhost",
			EnvVar:      "PGHOST",
			Destination: &host,
		},
		cli.IntFlag{
			Name:        "port",
			Usage:       "postgresql port",
			Value:       5432,
			EnvVar:      "PGPORT",
			Destination: &port,
		},
		cli.StringFlag{
			Name:        "user",
			Usage:       "postgresql user",
			Value:       "postgres",
			EnvVar:      "PGUSER",
			Destination: &user,
		},
		cli.StringFlag{
			Name:        "password",
			Usage:       "postgresql password",
			Value:       "postgres",
			EnvVar:      "PGPASSWORD",
			Destination: &password,
		},
	}

	App.Commands = []cli.Command{
		{
			Name:    "reset_schema",
			Aliases: []string{"r"},
			Usage:   "Drop schema, create schema",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "schema, s",
					Usage: "Owner name",
					Value: "public",
				},
				cli.StringFlag{
					Name:  "user, u",
					Usage: "Owner name",
				},
				cli.BoolFlag{
					Name:  "no_drop",
					Usage: "Don't drop, just create",
				},
				cli.BoolFlag{
					Name:  "no_create",
					Usage: "Don't create, just drop",
				},
			},
			Action: func(c *cli.Context) error {
				if c.NArg() > 0 {
					dbname = c.Args().Get(0)
					if dbname == "" {
						return cli.NewExitError("no db name provided.", 2)
					}
					db, err := connect()
					if err != nil {
						return cli.NewExitError("unable to connect to postgresql: "+err.Error(), 3)
					}

					schema := c.String("schema")
					quoted := pq.QuoteIdentifier(schema)

					no_create := c.Bool("no_create")
					no_drop := c.Bool("no_drop")

					if !no_drop {
						_, err = db.Exec(fmt.Sprintf("DROP SCHEMA %s CASCADE", quoted))
						if err != nil {
							return cli.NewExitError("unable to drop schema "+schema+": "+err.Error(), 4)
						}
					}

					if !no_create {
						user := c.String("user")
						quoted_user := pq.QuoteIdentifier(user)
						if user != "" {
							_, err = db.Exec(fmt.Sprintf("CREATE SCHEMA %s AUTHORIZATION %s", quoted, quoted_user))
						} else {
							_, err = db.Exec(fmt.Sprintf("CREATE SCHEMA %s", quoted))
						}
						if err != nil {
							return cli.NewExitError("unable to create schema "+schema+" with owner "+user+": "+err.Error(), 5)
						}
					}

				} else {
					return cli.NewExitError("no db name provided.", 1)
				}
				return nil
			},
		},
	}

	sort.Sort(cli.FlagsByName(App.Flags))
	sort.Sort(cli.CommandsByName(App.Commands))
}

func main() {
	App.Run(os.Args)
}
