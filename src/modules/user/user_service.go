package user

import (
	"capstonea03/be/src/libs/db/sql"
	ue "capstonea03/be/src/modules/user/user_entity"

	"github.com/google/uuid"
)

func (*Module) countUserWithAdminRoleService() (*int64, error) {
	return ue.UserRepository().Count(&sql.CountOptions{
		Where: &[]sql.Where{
			{
				Query: "role = ?",
				Args:  []interface{}{"ADMIN"},
			},
		},
	})
}

func (*Module) getUserService(id *uuid.UUID) (*ue.UserModel, error) {
	return ue.UserRepository().FindOne(&sql.FindOneOptions{
		Where: &[]sql.Where{
			{
				Query: "id = ?",
				Args:  []interface{}{id},
			},
		},
	})
}

func (*Module) updateUserService(id *uuid.UUID, data *ue.UserModel) (*ue.UserModel, error) {
	if _, err := ue.UserRepository().Update(data, &sql.UpdateOptions{
		Where: &[]sql.Where{
			{
				Query: "id = ?",
				Args:  []interface{}{id},
			},
		},
	}); err != nil {
		return nil, err
	}

	data, err := ue.UserRepository().FindOne(&sql.FindOneOptions{
		Where: &[]sql.Where{
			{
				Query: "id = ?",
				Args:  []interface{}{id},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (*Module) deleteUserService(id *uuid.UUID) error {
	return ue.UserRepository().Destroy(&ue.UserModel{
		Model: sql.Model{
			ID: id,
		},
	})
}
