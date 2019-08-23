package azns

import (
  "github.com/go-pg/pg/orm"

  . "github.com/Liquid-Labs/lc-containers-model/go/containers"
  . "github.com/Liquid-Labs/lc-entities-model/go/entities"
  . "github.com/Liquid-Labs/lc-users-model/go/users"
  . "github.com/Liquid-Labs/terror/go/terror"
)

type UserGroup struct {
  tableName struct{} `sql:"select:user_groups_join_entities,alias:user_group"`
  Container
}

func (ug *UserGroup) GetResourceName() ResourceName {
  return ResourceName(`user-groups`)
}

func (ug *UserGroup) CreateRaw(db orm.DB) Terror {
  if err := ug.Container.CreateRaw(db); err != nil {
    return ServerError(`Error while creating user groups record.`, err)
  } else {
    qs := db.Model(&Subject{Entity{ID:ug.GetID()}}).
      ExcludeColumn(EntityFields...)
    if _, err := qs.Insert(); err != nil {
      return ServerError(`Problem creating user group record.`, err)
    } else {
      qu := db.Model(ug).ExcludeColumn(ContainerFields...)
      if _, err := qu.Insert(); err != nil {
        return ServerError(`Problem creating user group record.`, err)
      }
    }
  }
  return nil
}
