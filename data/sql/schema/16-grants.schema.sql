CREATE TABLE grants (
  id BIGINT,
  subject BIGINT,
  azn BIGINT NOT NULL,
  target BIGINT,
  -- cookies VARCHAR(128)[], -- e.g. : "LIMIT:10" on create, etc.
  CONSTRAINT grants_key PRIMARY KEY ( id ),
  CONSTRAINT grants_subject_refs_users FOREIGN KEY ( subject ) REFERENCES users ( id ),
  CONSTRAINT grants_azn_refs_azns FOREIGN KEY ( azn ) REFERENCES azns ( id ),
  CONSTRAINT grants_target_refs_entities FOREIGN KEY ( target ) REFERENCES entities ( id ),
  CONSTRAINT grants_single_grant UNIQUE ( subject, azn, target )
);

CREATE INDEX azns_index_subject ON grants ( subject );
CREATE INDEX azns_index_target ON grants ( target );
