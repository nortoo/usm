package usm

import (
	"github.com/nortoo/usm/model"
	"github.com/nortoo/usm/types"
	"github.com/pkg/errors"
	"gorm.io/gorm/clause"
)

func (c *Client) CreateRole(r *model.Role) error {
	err := c.db.Create(r).Error
	if err != nil {
		return err
	}

	// assign permissions to the role.
	for _, p := range r.Permissions {
		err = c.addPolicy(r.Name, r.Application.Name, p.Resource, p.Action)
		if err != nil {
			// if failed to create policy, delete the role and return error.
			// this is to ensure the role is not created with incomplete policies.
			_ = c.clearPolicy(r.Name)
			_ = c.DeleteRole(r)
			return errors.Errorf("failed to create policy for role %s: %v", r.Name, err)
		}
	}
	return nil
}

func (c *Client) DeleteRole(r *model.Role) error {
	role, err := c.GetRole(r)
	if err != nil {
		return err
	}
	err = c.db.Select(clause.Associations).Delete(role).Error
	if err != nil {
		return err
	}
	return c.clearPolicy(role.Name)
}

func (c *Client) UpdateRole(r *model.Role, cols ...string) error {
	err := c.db.Model(r).Select(cols).Updates(r).Error
	if err != nil {
		return err
	}

	// update policies for the role.
	_ = c.clearPolicy(r.Name)
	for _, p := range r.Permissions {
		_ = c.addPolicy(r.Name, r.Application.Name, p.Resource, p.Action)
	}
	return nil
}

// GetRole retrieves a role by its attributes.
// If the given associations are provided in the query model,
// it will preload the association.
// For example, if the Application is not nil in the query model r,
// it will preload the Application association.
func (c *Client) GetRole(r *model.Role, cols ...interface{}) (*model.Role, error) {
	tx := c.db.Where(r, cols)
	if r.Application != nil {
		tx.Preload("Application")
	}
	if r.Menus != nil {
		tx.Preload("Menus")
	}
	if r.Permissions != nil {
		tx.Preload("Permissions")
	}
	err := tx.First(r).Error
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (c *Client) ListRoles(q *types.QueryRoleOptions) (ret []*model.Role, total int64, err error) {
	tx := c.db
	if len(q.IsDefault) > 0 {
		tx = tx.Where("is_default IN ?", q.IsDefault)
	}
	if q.Pagination != nil {
		if q.WithTotal {
			err = tx.Model(&model.Role{}).Count(&total).Error
			if err != nil || total == 0 {
				return
			}
		}
		tx.Limit(q.Pagination.PageSize).Offset((q.Pagination.Page - 1) * q.Pagination.PageSize)
	}
	err = tx.Find(&ret).Error
	return
}
