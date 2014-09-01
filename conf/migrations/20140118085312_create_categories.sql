
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
create table categories(
       id int NOT NULL DEFAULT nextval('category_id_seq'),

       label varchar(255), 
       spending_limit bigint,
       limit_interval int,          -- Maybe in days...

       account_group_id varchar(255),

       -- Server unix Timestamp when record is created.  
       created bigint NOT NULL DEFAULT cast(extract(epoch from current_timestamp) as integer),

       PRIMARY KEY(id)
);

CREATE INDEX ON categories (account_group_id);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE categories;
