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
	return ue.Repository().Count(&sql.CountOptions{
		Where: &[]sql.Where{
			{
				Query: "role = ?",
				Args:  []interface{}{"ADMIN"},
			},
		},
	})
}

func (m *Module) getUserListService(search *searchOption, pagination *paginationOption) (*[]*ue.UserModel, *paginationQuery, error) {
	where := new([]sql.FindAllWhere)
	limit := new(int)

	if search != nil {
		if search.byRole != nil {
			*where = append(*where, sql.FindAllWhere{
				Where: sql.Where{
					Query: "role = ?",
					Args:  []interface{}{search.byRole},
				},
				IncludeInCount: true,
			})
		}
	}

	if pagination != nil {
		if pagination.lastID != nil {
			userData, err := m.getUserService(pagination.lastID)
			if err != nil {
				return nil, nil, err
			}
			*where = append(*where, sql.FindAllWhere{
				Where: sql.Where{
					Query: "created_at < ?",
					Args:  []interface{}{userData.CreatedAt},
				},
				IncludeInCount: false,
			})
		}

		if pagination.limit != nil {
			*limit = *pagination.limit
		}
	}

	data, page, err := ue.Repository().FindAll(&sql.FindAllOptions{
		Where: where,
		Limit: limit,
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
	return ue.Repository().FindOne(&sql.FindOneOptions{
		Where: &[]sql.Where{
			{
				Query: "id = ?",
				Args:  []interface{}{id},
			},
		},
	})
}

func (m *Module) updateUserService(id *uuid.UUID, data *ue.UserModel) (*ue.UserModel, error) {
	if _, err := ue.Repository().Update(data, &sql.UpdateOptions{
		Where: &[]sql.Where{
			{
				Query: "id = ?",
				Args:  []interface{}{id},
			},
		},
	}); err != nil {
		return nil, err
	}

	return m.getUserService(id)
}

func (*Module) deleteUserService(id *uuid.UUID) error {
	return ue.Repository().Destroy(&sql.DestroyOptions{
		Where: &[]sql.Where{
			{
				Query: "id = ?",
				Args:  []interface{}{id},
			},
		},
	})
}
