package usm

import (
	"testing"

	"github.com/nortoo/usm/model"
	"github.com/nortoo/usm/types"
	"gorm.io/gorm"
)

func TestClient_UserAll(t *testing.T) {
	app := &model.Application{
		Name:      "web",
		APPID:     "abc",
		SecretKey: "abc",
	}
	err := client.CreateApplication(app)
	if err != nil {
		t.Fatal("failed to create application:", err)
	}
	defer func() {
		if err = client.DeleteApplication(app); err != nil {
			t.Log("failed to delete application:", err)
		}
	}()

	menu := &model.Menu{
		ParentID: 0,
		Name:     "user",
		Path:     "/users",
	}
	err = client.CreateMenu(menu)
	if err != nil {
		t.Fatal("failed to create menu:", err)
	}
	defer func() {
		if err = client.DeleteMenu(menu); err != nil {
			t.Log("failed to delete menu:", err)
		}
	}()

	permission := &model.Permission{
		Action:   "GET",
		Resource: "/api/v1/users",
		Comment:  "this is a permission for users that allow user to query user list.",
	}
	err = client.CreatePermission(permission)
	if err != nil {
		t.Fatal("failed to create permission:", err)
	}
	defer func() {
		if err = client.DeletePermission(permission); err != nil {
			t.Log("failed to delete permission:", err)
		}
	}()

	roles := []*model.Role{
		{
			Name:        "admin",
			Comment:     "this is a test role.",
			Application: app,
			Menus:       []*model.Menu{menu},
			Permissions: []*model.Permission{permission},
		},
		{
			Name:        "user",
			Comment:     "this is a test role for user.",
			Application: app,
			Menus:       []*model.Menu{menu},
			Permissions: []*model.Permission{permission},
		},
	}
	for _, role := range roles {
		err = client.CreateRole(role)
		if err != nil {
			t.Fatal("failed to create role:", err)
		}
		defer func() {
			if err = client.DeleteRole(role); err != nil {
				t.Log("failed to delete role:", err)
			}
		}()
	}

	group := &model.Group{
		Name:    "test group",
		Comment: "",
	}
	err = client.CreateGroup(group)
	if err != nil {
		t.Fatal("failed to create group:", err)
	}
	defer func() {
		if err = client.DeleteGroup(group); err != nil {
			t.Log("failed to delete group:", err)
		}
	}()

	samples := []*model.User{
		{
			Username: "user-1",
			Password: "password1",
			Email:    "test@example.com",
			Mobile:   "1234567890",
			Roles:    []*model.Role{roles[0]},
			Groups:   []*model.Group{group},
			State:    0,
		},
		{
			Username: "user-2",
			Password: "password2",
			Email:    "tes2@example.com",
			Mobile:   "0987654321",
			Roles:    []*model.Role{roles[1]},
			Groups:   []*model.Group{group},
		},
	}
	for _, sample := range samples {
		err := client.CreateUser(sample)
		if err != nil {
			t.Fatal("failed to create user:", err)
		}
	}

	for _, sample := range samples {
		user, err := client.GetUser(&model.User{
			Username: sample.Username,
			Roles:    []*model.Role{},
			Groups:   []*model.Group{},
		}, "Username")
		if err != nil {
			t.Fatal("failed to get user:", err)
		}
		t.Logf("origin user: %+v", user)

		user.Username = user.Username + "-updated"
		user.Email = user.Email + "-updated"
		user.Roles = roles
		err = client.UpdateUser(user, "Username", "Email", "Roles")
		if err != nil {
			t.Fatal("failed to update user:", err)
		}

		newUser, err := client.GetUser(&model.User{Model: gorm.Model{ID: user.ID},
			Roles:  []*model.Role{},
			Groups: []*model.Group{},
		}, "ID")
		if err != nil {
			t.Fatal("failed to get updated user:", err)
		}
		t.Logf("updated user: %+v", user)

		if newUser.Username != user.Username ||
			newUser.Email != user.Email ||
			len(newUser.Roles) != len(roles) {
			t.Fatal("user update failed")
		}
	}

	users, total, err := client.ListUsers(&types.QueryUserOptions{
		Username: "us",
		Pagination: &types.Pagination{
			Page:     1,
			PageSize: 10,
		},
	})
	if err != nil {
		t.Fatal("failed to list users:", err)
	}
	if total != int64(len(samples)) {
		t.Fatalf("expected total %d, got %d", len(samples), total)
	}
	t.Logf("Total users: %d", total)

	for _, user := range users {
		err = client.DeleteUser(user)
		if err != nil {
			t.Fatal("failed to delete user:", err)
		}

		u, err := client.GetUser(&model.User{Model: gorm.Model{ID: user.ID}}, "ID")
		if err == nil {
			t.Fatalf("expected error when getting deleted user, but got none: %v", u)
		}
	}
}
