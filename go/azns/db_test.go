package azns_test

import (
  "log"
  "math/rand"
  "os"
  "testing"
  "time"

  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/require"
  "github.com/stretchr/testify/suite"

  . "github.com/Liquid-Labs/lc-entities-model/go/entities"
  "github.com/Liquid-Labs/lc-rdb-service/go/rdb"
  "github.com/Liquid-Labs/terror/go/terror"
  . "github.com/Liquid-Labs/lc-users-model/go/users"

  /* pkg2test*/ "github.com/Liquid-Labs/lc-authorizations-model/go/azns"
)

type TestUser struct {
  User
}
func (tu *TestUser) GetResourceName() ResourceName {
  return ResourceName(`testusers`)
}

type TestThing struct {
  Entity
}
func (tt *TestThing) GetResourceName() ResourceName {
  return ResourceName(`things`)
}

const (
  name = `John Doe`
  desc = `desc`
  legalID = `555-55-5555`
  legalIDType = `SSN`
  active = true
)

func init() {
  terror.EchoErrorLog()
  rand.Seed(time.Now().UnixNano())
}

const runes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_./"
const aznLength = 16

func randStringBytes() string {
    b := make([]byte, aznLength)
    for i := range b {
        b[i] = runes[rand.Int63() % int64(len(runes))]
    }
    return string(b)
}

type GrantIntegrationSuite struct {
  suite.Suite
  User1  *User
  Thing1 *Entity
}
func (s *GrantIntegrationSuite) SetupSuite() {
  db := rdb.Connect()

  authID := randStringBytes()
  s.User1 = NewUser(&TestUser{}, name, desc, authID, legalID, legalIDType, active)
  require.NoError(s.T(), s.User1.Create(db))
  // log.Printf("User1: %s", s.User1.GetID())

  s.Thing1 = NewEntity(&TestThing{}, name, desc, s.User1.GetID(), false)
  require.NoError(s.T(), db.Insert(s.Thing1))
  // log.Printf("Thing1: %s", s.Thing1.GetID())
}
/*func (s *GrantIntegrationSuite) SetupTest() {

}*/
func TestGrantIntegrationSuite(t *testing.T) {
  if os.Getenv(`SKIP_INTEGRATION`) == `true` {
    t.Skip()
  } else {
    suite.Run(t, new(GrantIntegrationSuite))
  }
}

func (s *GrantIntegrationSuite) TestCapabilityByOwnership() {
  CapResponse, err := azns.CheckCapability(s.User1.GetID(), azns.AznBasicUpdate.ID, s.Thing1.GetID(), rdb.Connect())
  require.NoError(s.T(), err)
  assert.Equal(s.T(), true, CapResponse.IsGranted())
  assert.Equal(s.T(), azns.JsonB(nil), CapResponse.GetCookie())
  assert.Equal(s.T(), true, CapResponse.IsByOwnership())
  assert.Equal(s.T(), (*azns.Grant)(nil), CapResponse.GetGrant())
}
