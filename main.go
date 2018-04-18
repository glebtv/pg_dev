package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
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

func Execute(command string, arg ...string) {
	cmd := exec.Command(command, arg...)

	stdoutIn, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	stderrIn, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}

	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	go func() {
		io.Copy(os.Stdout, stdoutIn)
	}()

	go func() {
		io.Copy(os.Stderr, stderrIn)
	}()

	err = cmd.Wait()
	if err != nil {
		log.Fatalf(command+" "+strings.Join(arg, " ")+" failed with %s\n", err)
	}
}

func getConnStr() string {
	connstr := make([]string, 0)
	if host != "" {
		connstr = append(connstr, "host="+Quote(host))
	}
	connstr = append(connstr, "port="+strconv.Itoa(port))
	connstr = append(connstr, "user="+Quote(user))

	if dbname != "" {
		connstr = append(connstr, "dbname="+Quote(dbname))
	}

	if password != "" {
		connstr = append(connstr, "password="+Quote(password))
	}

	connstr = append(connstr, "sslmode=disable")
	return strings.Join(connstr, " ")
}

func Quote(name string) string {
	end := strings.IndexRune(name, 0)
	if end > -1 {
		name = name[:end]
	}
	return `'` + strings.Replace(name, `'`, `\'`, -1) + `'`
}

func connect() (*sql.DB, error) {
	connStr := getConnStr()
	fmt.Println("connecting:", connStr)
	return sql.Open("postgres", connStr)
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
			Value:       "",
			EnvVar:      "PGPASSWORD",
			Destination: &password,
		},
		cli.BoolFlag{
			Name:  "migrate",
			Usage: "Run rails migrations",
		},
		cli.BoolFlag{
			Name:  "seed",
			Usage: "Run rails seeds",
		},
	}

	App.After = func(c *cli.Context) error {
		path, err := exec.LookPath("bundle")
		if err != nil {
			log.Fatal("bundler not found")
			return nil
		}

		fmt.Printf("bundler is available at %s\n", path)

		migrate := c.Bool("migrate")
		if migrate {
			fmt.Printf("running migrations\n")
			Execute(path, "exec", "rake db:migrate")
		}
		seed := c.Bool("seed")
		if seed {
			fmt.Printf("running seeds\n")
			Execute(path, "exec", "rake db:seed")
		}
		fmt.Printf("done!\n")
		return nil
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
					Name:  "owner, o",
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
				if c.NArg() <= 0 {
					return cli.NewExitError("no db name provided.", 1)
				}

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
					q := fmt.Sprintf("DROP SCHEMA %s CASCADE", quoted)
					log.Println(q)
					_, err = db.Exec(q)
					if err != nil {
						return cli.NewExitError("unable to drop schema "+schema+": "+err.Error(), 4)
					}
				}

				if !no_create {
					user := c.String("owner")
					quoted_user := pq.QuoteIdentifier(user)
					if user != "" {
						q := fmt.Sprintf("CREATE SCHEMA %s AUTHORIZATION %s", quoted, quoted_user)
						log.Println(q)
						_, err = db.Exec(q)
					} else {
						q := fmt.Sprintf("CREATE SCHEMA %s", quoted)
						log.Println(q)
						_, err = db.Exec(q)
					}
					if err != nil {
						return cli.NewExitError("unable to create schema "+schema+" with owner "+user+": "+err.Error(), 5)
					}
				}

				return nil
			},
		},
		{
			Name:    "create",
			Aliases: []string{"c"},
			Usage:   "Create user with password {user}, create database {user}_development, and grant him full privileges",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "set_password",
					Usage: "Set new user password, default {user}",
				},
				cli.StringFlag{
					Name:  "dbname, db",
					Usage: "Database name, default {user}_development",
				},
				cli.StringFlag{
					Name:  "auth_db_name",
					Usage: "Authentication database name, default postgres",
					Value: "postgres",
				},
			},
			Action: func(c *cli.Context) error {
				if c.NArg() <= 0 {
					return cli.NewExitError("no user name provided.", 1)
				}

				uname := c.Args().Get(0)
				if uname == "" {
					return cli.NewExitError("no user name provided.", 2)
				}

				dbname = c.String("auth_db_name")
				db, err := connect()
				if err != nil {
					return cli.NewExitError("unable to connect to postgresql: "+err.Error(), 3)
				}

				password := c.String("set_password")
				if password == "" {
					password = uname
				}
				password_quoted := Quote(password)
				q := fmt.Sprintf("CREATE USER %s WITH PASSWORD %s", uname, password_quoted)
				log.Println(q)
				_, err = db.Exec(q)
				if err != nil {
					log.Println("unable to create user " + uname + ": " + err.Error())
					//return cli.NewExitError("unable to create user "+uname+": "+err.Error(), 4)
				}

				new_db_name := c.String("dbname")
				if new_db_name == "" {
					new_db_name = uname + "_development"
				}
				dbname_quoted := pq.QuoteIdentifier(new_db_name)
				uname_quoted := pq.QuoteIdentifier(password)

				q = fmt.Sprintf("CREATE DATABASE %s", dbname_quoted)
				log.Println(q)
				_, err = db.Exec(q)
				if err != nil {
					return cli.NewExitError("unable to create db "+dbname+": "+err.Error(), 5)
				}

				q = fmt.Sprintf("GRANT ALL ON DATABASE %s TO %s;", dbname_quoted, uname_quoted)
				log.Println(q)
				_, err = db.Exec(q)
				if err != nil {
					return cli.NewExitError("unable to grant all on db "+dbname+" to user "+uname+": "+err.Error(), 6)
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
