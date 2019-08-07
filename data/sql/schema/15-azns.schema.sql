CREATE TABLE azns (
  id BIGINT,
  name VARCHAR(128) NOT NULL UNIQUE,
  CONSTRAINT azns_key PRIMARY KEY ( id )
);

CREATE INDEX azns_index_name ON azns ( name );

-- Notice there is no create; that's a per-type subject-less grant.
INSERT INTO azns (0, '/entities/read');
INSERT INTO azns (1, '/entities/read-sensitive');
INSERT INTO azns (2, '/entities/write');
INSERT INTO azns (3, '/entities/archive');
INSERT INTO azns (4, '/entities/delete');
INSERT INTO azns (5, '/entities/grant');
