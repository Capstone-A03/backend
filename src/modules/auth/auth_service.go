package auth

import (
	"capstonea03/be/src/libs/db/sql"
	ue "capstonea03/be/src/modules/user/user_entity"

	"github.com/google/uuid"
)

func (m *Module) getUserService(id *uuid.UUID) (*ue.UserModel, error) {
	return ue.Repository().FindOne(&sql.FindOneOptions{
		Where: &[]sql.Where{
			{
				Query: "id = ?",
				Args:  []interface{}{id},
			},
		},
	})
}

func (m *Module) getUserByUsernameService(username *string) (*ue.UserModel, error) {
	return ue.Repository().FindOne(&sql.FindOneOptions{
		Where: &[]sql.Where{
			{
				Query: "username = ?",
				Args:  []interface{}{username},
			},
		},
	})
}

func (m *Module) addUserService(data *ue.UserModel) (*ue.UserModel, error) {
	return ue.Repository().Create(data)
}
