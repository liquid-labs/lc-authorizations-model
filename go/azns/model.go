package azns

import (
  "log"

  . "github.com/Liquid-Labs/lc-entities-model/go/entities"
)

type AznName string

type Authorization struct {
  ID   int
  Name AznName
}

func NewAuthorizatio(name AznName) *Authorization {
  return &Authorization{0, name}
}

func (a *Authorization) GetID() int { return a.ID }

func (a *Authorization) GetName() AznName { return a.Name }

// Basic authorizations.
// * 'create' is resource specific.
// * 'read-sensitive' is a placeholder for now.
// * We differentiate between 'archive' (common) and 'delete' (rare).
var (
  AznBasicRead = Authorization{1, AznName(`/entities/read`)}
  AznBasicReadSensitive = Authorization{2, AznName(`/entities/read-sensitive`)}
  AznBasicUpdate = Authorization{3, AznName(`/entities/update`)}
  AznBasicArchive = Authorization{4, AznName(`/entities/archive`)}
  AznBasicDelete = Authorization{5, AznName(`/entities/delete`)}
  AznBasicGrant = Authorization{6, AznName(`/entities/grant`)}
)

type JsonB *map[string]interface{}

type Grant struct {
  ID      int64   `sql:",pk"`
  Subject EID
  AznName AznName `sql:"-"`
  Azn     int
  Target  EID
  Cookie  JsonB
}

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

type CapResponse struct {
  Granted     bool
  Cookie      JsonB
  ByOwnership bool
  ByGrant     *Grant
}

func (r *CapResponse) IsGranted() bool { return r.Granted }

func (r *CapResponse) GetCookie() JsonB { return r.Cookie }

func (r *CapResponse) IsByOwnership() bool { return r.ByOwnership }

func (r *CapResponse) GetGrant() *Grant { return r.ByGrant }

var NoSuchCapRespose = &CapResponse{false, nil, false, nil}
