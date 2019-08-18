package azns

import (
  "log"

  "github.com/go-pg/pg/orm"

  . "github.com/Liquid-Labs/lc-entities-model/go/entities"
  . "github.com/Liquid-Labs/terror/go/terror"
)

// TODO: can we put this someplace else?
type JsonB *map[string]interface{}

// ## Grant model

// A Grant is a contingent capability given to a User to act over a target or targets for a specific purpose. E.g., the owner of a document may grant another use the rights to update that document.
// A grant may contain a 'cookie'. Cookies parameters are evaluated at the service rather than the authorization level and may result in an otherwise authorized action be de-authorized. E.g., a cookie may impose quota or time restrictions on an action.
type Grant struct {
  ID      int64   `sql:",pk"`
  Subject EID
  AznName AznName `sql:"-"`
  Azn     int
  Target  EID
  Cookie  JsonB
}

// TODO: if we get a good UC, we could use cookies to implement negative grants. The lower-level grant will be found first, and the special 'NEGATIVE' or 'RESCEND' cookie checked by the base library.

func NewGrant(subject EID, aznRef interface{}, target EID, cookie JsonB) *Grant {
  switch t := aznRef.(type) {
  case int:
    return &Grant{0, subject, ``, aznRef.(int), target, cookie}
  case string:
    return &Grant{0, subject, AznName(aznRef.(string)), 0, target, cookie}
  case AznName:
    return &Grant{0, subject, aznRef.(AznName), 0, target, cookie}
  default:
    log.Panicf(`Unknown type '%s' for 'azn reference'.`, t)
    return nil
  }
}

func (g *Grant) GetID() int64 { return g.ID }

func (g *Grant) GetSubject() EID { return g.Subject }

func (g *Grant) GetAznName() AznName { return g.AznName }

func (g *Grant) GetAzn() int { return g.Azn }

func (g *Grant) GetTarget() EID { return g.Target }

func (g *Grant) GetCookie() JsonB { return g.Cookie }
func (g *Grant) SetCookie(c JsonB) { g.Cookie = c }

// ## DB methods

func (g *Grant) CreateRaw(db orm.DB) Terror {
  if err := db.Insert(g); err != nil {
    return ServerError(`Problem creating grant.`, err)
  }
  return nil
}

func (g *Grant) UpdateRaw(db orm.DB) Terror {
  // TODO: not quite true; we should allow cookies to be updated. This will require coordination with the schema.
  return MethodNotAllowedError(`Grants cannot be updated, only revoked (deleted).`)
}

func (g *Grant) ArchiveRaw(db orm.DB) Terror {
  return MethodNotAllowedError(`Grants cannot be ardived, only deleted.`)
}

func (g *Grant) DeleteRaw(db orm.DB) Terror {
  if err := db.Delete(g); err != nil {
    return ServerError(`Problem deleting grant.`, err)
  }
  return nil
}
