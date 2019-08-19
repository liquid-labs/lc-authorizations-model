package azns_test

import (
  // "log"
  "os"
  "testing"

  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/require"
  "github.com/stretchr/testify/suite"

  . "github.com/Liquid-Labs/lc-containers-model/go/containers"
  . "github.com/Liquid-Labs/lc-entities-model/go/entities"
  "github.com/Liquid-Labs/lc-rdb-service/go/rdb"
  "github.com/Liquid-Labs/strkit/go/strkit"
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
}

type GrantIntegrationSuite struct {
  suite.Suite
  User1            *User
  Thing1A          *Entity
  Thing1B          *Entity
  Thing1C          *Entity
  Thing1Outer      *Entity
  Thing1GroupA     *Container
  Thing1GroupC     *Container
  Thing1GroupOuter *Container
  User1Group       *azns.UserGroup
  User2            *User
  Thing2           *Entity
}
// SetupSuite sets up the following capabilities:
//
// * User1 owns Thing1, Thing1Outer, Thing1GroupA & C, Thing1GroupOuter, and User1Group
// * User2 owns Thing2
// * User1 has direct grant over Thing2
// * User2 has grant over Thing1GroupA containing Thing1A
// * User2 in is User1Group with grant Thing1B
// * User2 in is User1Group with grant Thing1GroupC containing Thing1C
// * Thing1GroupOuter contains Thing1Outer, Thing1A, Thing1B, Thing1C, Thing1GroupA, and Thing1GroupC.
func (s *GrantIntegrationSuite) SetupSuite() {
  db := rdb.Connect()

  // setup base objects
  authID1 := strkit.RandString(strkit.LettersAndNumbers, 16)
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
  s.Thing1Outer = NewEntity(`entities`, `Thing1Outer`, ``, s.User1.GetID(), false)
  require.NoError(s.T(), CreateEntityRaw(s.Thing1Outer, db))
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
  s.Thing1GroupOuter = &Container{
    Entity:Entity{ResourceName: `containers`, Name:`Thing1GroupOuter`, OwnerID:s.User1.GetID()},
    Members:[]*Entity{s.Thing1Outer, s.Thing1A, s.Thing1B, s.Thing1C, &s.Thing1GroupA.Entity, &s.Thing1GroupC.Entity},
  }
  require.NoError(s.T(), s.Thing1GroupOuter.CreateRaw(db))
  // log.Printf("Thing2: %s", s.Thing2Group.GetID())

  authID2 := strkit.RandString(strkit.LettersAndNumbers, 16)
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
  for _, g := range grants { require.NoError(s.T(), g.CreateRaw(db)) }
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

func (s *GrantIntegrationSuite) TestNoCapabilityToOuterGroup() {
  for _, target := range []EID{s.Thing1Outer.GetID(), s.Thing1GroupOuter.GetID()} {
    CapResponse, err := azns.CheckCapability(s.User2.GetID(), azns.AznBasicUpdate.ID, target, rdb.Connect())
    require.NoError(s.T(), err)
    assert.Equal(s.T(), false, CapResponse.IsGranted())
    assert.Equal(s.T(), azns.JsonB(nil), CapResponse.GetCookie())
    assert.Equal(s.T(), false, CapResponse.IsByOwnership())
    assert.Equal(s.T(), false, CapResponse.IsByGrant())
  }
}
