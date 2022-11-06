-- migrate:up
create table users (
  id varchar(255),
  name varchar(255)
);

-- migrate:down
drop table if exists users
