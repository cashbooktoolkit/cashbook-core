
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE SEQUENCE txn_id_seq;

create table txns(
       id bigint NOT NULL DEFAULT nextval('txn_id_seq'),

       txn_type varchar(2),           -- W || D

       amount int,                    --  pennies
       description varchar(255),
       occurred_at timestamp,         -- Int64 is also possibility and store unixtime, but Rails

       txn_host_type varchar(50),     -- Code/type straight from core system
       trace_number varchar(255),     -- From the host
       combination_key varchar(255),  -- from the host system, TODO Talk to Rob what's this for again

       category_id int,               -- FK to categories table - null if system TG

       system_txn_group_id int,       -- FK to system TxnGroup
       txn_group_id int,              -- FK to a TxnGroup
       txn_group_type varchar(255),   -- Retail, Bill Payment, Loan, Cheque, etc

       classification varchar(255),   -- Transportation, Groceries, Department, Dining, etc

       account_group_id varchar(255), -- Account groups represent members

       created bigint,                -- Server unix Timestamp when record is created.

       PRIMARY KEY(id),
       UNIQUE (trace_number, combination_key),       -- TODO Just this or this and combo_key?
       CHECK (txn_type in ('W', 'D'))
);

CREATE INDEX ON txns (txn_type);
CREATE INDEX ON txns (occurred_at);
CREATE INDEX ON txns (txn_host_type);
CREATE INDEX ON txns (category_id);
CREATE INDEX ON txns (system_txn_group_id);
CREATE INDEX ON txns (txn_group_id);
CREATE INDEX ON txns (txn_group_type);
CREATE INDEX ON txns (classification);
CREATE INDEX ON txns (account_group_id);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE txns;
