# Database Setup

In this chapter, you will setup the database for use with Cashbook Toolkit.

### Download and install PostgreSQL

PostgreSQL is currently the only supported DBMS.  It can be obtained from

http://www.postgresql.org/

On the Mac, we recommend using

http://postgresapp.com/

### Support for multiple environments

The software fully supports working with multiple database environments. This is normally used for development and production.  The default environment is 'development'.  This can be overridden by a command line argument that explicitly sets the environment to use.

Create at least a development database.  I recommend a convention of cashbook_{envname} where envname is 'development' or 'production' (or 'test', 'staging', etc)

Please refer to http://www.postgresql.org/docs/9.1/static/app-createdb.html for information on creating databases.

### The database configuration file

The database configuration settings live in a yaml file located in **conf/dbconf.yml**.  There is a sample database configuration file shipped in the **conf/** directory.

Copy this file to **conf/dbconf.yml** and edit to reflect the correct database credentials and database names you used when setting up the database.

### Loading and updating the schema

Support for schema management comes built-in, so whether it's loading the initial schema or applying updates when the software is updated, the process is the same:

    ./cashbook migrate:up

It also supports downgrading:

    ./cashbook migrate:down

Or you can see which migrations have already been applied:

    ./cashbook migrate:status

