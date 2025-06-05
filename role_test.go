package usm

import (
	"testing"

	"github.com/nortoo/usm/model"
	"github.com/nortoo/usm/types"
)

func TestClient_RoleAll(t *testing.T) {
	app := &model.Application{
		Name:      "web-test",
		APPID:     "test-id",
		SecretKey: "test-key",
	}
	err := client.CreateApplication(app)
	if err != nil {
		t.Fatal("failed to create application:", err)
	}
	defer func() {
		// Clean up the application created for this test
		if err := client.DeleteApplication(app); err != nil {
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

	samples := []*model.Role{
		{
			Name:        "role-1",
			Comment:     "role 1 comment",
			Application: app,
			Menus:       []*model.Menu{menu},
		},
		{
			Name:        "role-2",
			Comment:     "role 2 comment",
			Application: app,
			Permissions: []*model.Permission{permission},
		},
		{
			Name:        "role-3",
			Comment:     "role 3 comment",
			Application: app,
			Menus:       []*model.Menu{menu},
			Permissions: []*model.Permission{permission},
		},
	}

	for _, sample := range samples {
		err := client.CreateRole(sample)
		if err != nil {
			t.Fatal("failed to create role:", err)
		}
	}
	defer func() {
		for _, sample := range samples {
			if err := client.clearPolicy(sample.Name); err != nil {
				t.Log("failed to clear policy:", err)
			}
		}
	}()

	for _, sample := range samples {
		role, err := client.GetRole(&model.Role{
			Name:        sample.Name,
			Application: &model.Application{},
			Menus:       []*model.Menu{},
			Permissions: []*model.Permission{},
		}, "Name")
		if err != nil {
			t.Fatal("failed to get role:", err)
		}

		for _, perm := range role.Permissions {
			if ok, err := client.Authorize(role.Name, role.Application.Name, perm.Resource, perm.Action); err != nil || !ok {
				t.Fatal("failed to authorize role:", err)
			}
		}

		role.Name = role.Name + "-updated"
		role.Comment = role.Comment + " updated"
		err = client.UpdateRole(role, "Name", "Comment")
		if err != nil {
			t.Fatal("failed to update role:", err)
		}

		newRole, err := client.GetRole(&model.Role{ID: role.ID}, "ID")
		if err != nil {
			t.Fatal("failed to get updated role:", err)
		}

		if newRole.Name != role.Name || newRole.Comment != role.Comment {
			t.Fatal("role update failed")
		}
	}

	roleList, total, err := client.ListRoles(&types.QueryRoleOptions{Pagination: &types.Pagination{
		Page:     1,
		PageSize: 10,
	}})
	if err != nil {
		t.Fatal("failed to list roles:", err)
	}
	if total != int64(len(samples)) {
		t.Fatal("role list count mismatch")
	}
	t.Logf("role list count: %d", len(roleList))

	for _, role := range roleList {
		err := client.DeleteRole(role)
		if err != nil {
			t.Fatal("failed to delete role:", err)
		}

		r, err := client.GetRole(role, "ID")
		if err == nil {
			t.Fatalf("expected error when getting deleted role, but got none: %v", r)
		}
	}
}
