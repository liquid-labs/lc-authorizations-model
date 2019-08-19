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
  . "github.com/Liquid-Labs/lc-containers-model/go/containers"
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

  // other stuff we don't own
  otherAuthID := strkit.RandString(strkit.LettersAndNumbers, 16)
  otherUser := NewUser(`users`, `Other user`, ``, otherAuthID, `444-55-5555`, `SSN`, true)
  require.NoError(s.T(), otherUser.CreateRaw(db))
  // log.Printf("Other user: %s", s.User.GetID())

  ungranted := NewEntity(`entities`, `Other Thing No Grant`, ``, otherUser.GetID(), false)
  require.NoError(s.T(), CreateEntityRaw(ungranted, db))

  directThing := NewEntity(`entities`, `Other Thing Direct Grant`, ``, otherUser.GetID(), false)
  require.NoError(s.T(), CreateEntityRaw(directThing, db))
  directGrant := azns.NewGrant(s.User.GetID(), azns.AznBasicRead.ID, directThing.GetID(), nil)
  require.NoError(s.T(), directGrant.CreateRaw(db))

  containedThing := NewEntity(`entities`, `Other Thing in Container`, ``, otherUser.GetID(), false)
  require.NoError(s.T(), CreateEntityRaw(containedThing, db))
  thingGroup := &Container{
    Entity:Entity{ResourceName: `containers`, Name:`Thing Group`, OwnerID:otherUser.GetID()},
    Members:[]*Entity{containedThing},
  }
  require.NoError(s.T(), thingGroup.CreateRaw(db))
  containedGrant := azns.NewGrant(s.User.GetID(), azns.AznBasicRead.ID, thingGroup.GetID(), nil)
  require.NoError(s.T(), containedGrant.CreateRaw(db))

  userGroup := &azns.UserGroup{
    Container:Container{
      Entity:Entity{ResourceName: `containers`, Name:`User Group`, OwnerID:otherUser.GetID()},
      Members:[]*Entity{&s.User.Entity},
    },
  }
  require.NoError(s.T(), userGroup.CreateRaw(db))
  userThing := NewEntity(`entities`, `Other Thing User Grant`, ``, otherUser.GetID(), false)
  require.NoError(s.T(), CreateEntityRaw(userThing, db))
  userGrant := azns.NewGrant(userGroup.GetID(), azns.AznBasicRead.ID, userThing.GetID(), nil)
  require.NoError(s.T(), userGrant.CreateRaw(db))

  doubleThing := NewEntity(`entities`, `Other User Thing in Container`, ``, otherUser.GetID(), false)
  require.NoError(s.T(), CreateEntityRaw(doubleThing, db))
  doubleGroup := &Container{
    Entity:Entity{ResourceName: `containers`, Name:`Double Group`, OwnerID:otherUser.GetID()},
    Members:[]*Entity{doubleThing},
  }
  require.NoError(s.T(), doubleGroup.CreateRaw(db))
  doubleGrant := azns.NewGrant(userGroup.GetID(), azns.AznBasicRead.ID, doubleGroup.GetID(), nil)
  require.NoError(s.T(), doubleGrant.CreateRaw(db))

  ctx := context.Background()
  authenticator := &auth.Authenticator{}
  authenticator.SetAznID(authID)
  s.CTX = context.WithValue(ctx, auth.AuthenticatorKey, authenticator)
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

func (s *ItemsIntegrationSuite) TestListSharedItems() {
  items := make([]Entity, 0)
  count, err := azns.ListSharedItemsQuery(&items, rdb.Connect(), azns.PageRequest{0, 3}, s.CTX)
  require.NoError(s.T(), err)
  assert.Equal(s.T(), 6, count)
  assert.Equal(s.T(), 3, len(items))
  assert.Equal(s.T(), `Double Group`, items[0].GetName())
  assert.Equal(s.T(), `Other Thing in Container`, items[2].GetName())

  count, err = azns.ListSharedItemsQuery(&items, rdb.Connect(), azns.PageRequest{1, 3}, s.CTX)
  require.NoError(s.T(), err)
  assert.Equal(s.T(), 6, count)
  assert.Equal(s.T(), 3, len(items))
  assert.Equal(s.T(), `Other Thing User Grant`, items[0].GetName())
  assert.Equal(s.T(), `Thing Group`, items[2].GetName())
}
