package azns

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
  ByGrant     bool
}

func (r *CapResponse) IsGranted() bool { return r.Granted }

func (r *CapResponse) GetCookie() JsonB { return r.Cookie }

func (r *CapResponse) IsByOwnership() bool { return r.ByOwnership }

func (r *CapResponse) IsByGrant() bool { return r.ByGrant }

var NoSuchCapRespose = &CapResponse{false, nil, false, false}
