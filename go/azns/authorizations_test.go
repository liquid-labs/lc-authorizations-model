package azns_test

import (
  "testing"

  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/require"

  "github.com/Liquid-Labs/lc-rdb-service/go/rdb"
  "github.com/Liquid-Labs/strkit/go/strkit"

  /* pkg2test*/ "github.com/Liquid-Labs/lc-authorizations-model/go/azns"
)

func TestAuthorizaionsCreate(t *testing.T) {
  aznName := azns.AznName(`/tests/` + strkit.RandString(strkit.Letters, 12))
  a := azns.NewAuthorization(aznName)
  require.NoError(t, a.CreateRaw(rdb.Connect()))
  assert.Less(t, 999, a.GetID())
}
