package caps

import (
  "context"
  "log"

  "github.com/go-pg/pg/orm"
  "github.com/Liquid-Labs/lc-authentication-api/go/auth"
  "github.com/Liquid-Labs/terror/go/terror"
)

type AuthorizationResponse struct {
  Granted bool
  Cookie  interface{} // could be any JSON derived structure; string, int, float, map, or array.
}

func CheckAuthorization(subject interface{}, action string, target interface{}) (*AuthorizationResponse, terror.Terror) {
  switch subject.(type) {
  case string: // or entities.PublicID? Also could be AznID (which we should type...)
    return nil, nil
  case int64: // or entities.InternalID?
    return nil, nil
  default:
    return nil, nil
  }
}

func resolveAuthorization(authorization interface{}, query *orm.Query) *orm.Query {
  switch authorization.(type) {
  case int:
    return query.Where(`grants.authorization=?`, authorization)
  case int64:
    return query.Where(`grants.authorization=?`, authorization)
  case string:
    return query.
      Join(`JOIN azns ON grants.authorization=azns.id`).
      Where(`azns.name=?`, authorization)
  default:
    return query
  }
}

func AuthorizedModel(baseQuery *orm.Query, accessRoute /*azn.*/AccessRoute, authorization interface{}, ctx context.Context) *orm.Query {
  if accessRoute == AccessPublic {
    return authorizedPublicModel(baseQuery, authorization)
  } else {
    authenticator := auth.GetAuthenticator(ctx)
    if !authenticator.IsRequestAuthenticated() {
      return authorizedPublicModel(baseQuery, authorization)
    }
    // else, the request is authenticated
    if accessRoute == AccessRoot {
      if authenticator.HasAllClaims(`root`) {
        return baseQuery
      } else {
        log.Panicf(`Cannot make 'root' request as non-root user.`)
        return nil
      }
    } else if accessRoute == AccessGrant {
      return authorizedGrantModel(baseQuery, authorization, ctx)
    } else {
      log.Panicf(`Unmatched 'access route' value: '%d'`, accessRoute)
      return nil
    }
  }
}

func authorizedPublicModel(q *orm.Query, authorization interface{}) *orm.Query {
  q = q.
    Join("JOIN grants ON grants.target=e.id").
    Where("read_public=TRUE OR grants.subject IS NULL")
  return resolveAuthorization(authorization, q)
}

func authorizedGrantModel(q *orm.Query, authorization interface{}, ctx context.Context) *orm.Query {
  authID := auth.GetAuthenticator(ctx).GetAznID()

  var f orm.Formatter
  recursiveCTE := string(f.FormatQuery(nil,
    `WITH RECURSIVE group(id) AS (
      SELECT agm.group AS id FROM azn_group_members agm JOIN users u ON agm.member=u.id WHERE u.auth_id=?
    UNION
      SELECT agm.group AS id FROM azn_group_members agm WHERE agm.member=group.id`,
    authID))

  query := q.
    WrapWith(recursiveCTE).
    Join(`JOIN group`).
    Join(`JOIN container`).
    Join(`JOIN users u ON u.auth_id=?`, authID).
    Join(`JOIN grants ON (grants.subject IS NULL OR grants.subject=u.id OR grants.subject=group.id) AND (grants.target=e.id OR e.containers @> ARRAY[azn_grant.target])`).
    Where(`containers @> ARRAY['s']::varchar[]`).
    Group(`entities.id`)

  return resolveAuthorization(authorization, query)
}
