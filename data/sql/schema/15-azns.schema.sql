CREATE TABLE azns (
  id BIGINT,
  name VARCHAR(64) NOT NULL UNIQUE,
  CONSTRAINT azns_key PRIMARY KEY ( id )
);

CREATE INDEX azns_index_name ON azns ( name );
