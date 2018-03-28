package main

import (
	"database/sql"
	"os"
	"sort"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
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

func connect() (sql.Db, error) {
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
					_, err := db.Query("DROP SCHEMA $1 CASCADE", schema)
					if err != nil {
						return cli.NewExitError("unable to drop schema "+schema+": "+err.Error(), 4)
					}

					user := c.String("user")
					_, err := db.Query("CREATE SCHEMA $1 OWNED BY $2", schema, user)
					if err != nil {
						return cli.NewExitError("unable to create schema "+schema+" with owner "+user+": "+err.Error(), 5)
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
