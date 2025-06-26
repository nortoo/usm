package usm

import (
	"testing"

	"github.com/nortoo/usm/model"
	"github.com/nortoo/usm/types"
)

func TestClient_GroupAll(t *testing.T) {
	samples := []*model.Group{
		{
			Name: "group-1",
		},
		{
			Name: "group-2",
		},
		{
			Name: "group-3",
		},
		{
			Name: "group-4",
		},
		{
			Name: "group-5",
		},
	}

	for _, sample := range samples {
		err := client.CreateGroup(sample)
		if err != nil {
			t.Fatal("failed to create group:", err)
		}
	}

	for _, sample := range samples {
		g, err := client.GetGroup(&model.Group{
			Name: sample.Name,
		}, "Name")
		if err != nil {
			t.Fatal("failed to get group:", err)
		}

		g.Name = g.Name + "-updated"
		err = client.UpdateGroup(g, "Name")
		if err != nil {
			t.Fatal("failed to update group:", err)
		}

		newGroup, err := client.GetGroup(&model.Group{ID: g.ID}, "ID")
		if err != nil {
			t.Fatal("failed to get updated group:", err)
		}
		if newGroup.Name != g.Name {
			t.Fatalf("expected group name %s, got %s", g.Name, newGroup.Name)
		}
	}

	groups, total, err := client.ListGroups(&types.QueryGroupOptions{Pagination: &types.Pagination{
		Page:     1,
		PageSize: 10,
	},
		WithTotal: true,
	})
	if err != nil {
		t.Fatal("failed to list groups:", err)
	}
	if total != int64(len(groups)) {
		t.Fatalf("expected total %d, got %d", len(samples), total)
	}
	t.Logf("Total groups: %d", total)

	for _, group := range groups {
		err := client.DeleteGroup(group)
		if err != nil {
			t.Fatal("failed to delete group:", err)
		}

		g, err := client.GetGroup(group, "ID")
		if err == nil {
			t.Fatalf("expected error when getting deleted group, but got none: %v", g)
		}
	}
}
