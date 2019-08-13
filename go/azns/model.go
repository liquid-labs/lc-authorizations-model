package azns

type Authorization struct {
  Name string
  ID   int
}

// These mirror the basic inserts.
var (
  EntityRead = Authorization{`/entities/read`, 0}
  EntityReadSensitive = Authorization{`/entities/read-sensitive`, 1}
  EntityUpdate = Authorization{`/entities/update`, 2}
  EntityArchive = Authorization{`/entities/archive`, 3}
  EntityDelete = Authorization{`/entities/delete`, 4}
  EntityGrant = Authorization{`/entities/grant`, 5}
)
