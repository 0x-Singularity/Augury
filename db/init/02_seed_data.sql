-- db/init/02_seed_data.sql
INSERT INTO ioc_query_log (ioc, result_count, user_name)
VALUES ('example.com', 42, 'alice'),
       ('1.2.3.4',     17, 'bob');
