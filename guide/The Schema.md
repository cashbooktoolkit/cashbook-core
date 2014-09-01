# The Schema

Describes the tables, fields and datatypes of the schema and any special notes.

**List of relations**

| Name | Type|
| -- | -- |
| txns                    | table    |
| txn_groups              | table    |
| categories              | table    |
| txn_group_id_seq        | sequence |
| txn_id_seq              | sequence |

** Txns **

|Column| Type | Modifiers |
| -- | -- | -- |
|id                  | bigint                      | not null default nextval('txn_id_seq'::regclass) |
| txn_type            | character varying(2)        |-|
| amount              | integer                     |-|
| description         | character varying(255)      |-|
| occurred_at         | timestamp without time zone |-|
| txn_host_type       | character varying(50)       |-|
| trace_number        | character varying(255)      |-|
| combination_key     | character varying(255)      |-|
| category_id         | integer                     |-|
| system_txn_group_id | integer                     |-|
| txn_group_id        | integer                     |-|
| txn_group_type      | character varying(255)      |-|
| classification      | character varying(255)      |-|
| account_group_id    | character varying(255)      |-|
| created             | bigint                      |-|

Indexes:

    "txns_pkey" PRIMARY KEY, btree (id)
    "txns_trace_number_combination_key_key" UNIQUE CONSTRAINT, btree (trace_number, combination_key)
    "txns_account_group_id_idx" btree (account_group_id)
    "txns_category_id_idx" btree (category_id)
    "txns_classification_idx" btree (classification)
    "txns_occurred_at_idx" btree (occurred_at)
    "txns_system_txn_group_id_idx" btree (system_txn_group_id)
    "txns_txn_group_id_idx" btree (txn_group_id)
    "txns_txn_group_type_idx" btree (txn_group_type)
    "txns_txn_host_type_idx" btree (txn_host_type)
    "txns_txn_type_idx" btree (txn_type)

Check constraints:

    "txns_txn_type_check" CHECK (txn_type::text = ANY (ARRAY['W'::character varying, 'D'::character varying]::text[]))

**txn_groups**

|Column| Type | Modifiers |
| -- | -- | -- |
|id                  | integer                | not null default nextval('txn_group_id_seq'::regclass)|
| group_type          | character varying(255) |-|
| label               | character varying(255) |-|
| description         | character varying(255) |-|
| classification      | character varying(255) |-|
| system_txn_group_id | integer                |-|
| category_id         | integer                |-|
| account_group_id    | character varying(255) |-|
| created             | bigint                 |-|

Indexes:

    "txn_groups_pkey" PRIMARY KEY, btree (id)
    "txn_groups_label_group_type_account_group_id_key" UNIQUE CONSTRAINT, btree (label, group_type, account_group_id)
    "txn_groups_account_group_id_idx" btree (account_group_id)
    "txn_groups_classification_idx" btree (classification)
    "txn_groups_group_type_idx" btree (group_type)
    "txn_groups_system_txn_group_id_idx" btree (system_txn_group_id)

**Categories**

|Column| Type | Modifiers |
| -- | -- | -- |
|id | integer | not null|
|label            | character varying(255) |-|
|spending_limit   | bigint                 |-|
|limit_interval   | integer                |-|
|account_group_id | character varying(255) |-|
|created          | bigint                 |-|

Indexes:

    "categories_pkey" PRIMARY KEY, btree (id)
    "categories_account_group_id_idx" btree (account_group_id)
