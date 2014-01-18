package modelview

type PaginatedListView struct {
	cursor string `json:",omitempty"`
	items  []interface{}
}

func (lv BasicPaginatedListView) Cursor() {
	return lv.cursor
}

func (lv BasicPaginatedListView) Items() {
	return lv.items
}

func (lv *BasicPaginatedListView) SetCursor(cursor string) {
	lv.cursor = cursor
}

func (lv *BasicPaginatedListView) SetItems(items []interface{}) {
	lv.items = items
}
