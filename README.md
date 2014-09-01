
Cashbook Core is part of http://cashbooktoolkit.com/

## About

Cashbook is a toolkit for reading, analyzing and reporting of financial
transactions.  It is designed to be used as part of a Credit Union's
PFM (Personal Finance Manager) strategy.

## Initial Setup (One time)

  1. Install PostgreSQL (or have an existing DB server available)
  2. Check the usename and password for the database in
     conf/dbconf.yml and adjust as needed.
  3. Create the database 'cashbook_devel' See:
     http://www.postgresql.org/docs/9.1/static/manage-ag-createdb.html
  4. Load the schema (Note this is also the command to update the schema):
     ./cashbook migrate:up
     
## Quick Start

  1. Perform the initial setup
  2. Load some data into the Cashbook database 
     TODO Add example
  3. Run a report
     TODO Add example

## Documentation

The software is as self documenting as possible.  If you run the cashbook
tool with no arguments it will print it's usage and exit.

Additional (but somewhat incomplete at this time) information can be found in guide folder.

The guide is written so that it can be processed by http://www.gitbook.io/ into a nice
html version, but the Markdown files can also be read as is.

## Contact

Bug, comments or questions? 

Email: vince@sourdoughlabs.com
Web:   http://cashbooktoolkit.com/

## Contributing

1. Fork it ( https://github.com/sourdoughlabs/cashbook-core/fork )
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create a new Pull Request

## TODO

* Deal with Database Nullable columns (currently uses the zero value as null)
* Threading/goroutines
* Documentation

## License/Terms of Use

   Cashbook is Copyright (C) 2014  Sourdough Labs Research and Development Corp.

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.

