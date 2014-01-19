package modelview

import "github.com/gopns/gopns/model"

type PaginatedListView struct {
	cursor string `json:",omitempty"`
	items  []interface{}
}

func NewPaginatedDeviceListView(devices []model.Device, cursor string) *PaginatedListView {
	items := make([]interface{}, len(devices))
	for i := range devices {
		items[i] = devices[i]
	}
	return &PaginatedListView{items: items, cursor: cursor}
}

func (lv PaginatedListView) Cursor() string {
	return lv.cursor
}

func (lv PaginatedListView) Items() []interface{} {
	return lv.items
}

func (lv *PaginatedListView) SetCursor(cursor string) {
	lv.cursor = cursor
}

func (lv *PaginatedListView) SetItems(items []interface{}) {
	lv.items = items
}
