package azns

import (
  "context"

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

func checkAuthentication(ctx context.Context) (string, Terror) {
  authenticator := auth.GetAuthenticator(ctx)
  if !authenticator.IsRequestAuthenticated() {
    return ``, UnauthenticatedError(`Non-Authenticated user cannot requested 'owned' items.`)
  }
  aznID := authenticator.GetAznID()
  if aznID == `` {
    return ``, ServerError(`Missing authorization ID for authenticated user.`, nil)
  }
  return aznID, nil
}

// ListItems retrieves the set of items selected by the base query to which the user has the necessary access rights according to the access route selected.
// The base query may be as simble as:
//
// list = make([]*FinalClass)
// query = db.Model(list)
//
// or it may include additional filter clauses.
func ListOwnedItems(q *orm.Query, pageRequest PageRequest, ctx context.Context) (int, Terror) {
  if q == nil {
    return 0, BadRequestError(`Request does not resolve to a base query.`)
  }

  aznID, err := checkAuthentication(ctx)
  if err != nil { return 0, err }
  pager := newPager(pageRequest)

  q.Context(ctx).
    Apply(pager.Pagination).
    Join(`JOIN users AS "owner" ON "owner".auth_id=?`, aznID).
    Where(`owner_id=owner.id`)

  if count, err := q.SelectAndCount(); err != nil {
    return 0, ServerError(`Problem retrieving item list.`, err)
  } else {
    return count, nil
  }
}
/*
import (
  "fmt"

  "github.com/Liquid-Labs/lc-authorizations-model/go/azns"
  "github.com/Liquid-Labs/lc-authorizations-model/go/caps"
  "github.com/Liquid-Labs/lc-entities-model/go/entities"
  "github.com/Liquid-Labs/lc-rdb-service/go/rdb"
  "github.com/Liquid-Labs/terror/go/terror"
)

// CreateItem will check user permissions via the indicated accessRoute and
// create a new record of the provided Entity (sub-type) where authorized.
func CreateItem(item entities.Identifiable, accessRoute caps.AccessRoute, ctx context.Context) terror.Terror {
  if item == nil {
    return terror.BadRequestError(`Entity for creation cannot be nil.`, nil)
  }

  authenticator := auth.GetAuthenticator(ctx)

  if authResponse, restErr := caps.CheckAuthorization(authenticator.GetAznID(), `/create/` + item.GetResourceName(), nil); restErr != nil {
    return restErr
  } else if !authResponse.Granted {
    // TODO: get helper to get us the name... method reciever for Entity?
    return terror.AuthorizationError(`User not authorized to create resource.`, nil)
  } else {
    if err := rdb.Connect().Insert(item); err != nil {
      return terror.ServerError(`Problem creating resource.`, err)
    } else {
      return nil
    }
  }
}

// GetItem will attempt to retrieve an Entity by either the public or internal
// ID. Which to use is determined by the 'id' type, which must be either a
// string (for public ID) or int64 (for internal ID). The base query is
// typically just 'db.Model(yourStruct)', where the struct used must embed
// Entity. GetItem adds hte necessary authorization checks to the provided
// base query.
func GetItem(id interface{}, baseQuery *orm.Query, accessRoute caps.AccessRoute, ctx context.Context) terror.Terror {
  if baseQuery == nil {
    return terror.BadRequestError(`Request does not resolve to a base query. Contact customer support if you believe this is a bug.`, nil)
  }

  query := baseQuery.Context(ctx)
  switch id.(type) {
  case string: // or entities.PublicID?
    query = baseQuery.Where(`e.pub_id=?`, id)
  case int64: // or entities.InternalID?
    query = baseQuery.Where(`e.ID=?`, id)
  default:
    return terror.BadRequestError(fmt.Sprintf(`Invalid identifier '%v' supplied to 'GetItem'.`, id), nil)
  }

  query = caps.AuthorizedModel(query, accessRoute, azns.EntityRead.ID, ctx)

  if err := query.Select(); err != nil {
    // Notice we don't return the ID because it may be a oddly formatted
    // internal ID, which should not be revealed.
    // TODO: we should log the info though.
    return terror.ServerError(`Problem retrieving entity.`, err)
  }
  return nil
}

// UpdateItem updates the entity provided the user has sufficient authorizations
// via the indicated access route.
func UpdateItem(item entities.Identifiable, accessRoute caps.AccessRoute, ctx context.Context) terror.Terror {
  if item == nil {
    return terror.BadRequestError(`No entity provided for update.`, nil)
  }

  query := rdb.Connect().Model(item).Context(ctx)
  query = caps.AuthorizedModel(query, accessRoute, azns.EntityRead.ID, ctx)

  var err error
  if item.GetID() != 0 {
    _, err = query.Update()
  } else if item.GetPublicID() != `` {
    _, err = query.Where(`entities.public_id=?public_id`).Update()
  } else {
    return terror.BadRequestError(`Provided item has no valid identifer.`, nil)
  }
  if err != nil {
    return terror.ServerError(`Problem updating entity.`, err)
  }
  return nil
}

// ArchiveItem performs a soft-delete of the indicated item provided the user
// has sufficient authorizations via the indicated access route.
func ArchiveItem(item *entities.Identifiable, accessRoute caps.AccessRoute, ctx context.Context) terror.Terror {
  if item == nil {
    return terror.BadRequestError(`No entity provided for archival.`, nil)
  }

  query := rdb.Connect().Model(item).Context(ctx)
  query = caps.AuthorizedModel(query, accessRoute, azns.EntityArchive.ID, ctx)

  if _, err := query.Delete(); err != nil {
    return terror.ServerError(`Problem updating entity.`, err)
  }
  return nil
}
*/
