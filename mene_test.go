package usm

import (
	"testing"

	"github.com/nortoo/usm/model"
)

func TestClient_MenuAll(t *testing.T) {
	samples := []*model.Menu{
		{
			Name:    "menu-1",
			Path:    "/menu-1",
			Comment: "menu 1 comment",
		},
		{
			Name:    "menu-2",
			Path:    "/menu-2",
			Comment: "menu 2 comment",
		},
		{
			Name:    "menu-3",
			Path:    "/menu-3",
			Comment: "menu 3 comment",
		},
		{
			Name:    "menu-4",
			Path:    "/menu-4",
			Comment: "menu 4 comment",
		},
		{
			Name:    "menu-5",
			Path:    "/menu-5",
			Comment: "menu 5 comment",
		},
	}

	for _, menu := range samples {
		err := client.CreateMenu(menu)
		if err != nil {
			t.Fatal("failed to create menu:", err)
		}
	}

	for _, menu := range samples {
		m, err := client.GetMenu(&model.Menu{
			Name: menu.Name,
		}, "Name")
		if err != nil {
			t.Fatal("failed to get menu:", err)
		}

		m.Name = m.Name + "-updated"
		m.Comment = m.Comment + " updated"
		err = client.UpdateMenu(m, "Name", "Comment")
		if err != nil {
			t.Fatal("failed to update menu:", err)
		}

		newMenu, err := client.GetMenu(&model.Menu{ID: m.ID}, "ID")
		if err != nil {
			t.Fatal("failed to get updated menu:", err)
		}
		if newMenu.Name != m.Name || newMenu.Comment != m.Comment {
			t.Fatalf("expected menu name %s and comment %s, got name %s and comment %s", m.Name, m.Comment, newMenu.Name, newMenu.Comment)
		}
	}

	menus, total, err := client.ListMenus(&QueryMenuOptions{Pagination: &model.Pagination{
		Page:     1,
		PageSize: 10,
	}})
	if err != nil {
		t.Fatal("failed to list menus:", err)
	}
	if total != int64(len(samples)) {
		t.Fatalf("expected total %d, got %d", len(samples), total)
	}
	t.Logf("Total menus: %d", total)

	for _, menu := range menus {
		err = client.DeleteMenu(menu)
		if err != nil {
			t.Fatal("failed to delete menu:", err)
		}

		m, err := client.GetMenu(menu, "ID")
		if err == nil {
			t.Fatalf("expected error when getting deleted menu, but got none: %v", m)
		}
	}
}
