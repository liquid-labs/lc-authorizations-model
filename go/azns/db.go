package azns

import (
  "log"
  "strconv"

  "github.com/go-pg/pg"
  "github.com/go-pg/pg/orm"

  . "github.com/Liquid-Labs/lc-entities-model/go/entities"
  . "github.com/Liquid-Labs/terror/go/terror"
)

var ownershipSelect =
  `SELECT
      ` + strconv.Itoa(AznRouteOwner.ID) +` AS route_id,
      '` + AznRoutePublic.Name + `' AS route_name,
      NULL AS cookie
    FROM entities AS entity
    WHERE entity.id=? AND entity.owner_id=?`

var grantSelect =
  `WITH RECURSIVE
    targs(id) AS (
        SELECT container_member.id AS id
          FROM container_members AS container_member
          WHERE container_member.member=?
      UNION
        SELECT container_member.id AS id
          FROM container_members AS container_member
            JOIN targs AS targ ON container_member.member=targ.id
    ),
    subjs(id) AS (
        SELECT user_group.id AS id
          FROM user_groups AS user_group
            JOIN container_members AS container_member ON user_group.id=container_member.id
          WHERE container_member.member=?
      UNION
        SELECT container_member.id AS id
          FROM container_members AS container_member
            JOIN subjs AS subj ON container_member.member=subj.id
      )
    SELECT
        ` + strconv.Itoa(AznRouteGrant.ID) + ` AS route_id,
        '` + AznRouteGrant.Name + `' AS route_name,
        "grant".cookie AS cookie
      FROM grants AS "grant"
        JOIN subjs AS subj ON "grant".subject=subj.id
        JOIN targs AS targ ON "grant".target=targ.id
      WHERE "grant".azn=?
      LIMIT 1`

var authSelect =
  `SELECT route_id, route_name, cookie FROM (
    (` + ownershipSelect + `)
    UNION
    (` + grantSelect + `)
  ) AS tmp LIMIT 1`

type capResults struct {
  RouteID   int
  RouteName string
  Cookie    JsonB
}

func CheckCapability(subject EID, aznRef interface{}, target EID, db orm.DB) (*CapResponse, Terror) {
  var err error
  capResult := &capResults{}
  switch t := aznRef.(type) {
  case int:
    _, err = db.QueryOne(capResult, authSelect, target, subject, target, subject, aznRef.(int))
  default:
    log.Panicf("Invalid azn reference type '%s'.", t)
  }

  if err != nil && err != pg.ErrNoRows {
    return nil, ServerError(`There was a problem checking for capability.`, err)
  } else if err == pg.ErrNoRows {
    return NoSuchCapRespose, nil
  } // else process response.

  return &CapResponse{
    true,
    capResult.Cookie,
    capResult.RouteID == AznRouteOwner.ID,
    nil,
  }, nil
}

/*
func (g *Grant) Create(db orm.DB) terror.Terror {
  // Check grant is properly formed
  if g.ID != 0 {
    return terror.BadRequestError(`Cannot set ID when creating grant.`, nil)
  }

  if g.AznID == 0 && g.AznName == `` {
    return terror.BadRequestError(`Grant must specify authorization.`, nil)
  } else if g.AznID != 0 && g.AznName != `` {
    return terror.BadRequestError(`Grant authorization overconstrained; designate either name or ID, not both`, nil)
  }

  var err error
  if g.AznName != `` {
    _, err =
      db.Exec(`INSERT INTO grant (subject, azn, target) SELECT ?, azn.id ? FROM azns as azn WHERE azn.name=?`,
        g.SubjectID,
        g.TargetID,
        g.AznName)
  } else {
    err = db.Insert(g)
  }
  if err != nil {
    return terror.BadRequestError(`Could not create grant. Check parameters.`, err)
  } else {
    return nil
  }
}
*/
