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

package main

import (
	"log"
	"time"
	"strings"
	"path"
	"html/template"
	"os"
	"bufio"
)

var reportUsage = "report [-membersfile=filename | -id=accountgroupid] | -system] -output=path -start=yyyy-mm-dd -end=yyyy-mm-dd template"

var reportCmd = &Command{
	Name:    "report",
	Usage:   reportUsage,
	Summary: "Generate reports",
	Help:    `
One of membersfile, system or id is required as is template.`,
	Run:     reportRun,
}

var membersFile string
var accountGroupId string
var systemReport bool
var outputDir string
var startDate  string
var endDate  string

func init() {
	reportCmd.Flag.StringVar(&membersFile, "membersfile", "", "A file containing a list of account group ids.")
	reportCmd.Flag.StringVar(&accountGroupId, "id", "", "Generate a report for a single id.")
	reportCmd.Flag.BoolVar(&systemReport, "system",false, "Generate a report for the aggregate of all ids.")
	reportCmd.Flag.StringVar(&outputDir, "output", ".", "Directory to store generated reports in.")
	reportCmd.Flag.StringVar(&startDate, "start", "", "Sets the start date to report on.")
	reportCmd.Flag.StringVar(&endDate, "end", "", "Sets the end date to report on.")
}

func reportRun(cmd *Command, args ...string) {

	if systemReport == false && membersFile == "" && accountGroupId == "" {
		printError("Missing option: One of membersfile, id or system is required\n", reportUsage)
		return
	}

	if endDate == "" || startDate == "" {
		printError("Missing option: -start and -end are required\n",reportUsage)
		return
	}

	if len(args) == 0 {
		printError("Missing template filename\n", reportUsage)
		return
	}

	var report = args[0]

	dbm := initDb()
	defer dbm.Db.Close()

	log.Print("Cashbook Report Generator")

	log.Printf("DB Connected.")

	defer timeTrack(time.Now(), "Generate Reports")

	log.Printf("Generating reports using '%s'", report)

	var pieces = strings.Split(report, ".")
	// Don't care about what's in-between really 
	// ext is just for the output file to match the template
	var ext = pieces[len(pieces) - 1] 


	var t = template.Must(template.New(report).Funcs(template.FuncMap{"currency": currency, 
		"capitalize": capitalize}).ParseFiles(report))

	reportStart, _ := time.Parse(shortForm, startDate + " 00:00:00")
	reportEnd, _ := time.Parse(shortForm, endDate + " 23:59:59")

	// TODO Add support for reading id's from a file.

	// TODO Add support for running the report for the organization.

	var id = accountGroupId

	var api = &ReportingApi{reportStart, reportEnd, id, dbm}

	// Do work.... TODO Lame, only supports a single id at the moment....
	renderReportToFile(t, ext, api)
}

func renderReportToFile(t *template.Template, ext string, api *ReportingApi) {

	var output = path.Join(outputDir, api.AccountGroupId + "." + ext);
	log.Printf("Generating %s", output)

	f, err := os.Create(output)
	if err != nil {
		log.Printf("Error creating '%s': %v", output, err)
		return
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	t.Execute(w, api)
	w.Flush()
}
