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

  . "github.com/Liquid-Labs/lc-containers-model/go/containers"
  . "github.com/Liquid-Labs/lc-entities-model/go/entities"
  "github.com/Liquid-Labs/lc-rdb-service/go/rdb"
  "github.com/Liquid-Labs/terror/go/terror"
  . "github.com/Liquid-Labs/lc-users-model/go/users"

  /* pkg2test*/ "github.com/Liquid-Labs/lc-authorizations-model/go/azns"
)

const (
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
  User1        *User
  Thing1A      *Entity
  Thing1B      *Entity
  Thing1C      *Entity
  Thing1GroupA *Container
  Thing1GroupC *Container
  User1Group   *azns.UserGroup
  User2        *User
  Thing2       *Entity
}
// SetupSuite sets up the following capabilities:
//
// * User1 owns Thing1, Thing1Group, and User1Group
// * User2 owns Thing2
// * User1 has direct rights over Thing2
func (s *GrantIntegrationSuite) SetupSuite() {
  db := rdb.Connect()

  // setup base objects
  authID1 := randStringBytes()
  s.User1 = NewUser(`users`, `User1`, ``, authID1, legalID, legalIDType, active)
  require.NoError(s.T(), s.User1.CreateRaw(db))
  // log.Printf("User1: %s", s.User1.GetID())
  s.Thing1A = NewEntity(`entities`, `Thing1A`, ``, s.User1.GetID(), false)
  require.NoError(s.T(), CreateEntityRaw(s.Thing1A, db))
  // log.Printf("Thing2: %s", s.Thing2.GetID())
  s.Thing1B = NewEntity(`entities`, `Thing1B`, ``, s.User1.GetID(), false)
  require.NoError(s.T(), CreateEntityRaw(s.Thing1B, db))
  // log.Printf("Thing2B: %s", s.Thing2B.GetID())
  s.Thing1C = NewEntity(`entities`, `Thing1C`, ``, s.User1.GetID(), false)
  require.NoError(s.T(), CreateEntityRaw(s.Thing1C, db))
  // log.Printf("Thing2B: %s", s.Thing2B.GetID())
  s.Thing1GroupA = &Container{
    Entity:Entity{ResourceName: `containers`, Name:`Thing1GroupA`, OwnerID:s.User1.GetID()},
    Members:[]*Entity{s.Thing1A},
  }
  require.NoError(s.T(), s.Thing1GroupA.CreateRaw(db))
  // log.Printf("Thing2: %s", s.Thing2Group.GetID())
  s.Thing1GroupC = &Container{
    Entity:Entity{ResourceName: `containers`, Name:`Thing1GroupC`, OwnerID:s.User1.GetID()},
    Members:[]*Entity{s.Thing1C},
  }
  require.NoError(s.T(), s.Thing1GroupC.CreateRaw(db))
  // log.Printf("Thing2: %s", s.Thing2Group.GetID())

  authID2 := randStringBytes()
  s.User2 = NewUser(`users`, `User2`, ``, authID2, legalID, legalIDType, active)
  require.NoError(s.T(), s.User2.CreateRaw(db))
  // log.Printf("User2: %s", s.User2.GetID())

  s.User1Group = &azns.UserGroup{
    Container:Container{
      Entity:Entity{ResourceName: `containers`, Name:`User1Group`, OwnerID:s.User1.GetID()},
      Members:[]*Entity{&s.User2.Entity},
    },
  }
  require.NoError(s.T(), s.User1Group.CreateRaw(db))
  // log.Printf("Thing2: %s", s.Thing2Group.GetID())

  s.Thing2 = NewEntity(`entities`, `Thing2`, ``, s.User2.GetID(), false)
  require.NoError(s.T(), CreateEntityRaw(s.Thing2, db))
  // log.Printf("Thing1: %s", s.Thing1.GetID())

  grants := []*azns.Grant{
    /* U1 > T2 */ azns.NewGrant(s.User1.GetID(), azns.AznBasicUpdate.ID, s.Thing2.GetID(), nil),
    /* (U2 > T1GA) > T1A */ azns.NewGrant(s.User2.GetID(), azns.AznBasicUpdate.ID, s.Thing1GroupA.GetID(), nil),
    /* U2 > (U1G > T1B) */ azns.NewGrant(s.User1Group.GetID(), azns.AznBasicUpdate.ID, s.Thing1B.GetID(), nil),
    /* U2 > (U1G > T1GC) > T1C */ azns.NewGrant(s.User1Group.GetID(), azns.AznBasicUpdate.ID, s.Thing1GroupC.GetID(), nil),
  }
  for i, g := range grants {
    log.Printf("doing %d\n\n", i)
    require.NoError(s.T(), g.CreateRaw(db))
  }
}

func TestGrantIntegrationSuite(t *testing.T) {
  if os.Getenv(`SKIP_INTEGRATION`) == `true` {
    t.Skip()
  } else {
    suite.Run(t, new(GrantIntegrationSuite))
  }
}

func (s *GrantIntegrationSuite) TestCapabilityByOwnership() {
  CapResponse, err := azns.CheckCapability(s.User1.GetID(), azns.AznBasicUpdate.ID, s.Thing1A.GetID(), rdb.Connect())
  require.NoError(s.T(), err)
  assert.Equal(s.T(), true, CapResponse.IsGranted())
  assert.Equal(s.T(), azns.JsonB(nil), CapResponse.GetCookie())
  assert.Equal(s.T(), true, CapResponse.IsByOwnership())
  assert.Equal(s.T(), false, CapResponse.IsByGrant())
}

func (s *GrantIntegrationSuite) TestCapabilityByDirectGrant() {
  CapResponse, err := azns.CheckCapability(s.User1.GetID(), azns.AznBasicUpdate.ID, s.Thing2.GetID(), rdb.Connect())
  require.NoError(s.T(), err)
  assert.Equal(s.T(), true, CapResponse.IsGranted())
  assert.Equal(s.T(), azns.JsonB(nil), CapResponse.GetCookie())
  assert.Equal(s.T(), false, CapResponse.IsByOwnership())
  assert.Equal(s.T(), true, CapResponse.IsByGrant())
}

func (s *GrantIntegrationSuite) TestCapabilityByIndirectTargetGrant() {
  CapResponse, err := azns.CheckCapability(s.User2.GetID(), azns.AznBasicUpdate.ID, s.Thing1A.GetID(), rdb.Connect())
  require.NoError(s.T(), err)
  assert.Equal(s.T(), true, CapResponse.IsGranted())
  assert.Equal(s.T(), azns.JsonB(nil), CapResponse.GetCookie())
  assert.Equal(s.T(), false, CapResponse.IsByOwnership())
  assert.Equal(s.T(), true, CapResponse.IsByGrant())
}

func (s *GrantIntegrationSuite) TestCapabilityByIndirectSubjectGrant() {
  CapResponse, err := azns.CheckCapability(s.User2.GetID(), azns.AznBasicUpdate.ID, s.Thing1B.GetID(), rdb.Connect())
  require.NoError(s.T(), err)
  assert.Equal(s.T(), true, CapResponse.IsGranted())
  assert.Equal(s.T(), azns.JsonB(nil), CapResponse.GetCookie())
  assert.Equal(s.T(), false, CapResponse.IsByOwnership())
  assert.Equal(s.T(), true, CapResponse.IsByGrant())
}

func (s *GrantIntegrationSuite) TestCapabilityByDoubleIndirectGrant() {
  CapResponse, err := azns.CheckCapability(s.User2.GetID(), azns.AznBasicUpdate.ID, s.Thing1C.GetID(), rdb.Connect())
  require.NoError(s.T(), err)
  assert.Equal(s.T(), true, CapResponse.IsGranted())
  assert.Equal(s.T(), azns.JsonB(nil), CapResponse.GetCookie())
  assert.Equal(s.T(), false, CapResponse.IsByOwnership())
  assert.Equal(s.T(), true, CapResponse.IsByGrant())
}
