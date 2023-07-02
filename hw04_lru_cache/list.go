package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len  int
	head *ListItem
	tail *ListItem

	List // Remove me after realization.
	// Place your code here.
}

func (l list) Len() int {
	return l.len
}

func (l list) Front() *ListItem {
	return l.head
}

func (l list) Back() *ListItem {
	return l.tail
}

func (l list) PushFront(v interface{}) *ListItem {
	return nil
}

func (l list) PushBack(v interface{}) *ListItem {
	return nil
}

func (l list) Remove(i *ListItem) {

}

func (l list) MoveToFront(i *ListItem) {

}

func NewList() List {
	// todo: re-use NewFilledList ?
	return new(list)
}

func NewFilledList(elems []interface{}) List {
	if len(elems) == 0 {
		return list{}
	}

	var cur, prev, next *ListItem = &ListItem{Value: elems[0]}, nil, nil
	list := new(list)
	list.head = cur

	for i := range elems {
		if i+1 < len(elems) {
			next = &ListItem{Value: elems[i+1]}
		}

		cur.Prev = prev
		cur.Next = next

		prev = cur
		cur = next
		next = nil
	}

	list.tail = prev
	list.len = len(elems)

	return list
}
