CREATE TABLE user_groups (
  id UUID,
  CONSTRAINT user_groups_key PRIMARY KEY ( id ),
  CONSTRAINT user_groups_refs_subjects FOREIGN KEY ( id ) REFERENCES subjects ( id )
);

CREATE TABLE user_group_members (
  id     UUID,
  member UUID,
  CONSTRAINT user_group_members_key PRIMARY KEY ( id, member ),
  CONSTRAINT user_group_members_member_refs_user FOREIGN KEY ( member ) REFERENCES users ( id )
);
