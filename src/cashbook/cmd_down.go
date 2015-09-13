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
	"bitbucket.org/liamstask/goose/lib/goose"
	"log"
)

var downCmd = &Command{
	Name:    "migrate:down",
	Usage:   "",
	Summary: "Roll back the version by 1",
	Help:    `down extended help here...`,
	Run:     downRun,
}

func downRun(cmd *Command, args ...string) {

	conf, err := dbConfFromFlags()
	if err != nil {
		log.Fatal(err)
	}

	current, err := goose.GetDBVersion(conf)
	if err != nil {
		log.Fatal(err)
	}

	previous, err := goose.GetPreviousDBVersion(conf.MigrationsDir, current)
	if err != nil {
		log.Fatal(err)
	}

	if err = goose.RunMigrations(conf, conf.MigrationsDir, previous); err != nil {
		log.Fatal(err)
	}
}
