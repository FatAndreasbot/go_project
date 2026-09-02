package postgresadapter_test

import (
	"slices"
	"strconv"
	"testing"

	postgresadapter "github.com/FatAndreasbot/go_project/user_service/adapters/outbound/postgres_adapter"
	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/google/uuid"
)

func getMockData() (oldPerms, newPerms, allPerms []*models.Permission) {
	allPerms = []*models.Permission{
		{ID: uuid.MustParse("dd612c86-2a0e-4d4c-8034-b87ba152380f"), Name: "0"},
		{ID: uuid.MustParse("aa2e7d57-96c3-4334-9e07-89f9b788b88f"), Name: "1"},
		{ID: uuid.MustParse("705d9b45-6137-4d14-a96c-8c5d47ca2586"), Name: "2"},
		{ID: uuid.MustParse("b7c99521-616c-4314-a6fd-bc59ea7c2781"), Name: "3"},
		{ID: uuid.MustParse("1cb052eb-ef3c-4ab4-a336-488780b67f13"), Name: "4"},
	}

	oldPerms = []*models.Permission{
		allPerms[1],
		allPerms[2],
		allPerms[3],
	}

	newPerms = []*models.Permission{
		allPerms[2],
		allPerms[3],
		allPerms[4],
	}

	return
}

func TestPermissionDifference1(t *testing.T) {
	oldPerms, newPerms, _ := getMockData()
	toAdd, toRemove := postgresadapter.PermissionDifference(oldPerms, newPerms)

	if !slices.Equal(toRemove, []*models.Permission{oldPerms[0]}) {
		t.Error(
			"toRemove",
			func() []string {
				var res []string
				for _, val := range toRemove {
					res = append(res, val.Name)
				}
				return res
			}(),
		)
	}

	if !slices.Equal(toAdd, []*models.Permission{newPerms[2]}) {
		t.Error(
			"toAdd",
			func() []string {
				var res []string
				for _, val := range toAdd {
					res = append(res, val.Name)
				}
				return res
			}(),
		)
	}
}

func TestPermissionDifference2(t *testing.T) {
	oldPerms, newPerms, _ := getMockData()

	newPerms = []*models.Permission{}

	toAdd, toRemove := postgresadapter.PermissionDifference(oldPerms, newPerms)

	slices.SortFunc(toRemove, func(a, b *models.Permission) int {
		a_number, _ := strconv.Atoi(a.Name)
		b_number, _ := strconv.Atoi(b.Name)
		return a_number - b_number
	})

	if !slices.Equal(toRemove, oldPerms) {
		t.Error(
			"toRemove",
			func() []string {
				var res []string
				for _, val := range toRemove {
					res = append(res, val.Name)
				}
				return res
			}(),
		)
	}

	if !slices.Equal(toAdd, []*models.Permission{}) {
		t.Error(
			"toAdd",
			func() []string {
				var res []string
				for _, val := range toAdd {
					res = append(res, val.Name)
				}
				return res
			}(),
		)
	}
}
