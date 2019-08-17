CREATE TABLE user_groups (
  id UUID,
  CONSTRAINT user_groups_key PRIMARY KEY ( id ),
  CONSTRAINT user_groups_refs_subjects FOREIGN KEY ( id ) REFERENCES subjects ( id ),
  CONSTRAINT user_groups_refs_containers FOREIGN KEY ( id ) REFERENCES containers ( id )
);
