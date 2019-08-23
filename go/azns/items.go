package azns

import (
  "context"

  "github.com/go-pg/pg"
  "github.com/go-pg/pg/orm"
  "github.com/go-pg/pg/urlvalues"

  "github.com/Liquid-Labs/lc-authentication-api/go/auth"
  . "github.com/Liquid-Labs/terror/go/terror"
)

type PageRequest struct {
  Page         int
  ItemsPerPage int
}
func newPager (pr PageRequest) urlvalues.Pager {
  pager := urlvalues.Pager{ Limit: pr.ItemsPerPage, MaxLimit: 100 }
  pager.SetPage(pr.Page)

  return pager
}

func prepQuery(q *orm.Query, pr PageRequest, ctx context.Context) {
  pager := newPager(pr)

  q.Context(ctx).
    Apply(pager.Pagination)
}

// ListItems retrieves the set of items selected by the base query to which the user has the necessary access rights according to the access route selected.
// The base query may be as simble as:
//
// list = make([]*FinalClass)
// query = db.Model(list)
//
// or it may include additional filter clauses.
func ListOwnedItems(q *orm.Query, pageRequest PageRequest, ctx context.Context) (int, Terror) {
  _, authID, err := auth.CheckAuthentication(ctx)
  if err != nil { return 0, err }

  prepQuery(q, pageRequest, ctx)

  q.Join(`JOIN users AS "owner" ON "owner".auth_id=?`, authID).
    Where(`owner_id=owner.id`)

  if count, err := q.SelectAndCount(); err != nil {
    return 0, ServerError(`Problem retrieving item list.`, err)
  } else {
    return count, nil
  }
}

func ListSharedItemsQuery(model interface{}, db orm.DB, pageRequest PageRequest, ctx context.Context) (int, Terror) {
  _, authID, terr := auth.CheckAuthentication(ctx)
  if terr != nil { return 0, terr }

  tm := db.Model(model).GetModel().Table()
  selectR := string(tm.FullNameForSelects)
  alias := string(tm.Alias)

  var count int
  _, err := db.Query(pg.Scan(&count),
    `WITH RECURSIVE
      subjs(id) AS (
        SELECT id FROM users WHERE auth_id=?
          UNION
        SELECT container_member.container_id AS id
          FROM container_members AS container_member
            JOIN subjs AS subj ON container_member.member=subj.id
      ),
      targs(id) AS (
        SELECT target FROM grants as "grant" JOIN subjs AS "subj" ON "grant".subject="subj".id WHERE "grant".azn=?
           UNION
         SELECT container_member.member AS id
           FROM container_members AS container_member
             JOIN targs AS targ ON container_member.container_id=targ.id
      )
      SELECT COUNT(*) FROM ? AS ? JOIN targs AS target ON target.id=?.id`,
    authID, AznBasicRead.ID, pg.F(selectR), pg.Q(alias), pg.Q(alias))
  if err != nil {
    return 0, ServerError(`Error selecting shared items.`, err)
  }

  _, err = db.Query(model,
    `WITH RECURSIVE
      subjs(id) AS (
        SELECT id FROM users WHERE auth_id=?
          UNION
        SELECT container_member.container_id AS id
          FROM container_members AS container_member
            JOIN subjs AS subj ON container_member.member=subj.id
      ),
      targs(id) AS (
        SELECT target FROM grants as "grant" JOIN subjs AS "subj" ON "grant".subject="subj".id WHERE "grant".azn=?
           UNION
         SELECT container_member.member AS id
           FROM container_members AS container_member
             JOIN targs AS targ ON container_member.container_id=targ.id
      )
      SELECT * FROM ? AS ? JOIN targs AS target ON target.id=?.id ORDER BY ?.name LIMIT ? OFFSET ?`,
    authID, AznBasicRead.ID, pg.F(selectR), pg.Q(alias), pg.Q(alias), pg.Q(alias), pageRequest.ItemsPerPage, pageRequest.ItemsPerPage * pageRequest.Page)
  if err != nil {
    return 0, ServerError(`Error selecting shared items.`, err)
  }

  return count, nil
}
