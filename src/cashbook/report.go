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
    "github.com/coopernurse/gorp"
	"log"
	"fmt"
)

const shortForm = "2006-01-02 15:04:05"

// A generlizing query api for building various kinds of reports out of PFM data
type ReportingApi struct {
	// Query for results on or after this time
	PeriodStart time.Time

	// Query for results before this time
	PeriodEnd time.Time

	// If non-nil, queries will be scoped to this account group, other wise will query FI wide
	AccountGroupId string

	dbmap *gorp.DbMap
}

type TxnSummary struct {
	Count int
	Sum int
}

func (r *ReportingApi) Remaining(deposits int, withdrawals int) int {
	return deposits - withdrawals 
}

// Return a count and sum for the set of transactions grouped in txn_type (Withdrawals and Deposits)
func (r *ReportingApi) TxnTypeSummary(txnType string) *TxnSummary {

	summary := TxnSummary{} 

	err := r.dbmap.SelectOne(&summary, "SELECT count(*) as count, coalesce(sum(amount), 0) as sum FROM txns WHERE account_group_id = :id AND txn_type = :type AND occurred_at BETWEEN :rstart AND :rend",
		map[string]interface{} { 
			"type": txnType,
			"id": r.AccountGroupId,
		    "rstart": r.PeriodStart,
			"rend": r.PeriodEnd})

	if err != nil {
		log.Printf("Error getting TxnGroupSummary: %v", err);
		return nil
	}

	return &summary
}

// Return a count and sum for the set of transactions grouped by txnGroupType (Bills, Retail, etc)
func (r *ReportingApi) TxnGroupTypeSummary(txnGroupType string) *TxnSummary {
	summary := TxnSummary{} 

	err := r.dbmap.SelectOne(&summary, "SELECT count(*) as count, coalesce(sum(amount), 0) as sum FROM txns WHERE account_group_id = :id AND txn_group_type = :type AND occurred_at BETWEEN :rstart AND :rend",
		map[string]interface{} { 
			"type": txnGroupType,
			"id": r.AccountGroupId,
		    "rstart": r.PeriodStart,
			"rend": r.PeriodEnd})

	if err != nil {
		log.Printf("Error getting TxnGroupSummary: %v", err);
		return nil
	}

	return &summary
}

// Return a count and sum for the set of transactions grouped by Indusrty Classification
func (r *ReportingApi) TxnGroupSummary(txnGroupId int64, txnType string) *TxnSummary {

	expenses := TxnSummary{} 
	deposits := TxnSummary{} // Sometimes there are deposits for the txnGroup (refunds, etc)

	err := r.dbmap.SelectOne(&expenses, "SELECT count(*) as count, coalesce(sum(amount), 0) as sum FROM txns WHERE txn_group_id = :id AND txn_type = :type AND occurred_at BETWEEN :rstart AND :rend",
		map[string]interface{} { 
			"type": txnType,
			"id": txnGroupId,
		    "rstart": r.PeriodStart,
			"rend": r.PeriodEnd})

	if err != nil {
		log.Printf("Error getting TxnGroupSummary: %v", err);
		return nil
	}

	err = r.dbmap.SelectOne(&deposits, "SELECT count(*) as count, coalesce(sum(amount), 0) as sum FROM txns WHERE txn_group_id = :id AND txn_type != :type AND occurred_at BETWEEN :rstart AND :rend",
		map[string]interface{} { 
			"type": txnType,
			"id": txnGroupId,
		    "rstart": r.PeriodStart,
			"rend": r.PeriodEnd})

	if err != nil {
		log.Printf("Error getting TxnGroupSummary: %v", err);
		return nil
	}
 
	expenses.Count = expenses.Count - deposits.Count
	expenses.Sum = expenses.Sum - deposits.Sum

	return &expenses
}


// Return the set of txn groups for the FI or for the AccountGroupId
func (r *ReportingApi) TxnGroups(group_type string) []TxnGroup {

	//Ideally, this should return only groups that have txns in the date range (or maybe two entry points, or a param)
	var groups []TxnGroup

	var _, err = r.dbmap.Select(&groups, "SELECT * FROM txn_groups WHERE group_type = :type AND account_group_id = :id",
		map[string]interface{} { 
			"type": group_type,
			"id": r.AccountGroupId})

	if err != nil {
		log.Printf("Error getting TxnGroups: %v", err);
		return nil
	}

	return groups
}

func (r *ReportingApi) StartDate() string {
	return r.PeriodStart.Format("January 2, 2006")
}

func (r *ReportingApi) EndDate() string {
	return r.PeriodEnd.Format("January 2, 2006")
}

func (r *ReportingApi) Api() *ReportingApi {
	return r
}

func (r *ReportingApi) SampleMethod(param string) string {
	return " [Sample] " + param
}


func currency(amount int) string {
	dollars := amount / 100
	cents := amount % 100

	return fmt.Sprintf("$%d.%02d", dollars, cents)
}

