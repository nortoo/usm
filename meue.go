package usm

import "github.com/nortoo/usm/model"

func (c *Client) CreateMenu(m *model.Menu) error {
	return c.db.Create(m).Error
}

func (c *Client) DeleteMenu(m *model.Menu) error {
	return c.db.Delete(m).Error
}

func (c *Client) UpdateMenu(m *model.Menu, cols ...string) error {
	return c.db.Model(m).Select(cols).Updates(m).Error
}

func (c *Client) GetMenu(m *model.Menu, cols ...interface{}) (*model.Menu, error) {
	err := c.db.Where(m, cols).First(m).Error
	if err != nil {
		return nil, err
	}
	return m, nil
}

type QueryMenuOptions struct {
	Pagination *model.Pagination
}

func (c *Client) ListMenus(q *QueryMenuOptions) (ret []*model.Menu, total int64, err error) {
	tx := c.db
	if q.Pagination != nil {
		err = tx.Model(&model.Menu{}).Count(&total).Error
		if err != nil || total == 0 {
			return
		}

		tx.Limit(q.Pagination.PageSize).Offset((q.Pagination.Page - 1) * q.Pagination.PageSize)
	}
	err = tx.Find(&ret).Error
	return
}
