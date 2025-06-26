package usm

import (
	"fmt"
	"testing"

	"github.com/nortoo/usm/model"
	"github.com/nortoo/usm/types"
)

func TestClient_ApplicationAll(t *testing.T) {
	samples := []*model.Application{
		{
			Name:      "test-1",
			APPID:     "aaaaaa",
			SecretKey: "aaaaaa",
			Comment:   "",
			State:     0,
		},
		{
			Name:      "test-2",
			APPID:     "bbbbbb",
			SecretKey: "bbbbbb",
			Comment:   "",
			State:     0,
		},
		{
			Name:      "test-3",
			APPID:     "cccccc",
			SecretKey: "cccccc",
			Comment:   "",
			State:     0,
		},
	}

	for _, sample := range samples {
		err := client.CreateApplication(sample)
		if err != nil {
			t.Fatal("failed to create application:", err)
		}
	}

	for _, sample := range samples {
		app, err := client.GetApplication(&model.Application{
			Name: sample.Name,
		}, "Name")
		if err != nil {
			t.Fatal("failed to get application:", err)
		}

		app.Name = fmt.Sprintf("%s-%s", sample.Name, "updated")
		app.State = 1
		err = client.UpdateApplication(app, "Name", "State")
		if err != nil {
			t.Fatal("failed to update application:", err)
		}

		newApp, err := client.GetApplication(&model.Application{
			ID: app.ID,
		}, "ID")
		if err != nil {
			t.Fatal("failed to get application:", err)
		}
		if newApp.Name != app.Name || newApp.State != app.State {
			t.Fatal("failed to update application")
		}
	}

	apps, total, err := client.ListApplications(&types.QueryApplicationOptions{Pagination: &types.Pagination{
		Page:     1,
		PageSize: 10,
	},
		WithTotal: true,
	})
	if err != nil {
		t.Fatal("failed to list applications:", err)
	}
	if total != int64(len(samples)) {
		t.Fatal("get the records total count wrong:", total)
	}
	t.Logf("Total applications: %d", total)

	for _, app := range apps {
		err := client.DeleteApplication(app)
		if err != nil {
			t.Fatal("failed to delete application:", err)
		}

		// Verify deletion
		deletedApp, err := client.GetApplication(app, "ID")
		if err == nil {
			t.Fatal("expected error when getting deleted application, but got none")
		}
		if deletedApp != nil {
			t.Fatal("expected deleted application to be nil, but got:", deletedApp)
		}
	}
}
