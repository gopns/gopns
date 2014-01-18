package modelview

type PaginatedListView interface {
	Cursor() string
	Name() string
	Items() []interface{}
}

type BasicPaginatedListView struct {
	cursor string
	name   string
	items  []interface{}
}

func (lv BasicPaginatedListView) Cursor() {
	return lv.cursor
}

func (lv BasicPaginatedListView) Name() {
	return lv.name
}

func (lv BasicPaginatedListView) Items() {
	return lv.items
}

func (lv *BasicPaginatedListView) SetCursor(cursor string) {
	lv.cursor = cursor
}

func (lv *BasicPaginatedListView) SetName(name string) {
	lv.name = name
}

func (lv *BasicPaginatedListView) SetItems(items []interface{}) {
	lv.items = items
}
