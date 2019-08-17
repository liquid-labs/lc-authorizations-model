package azns

type AznName string

type Authorization struct {
  Name AznName
  ID   int
}

// Basic authorizations.
// * 'create' is resource specific.
// * 'read-sensitive' is a placeholder for now.
// * We differentiate between 'archive' (common) and 'delete' (rare).
var (
  EntityRead = Authorization{`/entities/read`, 1}
  EntityReadSensitive = Authorization{`/entities/read-sensitive`, 2}
  EntityUpdate = Authorization{`/entities/update`, 3}
  EntityArchive = Authorization{`/entities/archive`, 4}
  EntityDelete = Authorization{`/entities/delete`, 5}
  EntityGrant = Authorization{`/entities/grant`, 6}
)
