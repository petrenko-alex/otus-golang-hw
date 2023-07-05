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

	// List // Remove me after realization.
	// Place your code here.
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	l.pushItemToFront(&ListItem{
		Value: v,
		Prev:  nil,
		Next:  nil,
	})

	return l.Front()
}

func (l *list) pushItemToFront(i *ListItem) {
	// connect with previous head
	if l.Front() != nil {
		i.Next = l.Front()
		l.Front().Prev = i
	}

	// set new head
	l.head = i

	// set new tail in case of empty list
	if l.Back() == nil {
		l.tail = i
	}

	l.len++
}

func (l *list) PushBack(v interface{}) *ListItem {
	l.pushItemToBack(&ListItem{
		Value: v,
		Prev:  nil,
		Next:  nil,
	})

	return l.Back()
}

func (l *list) pushItemToBack(i *ListItem) {
	// connect with previous tail
	if l.Back() != nil {
		i.Prev = l.Back()
		l.Back().Next = i
	}

	// set new tail
	l.tail = i

	// set new head in case of empty list
	if l.Front() == nil {
		l.head = i
	}

	l.len++
}

func (l *list) Remove(i *ListItem) {
	if i == nil {
		return
	}

	prevItem, nextItem := i.Prev, i.Next

	if prevItem != nil {
		prevItem.Next = nextItem
	}
	if nextItem != nil {
		nextItem.Prev = prevItem
	}

	i.Next, i.Prev = nil, nil

	// in case head is deleted
	if l.Front() == i {
		l.head = nextItem
	}

	// in case back is deleted
	if l.Back() == i {
		l.tail = prevItem
	}

	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == nil || l.len == 0 {
		return
	}

	l.Remove(i)
	l.pushItemToFront(i)
}

func NewList() List {
	return NewFilledList([]interface{}{})
}

func NewFilledList(elems []interface{}) List {
	if len(elems) == 0 {
		return &list{}
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
