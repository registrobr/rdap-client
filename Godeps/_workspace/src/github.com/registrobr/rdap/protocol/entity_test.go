package protocol

import (
	"reflect"
	"testing"
)

func TestEntityGetEntity(t *testing.T) {
	roleA := Entity{Handle: "A", Roles: []string{"role-A"}}
	roleB := Entity{Handle: "B", Roles: []string{"role-B"}}
	roleABC := Entity{Handle: "ABC", Roles: []string{"role-A", "role-B", "role-C"}}

	e := Entity{
		Entities: []Entity{roleA, roleB, roleABC},
	}

	data := []struct {
		description    string
		role           string
		expectedEntity Entity
		expectedFound  bool
	}{
		{
			description:    "it should return the first item found (1)",
			role:           "role-A",
			expectedEntity: roleA,
			expectedFound:  true,
		},
		{
			description:    "it should return the first item found (2)",
			role:           "role-B",
			expectedEntity: roleB,
			expectedFound:  true,
		},
		{
			description:    "it should return the first item found (3)",
			role:           "role-C",
			expectedEntity: roleABC,
			expectedFound:  true,
		},
		{
			description:   "it should not found",
			role:          "role-D",
			expectedFound: false,
		},
	}

	for i, item := range data {
		entity, found := e.GetEntity(item.role)

		if found != item.expectedFound {
			t.Errorf("[%d] %s: expected found “%t”", i, item.description, item.expectedFound)
		}

		if !reflect.DeepEqual(item.expectedEntity, entity) {
			t.Errorf("[%d] %s: unexpected entity returned. Expected “%s” and got “%s”", i, item.description, item.expectedEntity.Handle, entity.Handle)
		}
	}

}
