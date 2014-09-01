# Authoring Reports

### Templates

Cashbook Toolkit is written in Google Go (#golang) and uses the standard template package included with the standard library.  Somewhat good, but terse documentation on them is included with the go documents at:

http://golang.org/pkg/text/template/

A good place to start is by looking at the sample report:

./samplereport.html

### Generating PDF

The toolkit does not generate PDF directly.  However, if you format your reports as HTMl, you can use:

http://wkhtmltopdf.org/

to convert the generated HTML to PDF.

## API Reference

### Objects Returned by API Calls

There are two types of objects returned from the method calls on the Api object:

* Summary
* TxnGroup

#### Summary

Has just two fields:

* Count
* Sum

#### TxnGroup

Has the following fields:

* Label
* GroupType
* Description (Overrides label if non-nil) - user set.
* Classification string

### Available Methods

Please refer to the samplereport.html for a complete example of using these.

#### Api

The entry point into the api is accessible using this method.  Due to the way *#golang* templating works, you will need to store the return value of this method in a variable for later use.

*example*

    {{ $api := .Api }}

#### StartDate

Outputs the starting date for the reports date range.

*example*

     {{.StartDate}}

#### EndDate

Outputs the ending date for the reports date range.

*example*

     {{.EndDate}}

#### Remaining deposits withdrawal

Simply subtracts withdrawls from deposits.

*example*

#### TxnTypeSummary txnType

Return a Summary containing the count and sum for the set of transactions grouped by txnType

* txnType is one of: W or D

*example*

    {{$deposits := $api.TxnTypeSummary "D"}}
    {{$deposits.Count}}
    {{$deposits.Sum | currency}}

    Withdrawals work exactly the same.

#### TxnGroupTypeSummary txnGroupType

Return a Summary containing the count and sum for the set of transactions grouped by txnGroupType

* txnGroupType must be one of the groupTypes defined in your sef of matchers (but is typically one of Retail, Loans, Bill Payments, etc).

*example*

    {{ $retail := $api.TxnGroupTypeSummary "Retail"}}
    <h2>
        Retail Spending {{$retail.Sum | currency}} ({{$retail.Count}})
    </h2>

#### TxnGroups txnGroupType

Return an array of TxnGroup either for the financial institution or for the AccountGroupId that are grouped under txnGroupType

* txnGroupType must be one of the groupTypes defined in your sef of matchers (but is typically one of Retail, Loans, Bill Payments, etc).

*example*

     {{range .TxnGroups "Retail"}}

     Note: Range sets a new context for . (it becomes the current TxnGroup). This is why we have to store .Api into $api as shown earlier.


#### TxnGroupSummary txnGroupId txnType

Return a Summary containing the count and sum for the set of transactions grouped by txnGroupId and txnType

* txnGroupId is the database id of a transaction group
* txnType is one of W or D

*example*

    {{ $summary := $api.TxnGroupSummary .Id "W"}}
    <td style="text-align:left;">{{ $summary.Count }} </td>
    <td style="text-align:right;">{{ $summary.Sum | currency}}</td>

    Note: .Id is being called in the txnGroup in the range from the TxnGroups call.

### Helper functions

If you pipe (|) the output (just like in the bourne shell) of a method call (or just stored data) to these functions, they will format the data for output.

#### currency

Renders an amount as a nice currency string.

*example*

    {{ $summary.Sum | currency }}

#### Capitalize

Capitalizes the first character.

*example*

    {{.Label | capitalize}}


