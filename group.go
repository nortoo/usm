package usm

import "github.com/nortoo/usm/model"

func (c *Client) CreateGroup(g *model.Group) error {
	return c.db.Create(g).Error
}

func (c *Client) DeleteGroup(g *model.Group) error {
	return c.db.Delete(g).Error
}

func (c *Client) UpdateGroup(g *model.Group, cols ...string) error {
	return c.db.Model(g).Select(cols).Updates(g).Error
}

func (c *Client) GetGroup(g *model.Group, cols ...interface{}) (*model.Group, error) {
	err := c.db.Where(g, cols).First(g).Error
	if err != nil {
		return nil, err
	}
	return g, nil
}

type QueryGroupOptions struct {
	Pagination *model.Pagination
}

func (c *Client) ListGroups(q *QueryGroupOptions) (ret []*model.Group, total int64, err error) {
	tx := c.db
	if q.Pagination != nil {
		err = tx.Model(&model.Group{}).Count(&total).Error
		if err != nil || total == 0 {
			return
		}

		tx.Limit(q.Pagination.PageSize).Offset((q.Pagination.Page - 1) * q.Pagination.PageSize)
	}
	err = tx.Find(&ret).Error
	return
}
