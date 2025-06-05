package usm

import (
	"github.com/nortoo/usm/model"
	"github.com/nortoo/usm/types"
)

func (c *Client) CreatePermission(p *model.Permission) error {
	return c.db.Create(p).Error
}

func (c *Client) DeletePermission(p *model.Permission) error {
	return c.db.Delete(p).Error
}

func (c *Client) UpdatePermission(p *model.Permission, cols ...string) error {
	return c.db.Model(p).Select(cols).Updates(p).Error
}

func (c *Client) GetPermission(p *model.Permission, cols ...interface{}) (*model.Permission, error) {
	err := c.db.Where(p, cols).First(p).Error
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (c *Client) ListPermissions(q *types.QueryPermissionOptions) (ret []*model.Permission, total int64, err error) {
	tx := c.db
	if q.Pagination != nil {
		err = tx.Model(&model.Permission{}).Count(&total).Error
		if err != nil || total == 0 {
			return
		}

		tx.Limit(q.Pagination.PageSize).Offset((q.Pagination.Page - 1) * q.Pagination.PageSize)
	}
	err = tx.Find(&ret).Error
	return
}
