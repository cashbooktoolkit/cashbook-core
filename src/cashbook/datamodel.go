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
	"time"
)

// A financial transaction belonging to a set of related accounts (Member)
type Txn struct {
	Id int64
	TxnType string `db:"txn_type"`

	Amount int64
	Description string

	OccurredAt time.Time `db:"occurred_at"`
	
	TxnHostType string `db:"txn_host_type"`
	TraceNumber string `db:"trace_number"`
	CombinationKey string `db:"combination_key"`

	CategoryId int64 `db:"category_id"`

	SystemTxnGroupId int64 `db:"system_txn_group_id"`
	TxnGroupId int64 `db:"txn_group_id"`
	TxnGroupType string `db:"txn_group_type"`

	Classification string
	AccountGroupId string `db:"account_group_id"`
	
	Created int64 // Init with time.Now().UnixNano

	// These are used to route incoming Txn's and don't get persisted
	ActionCode string `db:"-"`
	Version string `db:"-"`
}

// Check to ensure the txn is a valid one.
func (txn *Txn) Valid() bool {
	return true // TODO Implement this method.
}

// A group of related transactions. Also references things that 
// would apply to all transactions in the group. 
type TxnGroup struct {
	Id int64
	GroupType string `db:"group_type"`

	Label string `db:"label"`  // This is the label/fingerprint

	Description string // If user sets will be used for display
	Classification string

	SystemTxnGroupId int64 `db:"system_txn_group_id"`
	CategoryId int64 `db:"category_id"`

	AccountGroupId string `db:"account_group_id"`
	
	Created int64 // Init with time.Now().UnixNano
}
