package user

import (
	"capstonea03/be/src/libs/db/sql"
	uc "capstonea03/be/src/modules/user/user_constant"
	ue "capstonea03/be/src/modules/user/user_entity"

	"github.com/google/uuid"
)

type searchOption struct {
	byRole *uc.Role
}

type paginationOption struct {
	lastID *uuid.UUID
	limit  *int
}

type paginationQuery struct {
	count *int
	limit *int
	total *int
}

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

func (m *Module) getUserListService(pagination *paginationOption, search *searchOption) (*[]*ue.UserModel, *paginationQuery, error) {
	where := []sql.FindAllWhere{}
	limit := 0

	if search != nil {
		if search.byRole != nil {
			where = append(where, sql.FindAllWhere{
				Where: sql.Where{
					Query: "role = ?",
					Args:  []interface{}{search.byRole},
				},
				IncludeInCount: true,
			})
		}
	}

	if pagination.lastID != nil && len(*pagination.lastID) > 0 {
		userData, err := m.getUserService(pagination.lastID)
		if err != nil {
			return nil, nil, err
		}
		where = append(where, sql.FindAllWhere{
			Where: sql.Where{
				Query: "created_at < ?",
				Args:  []interface{}{userData.CreatedAt},
			},
			IncludeInCount: false,
		})
	}

	if pagination.limit != nil && *pagination.limit > 0 {
		limit = *pagination.limit
	}

	data, page, err := ue.UserRepository().FindAll(&sql.FindAllOptions{
		Where: &where,
		Limit: &limit,
		Order: &[]string{"created_at desc"},
	})
	if err != nil {
		return nil, nil, err
	}

	return data, &paginationQuery{
		count: &page.Count,
		limit: &page.Limit,
		total: &page.Total,
	}, nil
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
	where := []sql.Where{
		{
			Query: "id = ?",
			Args:  []interface{}{id},
		},
	}

	if _, err := ue.UserRepository().Update(data, &sql.UpdateOptions{
		Where: &where,
	}); err != nil {
		return nil, err
	}

	data, err := ue.UserRepository().FindOne(&sql.FindOneOptions{
		Where: &where,
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
