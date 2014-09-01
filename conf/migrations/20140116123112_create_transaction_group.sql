
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE SEQUENCE txn_group_id_seq;

create table txn_groups(
       id int NOT NULL DEFAULT nextval('txn_group_id_seq'), 

       group_type varchar(255),       -- Retail, Bill Payment, Loan, Cheque, etc

       label varchar(255),            -- Comes from the transaction description
       description varchar(255),      -- Will use label if null - allows member to override label

       classification varchar(255),   -- Transportation, Groceries, Department, Dining, etc

       system_txn_group_id int,       -- FK to a System TxnGroup
       category_id int,               -- FK to categories table - null if system TG

       -- If account_group_id is null, then this transaction_group is a system group
       account_group_id varchar(255),  -- Account groups represent members

       created bigint,                 -- Unix timestamp

       PRIMARY KEY(id),
       UNIQUE (label, group_type, account_group_id)
);

CREATE INDEX ON txn_groups (group_type);
CREATE INDEX ON txn_groups (classification);
CREATE INDEX ON txn_groups (system_txn_group_id);
CREATE INDEX ON txn_groups (account_group_id);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE txn_groups;
