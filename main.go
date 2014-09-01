/*
 * This file is part of Cashbook, a tool to analyze and report on sets of financial transactions.
 *
 * Copyright (C) 2014  Sourdough Labs Research and Development Corp.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

// Cashbook is a program to analyse a series of formatted financial 
// transactions into interesting sets for later reporting.
package main

import (
	"database/sql"
	"io"
	"flag"
	"fmt"
    "log"
	"os"
	"strings"
	"text/template"
    "time"
	"unicode"
    "unicode/utf8"

	"bitbucket.org/liamstask/goose/lib/goose"
    "github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
)

// global options. available to any subcommands.
var flagPath = flag.String("path", "conf", "folder containing config files")
var flagEnv = flag.String("env", "development", "which DB environment to use")
var flagMatchers = flag.String("matchers", "conf/matchers.json", "file containing the matcher defs")

// helper to create a DBConf from the given flags
func dbConfFromFlags() (dbconf *goose.DBConf, err error) {
	return goose.NewDBConf(*flagPath, *flagEnv)
}

const Version = "2014.9"

var commands= []*Command{
	importCsvCmd,
	reportCmd,
	upCmd,
	downCmd,
	redoCmd,
	statusCmd,
}

func main() {

	flag.Usage = usage
	flag.Parse()
	
	args := flag.Args()
	if len(args) == 0 || args[0] == "-h" {
		flag.Usage()
		return
	}
	
	if args[0] == "help" {
		help(args[1:])
		return
	}

	if args[0] == "version" {
		fmt.Printf("Cashbook Version: %s \n", Version)
		return
	}
	

	var cmd *Command
	name := args[0]
	for _, c := range commands {
		if strings.HasPrefix(c.Name, name) {
			cmd = c
			break
		}
	}
	
	if cmd == nil {
		fmt.Printf("error: unknown command %q\n", name)
		flag.Usage()
		os.Exit(1)
	}

	cmd.Exec(args[1:])
}

// Initialize a connection to the database
func initDb() *gorp.DbMap {
	conf, err := dbConfFromFlags()
	if err != nil {
		log.Fatal(err)
	}

	//We're going to insist on Postgres here, the open string will come from the 
	// dbConf object.
	db, err := sql.Open("postgres", conf.Driver.OpenStr)
    checkErr(err, "sql.Open failed")

    // construct a gorp DbMap
    dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	// -- Populate the dbmap with table definitions.

    // add a table, setting the table name to 'posts' and
    // specifying that the Id property is an auto incrementing PK
    dbmap.AddTableWithName(Txn{}, "txns").SetKeys(true, "Id")
    dbmap.AddTableWithName(TxnGroup{}, "txn_groups").SetKeys(true, "Id")

    return dbmap
}

func checkErr(err error, msg string) {
    if err != nil {
        log.Fatalln(msg, err)
    }
}

func usage() {
	fmt.Print(usagePrefix)
	flag.PrintDefaults()
	usageTmpl.Execute(os.Stdout, commands)
}

var usagePrefix = `Cashbook (c) 2014 Sourdough Labs Research and Development Corp.

Usage:
    cashbook [options] <command> [subcommand options]

Options:
`

var usageTmpl = template.Must(template.New("usage").Parse(`
Commands are:{{range .}}
    {{.Name | printf "%-20s"}} {{.Summary}}{{end}}
    {{"version" | printf "%-20s"}} Print the version number and exit

Use "cashbook help [command]" for more information about a command.

`))

var helpTemplate = `{{.Help | trim}}
`

func capitalize(s string) string {
	if s == "" {
		return s
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToTitle(r)) + s[n:]
}

// tmpl executes the given template text on data, writing the result to w.
func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	t.Funcs(template.FuncMap{"trim": strings.TrimSpace, "capitalize": capitalize})
	template.Must(t.Parse(text))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}

func help(args []string) {
	if len(args) == 0 {
		usage()
		// not exit 2: succeeded at 'go help'.
		return
	}
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "usage: cashbook help command\n\nToo many arguments given.\n")
		os.Exit(2) // failed at 'go help'
	}

	arg := args[0]
	
	for _, cmd := range commands {
		if cmd.Name == arg {
			fmt.Printf("usage: cashbook %s", cmd.Usage)
			fmt.Print("\n\nOptions:\n\n")
			cmd.Flag.PrintDefaults()
			fmt.Print("\n")

			tmpl(os.Stdout, helpTemplate, cmd)
			// not exit 2: succeeded at 'go help cmd'.
			return
		}
	}
	
	fmt.Fprintf(os.Stderr, "Unknown help topic %#q.  Run 'cashbook help'.\n", arg)
	os.Exit(2) // failed at 'go help cmd'
}


func timeTrack(start time.Time, name string) {
    elapsed := time.Since(start)
    log.Printf("%s took %s", name, elapsed)
}
