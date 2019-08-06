CREATE TABLE grants (
  id BIGINT,
  subject BIGINT,
  azn BIGINT,
  target BIGINT,
  CONSTRAINT grants_key PRIMARY KEY ( id ),
  CONSTRAINT grants_subject_refs_users FOREIGN KEY ( subject ) REFERENCES users ( id ),
  CONSTRAINT grants_azn_refs_azns FOREIGN KEY ( azn ) REFERENCES azns ( id ),
  CONSTRAINT grants_target_refs_entities FOREIGN KEY ( target ) REFERENCES entities ( id )
);

CREATE INDEX azns_index_subject ON grants ( subject );
CREATE INDEX azns_index_target ON grants ( target );
