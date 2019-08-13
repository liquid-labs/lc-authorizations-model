package caps

type AccessRoute int
const (
  AccessPublic    AccessRoute = 0 // default
  AccessRoot      AccessRoute = 1
  AccessGrant     AccessRoute = 2
  // AccessAny       AccessRoute = 3 -- Not sure there's a UC for this.
)
