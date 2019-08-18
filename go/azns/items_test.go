package azns_test

import (
  "context"
  "fmt"
  "os"
  "testing"

  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/require"
  "github.com/stretchr/testify/suite"

  "github.com/Liquid-Labs/lc-authentication-api/go/auth"
  . "github.com/Liquid-Labs/lc-entities-model/go/entities"
  "github.com/Liquid-Labs/lc-rdb-service/go/rdb"
  . "github.com/Liquid-Labs/lc-users-model/go/users"
  "github.com/Liquid-Labs/strkit/go/strkit"
  "github.com/Liquid-Labs/terror/go/terror"

  /* pkg2test*/ "github.com/Liquid-Labs/lc-authorizations-model/go/azns"
)

func init() {
  terror.EchoErrorLog()
}

type ItemsIntegrationSuite struct {
  suite.Suite
  CTX  context.Context
  User *User
}

func (s *ItemsIntegrationSuite) SetupSuite() {
  db := rdb.Connect()

  // setup base objects
  authID := strkit.RandString(strkit.LettersAndNumbers, 16)
  s.User = NewUser(`users`, `User`, ``, authID, `555-55-5555`, `SSN`, true)
  require.NoError(s.T(), s.User.CreateRaw(db))
  // log.Printf("User: %s", s.User.GetID())

  for i := 1; i < 20; i += 1 {
    thing := NewEntity(`entities`, fmt.Sprintf(`Thing %02d`, i), ``, s.User.GetID(), false)
    require.NoError(s.T(), CreateEntityRaw(thing, db))
  }


  // ensure there's something that shouldn't show up in any selects
  otherAuthID := strkit.RandString(strkit.LettersAndNumbers, 16)
  otherUser := NewUser(`users`, `Other user`, ``, otherAuthID, `444-55-5555`, `SSN`, true)
  require.NoError(s.T(), otherUser.CreateRaw(db))
  // log.Printf("Other user: %s", s.User.GetID())
  otherThing := NewEntity(`entities`, `Other thing`, ``, otherUser.GetID(), false)
  require.NoError(s.T(), CreateEntityRaw(otherThing, db))

  ctx := context.Background()
  authenticator := &auth.Authenticator{}
  authenticator.SetAznID(authID)
  s.CTX = context.WithValue(ctx, auth.AuthenticatorKey, authenticator)

  /*
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

  grants := []*azns.Grant{*/
    ///* U1 > T2 */ azns.NewGrant(s.User1.GetID(), azns.AznBasicUpdate.ID, s.Thing2.GetID(), nil),
    //* (U2 > T1GA) > T1A */ azns.NewGrant(s.User2.GetID(), azns.AznBasicUpdate.ID, s.Thing1GroupA.GetID(), nil),
    //* U2 > (U1G > T1B) */ azns.NewGrant(s.User1Group.GetID(), azns.AznBasicUpdate.ID, s.Thing1B.GetID(), nil),
    //* U2 > (U1G > T1GC) > T1C */ azns.NewGrant(s.User1Group.GetID(), azns.AznBasicUpdate.ID, s.Thing1GroupC.GetID(), nil),
    /*
  }
  for _, g := range grants { require.NoError(s.T(), g.CreateRaw(db)) }*/
}

func TestItemsIntegrationSuite(t *testing.T) {
  if os.Getenv(`SKIP_INTEGRATION`) == `true` {
    t.Skip()
  } else {
    suite.Run(t, new(ItemsIntegrationSuite))
  }
}

func (s *ItemsIntegrationSuite) TestListOwnedItems() {
  things := make([]Entity, 0)
  q := rdb.Connect().Model(&things).Order(`entity.name ASC`)
  pr := azns.PageRequest{Page:0, ItemsPerPage:5}

  count, err := azns.ListOwnedItems(q, pr, s.CTX)
  require.NoError(s.T(), err)
  assert.Equal(s.T(), 20, count)
  assert.Equal(s.T(), 5, len(things))
  assert.Equal(s.T(), `Thing 01`, things[0].GetName())
  assert.Equal(s.T(), `Thing 05`, things[4].GetName())
}
