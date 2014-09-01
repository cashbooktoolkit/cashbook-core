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
	"database/sql"
    "github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
	"time"
)

// Process a single incoming transaction
func assign_txn_to_txn_group(txn *Txn, matchers MatchSet, dbm *gorp.DbMap) (*Txn) {

	if !txn.Valid() {
		log.Printf("Invalid txn: %v", txn)
		return txn
	}

	matcher := matchers.FindMatcher(txn.Description)

	if matcher != nil {

		label  := matcher.Label(txn.Description)

		group_type :=  matcher.GetGroupType() // Retail, Bill, Loan, Transfer, etc
		
		system_txn_group := find_or_create_system_txn_group(label, group_type, dbm)

		txn.SystemTxnGroupId = system_txn_group.Id
		txn.Classification = system_txn_group.Classification

		txn_group := find_or_create_txn_group(txn.AccountGroupId, system_txn_group, dbm)
		
		txn.TxnGroupId = txn_group.Id
		txn.TxnGroupType = txn_group.GroupType
		txn.CategoryId = txn_group.CategoryId
	}

	if txn.Id != 0 { // TODO Find a better way to do this.
		_, err := dbm.Update(txn)
		if err != nil {
			log.Printf("Error updating txn: %v", err)
		}
	} else {
		txn.Created = time.Now().UnixNano()
		err := dbm.Insert(txn)

		if err != nil {
			log.Printf("Error inserting txn: %v", err)
		}
	}

	return txn
}

func find_or_create_txn_group(agid string, system_group *TxnGroup, dbm *gorp.DbMap) (txn *TxnGroup) {

	// Pre-initialize a new record
	group := TxnGroup{
		GroupType: system_group.GroupType,
		Label: system_group.Label, 
		Classification: system_group.Classification, 
		SystemTxnGroupId: system_group.Id, 
		AccountGroupId: agid,
		Created: time.Now().UnixNano(),
	}

	err := dbm.SelectOne(&group, 
		"select * from txn_groups where account_group_id = :agi and group_type = :group_type and label = :label", 
		map[string]interface{} { 
			"agi": agid, 
			"group_type": system_group.GroupType, 
			"label": system_group.Label})

	if err != nil {
		if err ==  sql.ErrNoRows {
			err = dbm.Insert(&group)

			if err != nil {
				log.Printf("Error inserting new txn group %v", err)
			}

			// TODO Handle UNIQUE constraint violation (the 'Race Condition') by re-trying
		} else {
			log.Printf("Error selecting txn group %v", err)
		}
	}

	return &group
}


func find_or_create_system_txn_group(label string, group_type string, dbm *gorp.DbMap) (txnGroup *TxnGroup) {
	// Pre-initialize a new record
	group := TxnGroup{
		GroupType: group_type,
		Label: label, 
		Created: time.Now().UnixNano(),
	}

	err := dbm.SelectOne(&group, 
		"select * from txn_groups where account_group_id = '' and group_type = :group_type and label = :label", 
		map[string]interface{} { 
			"group_type": group.GroupType, 
			"label": group.Label})

	if err != nil {
		if err ==  sql.ErrNoRows {
			err = dbm.Insert(&group)

			if err != nil {
				log.Printf("Error inserting new system txn group %v", err)
			}

			// TODO Handle UNIQUE constraint violation (the 'Race Condition') by re-trying
		} else {
			log.Printf("Error selecting system txn group %v", err)
		}
	}

	return &group
}



/*
		if err, ok := err.(*pq.Error), ok {
			fmt.Println("pq error:", err.Code.Name())

			err := dbm.SelectOne(&group, 
				"select * from txn_groups where account_group_id = :agi and group_type = :group_type and label = :label", 
				map[string]interface{} { 
					"agi": agid, 
					"group_type": system_group.GroupType, 
					"label": system_group.Label})
			
			if err == sql.ErrNoRows {
				log.Println("Got postgress error, tried select one again, and again got ErrNoRows (very strange)")
			}
		}
*/


















