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
	"os"
	"io"
	"encoding/csv"
	"log"
	"time"
	"fmt"
	"strconv"
	"strings"
)

var importCsvCmd = &Command{
	Name:    "importcsv",
	Usage:   "",
	Summary: "Import a csv file",
	Help:    `import csv extended help here...`,
	Run:     importCsvRun,
}

const (
	TXNTYPE  = iota
	AMOUNT
	DESC
	OCCURREDAT
	TXNHOSTTYPE
	TRACENUMBER
	COMBINATIONKEY
	ACCOUNTGROUPID
)

// Date format:  RFC3339     = "2006-01-02T15:04:05Z07:00"

func importCsvRun(cmd *Command, args ...string) {

	log.Printf("Import CSV")

	if len(args) == 0 {
		log.Print("Missing filename of csv file, exiting")
		return
	}

	dbm := initDb()
	defer dbm.Db.Close()

	log.Printf("DB Connected.")

	matchers := match_set_from_file(*flagMatchers)

	file, err := os.Open(args[0])
    if err != nil {
        log.Println("Error:", err)
        return
    }
    defer file.Close()
    reader := csv.NewReader(file)
	reader.Comma = '\t'// Use tab-separated values

	defer timeTrack(time.Now(), "CSV Import")

	count := 0
	err_count := 0

    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        } else if err != nil {
            log.Println("Error:", err)
            return
        }
 		
		amount, err := strconv.ParseInt(strings.TrimSpace(record[AMOUNT]), 10, 64) // Pennies!

		if err != nil {
			log.Printf("error parsing amount: %v: record =  %v", err, record)
			err_count++
			continue
		}

		occurred_at, err := time.Parse("02 Jan 06 15:04:05", strings.TrimSpace(record[OCCURREDAT]))

		if err != nil {
			log.Printf("error parsing occured_at: %v: record =  %v", err, record)
			err_count++
			continue
		}

		txn := &Txn{
			TxnType: strings.TrimSpace(record[TXNTYPE]),
			Amount: amount, // Remember, in pennies.
			Description: strings.TrimSpace(record[DESC]),
			OccurredAt: occurred_at,
			TxnHostType: strings.TrimSpace(record[TXNHOSTTYPE]),
			TraceNumber:  strings.TrimSpace(record[TRACENUMBER]),
			CombinationKey: strings.TrimSpace(record[COMBINATIONKEY]),
			AccountGroupId: strings.TrimSpace(record[ACCOUNTGROUPID]),
		}

		_ = assign_txn_to_txn_group(txn, matchers, dbm)
		count++

		if count % 500 == 0 {
			fmt.Print(".")  // Show progress every 1000 records
		}
    }

	log.Printf("Imported %d records, %d bad records", count, err_count)
}










