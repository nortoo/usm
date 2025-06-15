package usm

import (
	"fmt"

	"github.com/nortoo/usm/model"
	"github.com/nortoo/usm/types"
)

func (c *Client) CreateUser(u *model.User) error {
	return c.db.Create(u).Error
}

func (c *Client) DeleteUser(u *model.User) error {
	// TODO: solve re-register user issue
	// which duplicated username or email could occur
	// when using a email from deleted user.
	return c.db.Delete(u).Error
}

func (c *Client) UpdateUser(u *model.User, cols ...string) error {
	return c.db.Model(u).Select(cols).Updates(u).Error
}

func (c *Client) GetUser(u *model.User, cols ...interface{}) (*model.User, error) {
	tx := c.db.Where(u, cols...)
	if u.Roles != nil {
		tx.Preload("Roles")
	}
	if u.Groups != nil {
		tx.Preload("Groups")
	}
	err := tx.First(u).Error
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (c *Client) ListUsers(q *types.QueryUserOptions) (ret []*model.User, total int64, err error) {
	tx := c.db.Model(&model.User{})
	if q.Username != "" {
		tx = tx.Where("username LIKE ?", fmt.Sprintf("%%%s%%", q.Username))
	}
	if q.Email != "" {
		tx = tx.Where("email LIKE ?", q.Email)
	}
	if q.Mobile != "" {
		tx = tx.Where("mobile LIKE ?", q.Mobile)
	}
	if len(q.States) > 0 {
		tx = tx.Where("state IN ?", q.States)
	}
	if q.RoleID != 0 {
		tx = tx.Joins("JOIN user_roles ON user_roles.user_id = users.id").
			Where("user_roles.role_id = ?", q.RoleID)
	}
	if q.GroupID != 0 {
		tx = tx.Joins("JOIN user_groups ON user_groups.user_id = users.id").
			Where("user_groups.group_id = ?", q.GroupID)
	}

	if q.Pagination != nil {
		err = tx.Count(&total).Error
		if err != nil || total == 0 {
			return
		}

		tx = tx.Limit(q.Pagination.PageSize).Offset((q.Pagination.Page - 1) * q.Pagination.PageSize)
	}

	err = tx.Find(&ret).Error
	return
}

func (c *Client) DoesUsernameExist(username string) (bool, error) {
	var exists bool
	sql := "SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)"
	err := c.db.Raw(sql, username).Scan(&exists).Error
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (c *Client) DoesEmailExist(email string) (bool, error) {
	var exists bool
	sql := "SELECT EXISTS(SELECT 1 FROM users WHERE email = ? AND deleted_at IS NULL)"
	err := c.db.Raw(sql, email).Scan(&exists).Error
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (c *Client) DoesMobileExist(mobile string) (bool, error) {
	var exists bool
	sql := "SELECT EXISTS(SELECT 1 FROM users WHERE mobile = ? and deleted_at IS NULL)"
	err := c.db.Raw(sql, mobile).Scan(&exists).Error
	if err != nil {
		return false, err
	}
	return exists, nil
}
