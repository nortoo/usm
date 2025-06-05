package usm

import (
	"github.com/nortoo/usm/model"
	"github.com/nortoo/usm/types"
)

func (c *Client) CreateApplication(app *model.Application) error {
	return c.db.Create(app).Error
}

func (c *Client) DeleteApplication(app *model.Application) error {
	return c.db.Delete(app).Error
}

func (c *Client) UpdateApplication(app *model.Application, cols ...string) error {
	return c.db.Model(app).Select(cols).Updates(app).Error
}

func (c *Client) GetApplication(app *model.Application, cols ...interface{}) (*model.Application, error) {
	err := c.db.Where(app, cols...).First(app).Error
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (c *Client) ListApplications(q *types.QueryApplicationOptions) (ret []*model.Application, total int64, err error) {
	tx := c.db
	if q.Pagination != nil {
		err = tx.Model(&model.Application{}).Count(&total).Error
		if err != nil || total == 0 {
			return
		}

		tx.Limit(q.Pagination.PageSize).Offset((q.Pagination.Page - 1) * q.Pagination.PageSize)
	}
	err = tx.Find(&ret).Error
	return
}
