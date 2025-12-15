package items_test

import (
	"testing"

	"tomokimura.jp/notes-app/backend/domain/items"
)

func Test_NewItem(t *testing.T) {
	item := items.NewItem("title", "content")
	if item == nil {
		t.Errorf("expected item to be not nil")
	}
}
