-- db/init/01_create_tables.sql   (PostgreSQL syntax)

CREATE TABLE ioc_query_log (
    id           SERIAL PRIMARY KEY,
    ioc          VARCHAR(255) NOT NULL,
    last_lookup  TIMESTAMPTZ  DEFAULT now(),
    result_count INT          NOT NULL,
    user_name    VARCHAR(255)
);
