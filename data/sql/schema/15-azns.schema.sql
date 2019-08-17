CREATE TABLE azns (
  id   INT,
  name VARCHAR(128) NOT NULL UNIQUE,
  CONSTRAINT azns_key PRIMARY KEY ( id )
);

CREATE INDEX azns_index_name ON azns ( name );

-- Notice there is no create; that's a per-type subject-less grant.
INSERT INTO azns VALUES (1, '/entities/read');
INSERT INTO azns VALUES (2, '/entities/read-sensitive');
INSERT INTO azns VALUES (3, '/entities/update');
INSERT INTO azns VALUES (4, '/entities/archive');
INSERT INTO azns VALUES (5, '/entities/delete');
INSERT INTO azns VALUES (6, '/entities/grant');
