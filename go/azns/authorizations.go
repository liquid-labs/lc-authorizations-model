package azns

import (
  "github.com/go-pg/pg/orm"

  . "github.com/Liquid-Labs/terror/go/terror"
)

type AznName string

type Authorization struct {
  tableName struct{} `sql:"azns,alias:azn"`
  ID        int      `sql:",pk"`
  Name      AznName
}

func NewAuthorization(name AznName) *Authorization {
  return &Authorization{ID: 0, Name: name}
}

func (a *Authorization) GetID() int { return a.ID }

func (a *Authorization) GetName() AznName { return a.Name }

// Basic authorizations.
// * 'create' is resource specific.
// * 'read-sensitive' is a placeholder for now.
// * We differentiate between 'archive' (common) and 'delete' (rare).
var (
  AznBasicRead = Authorization{ID: 1, Name: AznName(`/entities/read`)}
  AznBasicReadSensitive = Authorization{
    ID: 2,
    Name: AznName(`/entities/read-sensitive`)}
  AznBasicUpdate = Authorization{ID: 3, Name: AznName(`/entities/update`)}
  AznBasicArchive = Authorization{ID: 4, Name: AznName(`/entities/archive`)}
  AznBasicDelete = Authorization{ID: 5, Name: AznName(`/entities/delete`)}
  AznBasicGrant = Authorization{ID: 6, Name: AznName(`/entities/grant`)}
)

func (a *Authorization) CreateRaw(db orm.DB) Terror {
  if err := db.Insert(a); err != nil {
    return ServerError(`Problem creating authorization.`, err)
  }
  return nil
}

type AznRoute struct {
  ID   int
  Name string
}
var (
  AznRoutePublic = AznRoute{1, "public"}
  AznRouteOwner  = AznRoute{2, "ownership"}
  AznRouteGrant  = AznRoute{3, "grant"}
  AznRouteRoot   = AznRoute{4, "root"}
)
