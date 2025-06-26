package usm

import (
	"testing"

	"github.com/nortoo/usm/model"
	"github.com/nortoo/usm/types"
)

func TestClient_PermissionAll(t *testing.T) {
	samples := []*model.Permission{
		{
			Action:   "GET",
			Resource: "/api/v1/test-1",
		},
		{
			Action:   "POST",
			Resource: "/api/v1/test-2",
		},
		{
			Action:   "PUT",
			Resource: "/api/v1/test-3",
		},
		{
			Action:   "DELETE",
			Resource: "/api/v1/test-4",
		},
	}

	for _, sample := range samples {
		err := client.CreatePermission(sample)
		if err != nil {
			t.Fatal("failed to create permission:", err)
		}
	}

	for _, sample := range samples {
		p, err := client.GetPermission(&model.Permission{
			Action:   sample.Action,
			Resource: sample.Resource,
		}, "Action", "Resource")
		if err != nil {
			t.Fatal("failed to get permission:", err)
		}

		p.Action = p.Action + "-updated"
		p.Resource = p.Resource + "-updated"
		err = client.UpdatePermission(p, "Action", "Resource")
		if err != nil {
			t.Fatal("failed to update permission:", err)
		}

		newPerm, err := client.GetPermission(&model.Permission{ID: p.ID}, "ID")
		if err != nil {
			t.Fatal("failed to get updated permission:", err)
		}

		if newPerm.Action != p.Action || newPerm.Resource != p.Resource {
			t.Fatal("permission update failed")
		}
	}

	// List permissions
	permissions, total, err := client.ListPermissions(&types.QueryPermissionOptions{
		Pagination: &types.Pagination{
			Page:     1,
			PageSize: 10,
		},
		WithTotal: true,
	})
	if err != nil {
		t.Fatal("failed to list permissions:", err)
	}
	if total != int64(len(samples)) {
		t.Fatalf("expected %d permissions, got %d", len(samples), total)
	}
	t.Logf("Total permissions: %d", total)

	for _, perm := range permissions {
		err = client.DeletePermission(perm)
		if err != nil {
			t.Fatal("failed to delete permission:", err)
		}

		p, err := client.GetPermission(perm, "ID")
		if err == nil {
			t.Fatalf("expected error when getting deleted permission, but got none: %v", p)
		}
	}
}
