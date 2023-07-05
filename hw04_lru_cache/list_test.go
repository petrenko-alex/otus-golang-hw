package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// todo: can be cycle?
// todo: test for nil
// todo: test errors?
// todo: test for other types

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
}

func TestNewFilledList(t *testing.T) {
	tests := []struct {
		name        string
		elems       []interface{}
		expectedLen int
	}{
		{
			name:        "new empty list",
			elems:       []interface{}{},
			expectedLen: 0,
		},
		{
			name:        "new one element list",
			elems:       []interface{}{1},
			expectedLen: 1,
		},
		{
			name:        "new filled list",
			elems:       []interface{}{1, 2, 3, 4, 5},
			expectedLen: 5,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			list := NewFilledList(test.elems)

			require.Equal(t, test.expectedLen, list.Len())
			require.True(t, compareElements(test.elems, list))
			require.True(t, checkAddresses(list))
		})
	}
}

func compareElements(elems []interface{}, list List) bool {
	if len(elems) != list.Len() {
		return false
	}

	listItem := list.Front()
	if listItem != nil && len(elems) == 0 {
		return false
	}

	for _, val := range elems {
		if listItem != nil && listItem.Value != val {
			return false
		}

		listItem = listItem.Next
	}

	return true
}

func checkAddresses(l List) bool {
	if l.Len() == 0 && l.Front() != nil && l.Back() != nil {
		return false
	}

	if l.Front() != nil && l.Front().Prev != nil { // no prev element for head
		return false
	}

	if l.Back() != nil && l.Back().Next != nil { // no next element for tail
		return false
	}

	el := l.Front()
	for el != nil {
		if el.Next != nil && el != el.Next.Prev {
			return false
		}

		el = el.Next
	}

	return true
}

func TestLen(t *testing.T) {
	tests := []struct {
		name     string
		list     List
		expected int
	}{
		{
			name:     "len of empty list",
			list:     NewFilledList([]interface{}{}),
			expected: 0,
		},
		{
			name:     "len of one element list",
			list:     NewFilledList([]interface{}{1}),
			expected: 1,
		},
		{
			name:     "len of normal list",
			list:     NewFilledList([]interface{}{1, 2, 3}),
			expected: 3,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expected, test.list.Len())
		})
	}
}

func TestFront(t *testing.T) {
	t.Run("front from empty list", func(t *testing.T) {
		list := NewFilledList([]interface{}{})

		require.Nil(t, list.Front())
	})

	t.Run("front from one element list", func(t *testing.T) {
		list := NewFilledList([]interface{}{1})

		head := list.Front()

		require.NotNil(t, head)
		require.Equal(t, 1, head.Value)
		require.Equal(t, head, list.Back())
	})

	t.Run("front from normal list", func(t *testing.T) {
		list := NewFilledList([]interface{}{5, 4, 3, 2, 1})

		head := list.Front()

		require.NotNil(t, head)
		require.Equal(t, 5, head.Value)
	})
}

func TestBack(t *testing.T) {
	t.Run("back from empty list", func(t *testing.T) {
		list := NewFilledList([]interface{}{})

		require.Nil(t, list.Back())
	})

	t.Run("back from one element list", func(t *testing.T) {
		list := NewFilledList([]interface{}{1})

		tail := list.Back()

		require.NotNil(t, tail)
		require.Equal(t, 1, tail.Value)
		require.Equal(t, tail, list.Front())
	})

	t.Run("back from normal list", func(t *testing.T) {
		list := NewFilledList([]interface{}{1, 2, 3, 4, 5})

		tail := list.Back()

		require.NotNil(t, tail)
		require.Equal(t, 5, tail.Value)
	})
}

func TestPushFront(t *testing.T) {
	tests := []struct {
		name          string
		list          List
		pushValue     int
		expectedLen   int
		expectedFront int
		expectedBack  int
	}{
		{
			name:          "push front to empty list",
			list:          NewFilledList([]interface{}{}),
			pushValue:     1,
			expectedLen:   1,
			expectedFront: 1,
			expectedBack:  1,
		},
		{
			name:          "push front to normal list",
			list:          NewFilledList([]interface{}{1, 2, 3}),
			pushValue:     0,
			expectedLen:   4,
			expectedFront: 0,
			expectedBack:  3,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pushedEl := test.list.PushFront(test.pushValue)

			require.Equal(t, test.expectedLen, test.list.Len())
			require.Equal(t, test.expectedFront, pushedEl.Value)
			require.Equal(t, test.expectedBack, test.list.Back().Value)
			require.True(t, checkAddresses(test.list))
		})
	}
}

func TestPushBack(t *testing.T) {
	tests := []struct {
		name          string
		list          List
		pushValue     int
		expectedLen   int
		expectedFront int
		expectedBack  int
	}{
		{
			name:          "push back to empty list",
			list:          NewFilledList([]interface{}{}),
			pushValue:     1,
			expectedLen:   1,
			expectedFront: 1,
			expectedBack:  1,
		},
		{
			name:          "push back to normal list",
			list:          NewFilledList([]interface{}{1, 2, 3}),
			pushValue:     4,
			expectedLen:   4,
			expectedFront: 1,
			expectedBack:  4,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pushedEl := test.list.PushBack(test.pushValue)

			require.Equal(t, test.expectedLen, test.list.Len())
			require.Equal(t, test.expectedFront, test.list.Front().Value)
			require.Equal(t, test.expectedBack, pushedEl.Value)
			require.True(t, checkAddresses(test.list))
		})
	}
}

func TestRemove(t *testing.T) {
	t.Run("remove from normal list", func(t *testing.T) {
		list := NewFilledList([]interface{}{1, 2, 3})
		elementToRemove := list.Back()

		list.Remove(elementToRemove)

		require.Equal(t, 1, list.Front().Value)
		require.Equal(t, 2, list.Back().Value)
		require.True(t, checkAddresses(list))
	})

	t.Run("remove only element", func(t *testing.T) {
		list := NewFilledList([]interface{}{1})
		elementToRemove := list.Back()

		list.Remove(elementToRemove)

		require.Equal(t, 0, list.Len())
		require.Nil(t, list.Front())
		require.Nil(t, list.Back())
		require.True(t, checkAddresses(list))
	})

	t.Run("remove from empty list", func(t *testing.T) {
		list := NewFilledList([]interface{}{})
		elementToRemove := list.Back()

		list.Remove(elementToRemove)

		require.Equal(t, 0, list.Len())
		require.Nil(t, list.Front())
		require.Nil(t, list.Back())
		require.True(t, checkAddresses(list))
	})
}

func TestMoveToFront(t *testing.T) {
	t.Run("Move on empty list", func(t *testing.T) {
		list := NewFilledList([]interface{}{})
		elementToMove := list.Back()

		list.MoveToFront(elementToMove)

		require.Equal(t, 0, list.Len())
		require.Nil(t, list.Front())
		require.Nil(t, list.Back())
		require.True(t, checkAddresses(list))
	})
	t.Run("Move in one element list", func(t *testing.T) {
		list := NewFilledList([]interface{}{1})
		elementToMove := list.Back()

		list.MoveToFront(elementToMove)

		require.Equal(t, 1, list.Front().Value)
		require.Equal(t, 1, list.Back().Value)
		require.True(t, checkAddresses(list))
	})
	t.Run("Move from middle", func(t *testing.T) {
		list := NewFilledList([]interface{}{1, 2, 3})
		elementToMove := list.Front().Next

		list.MoveToFront(elementToMove)

		require.Equal(t, 2, list.Front().Value)
		require.Equal(t, 1, list.Front().Next.Value)
		require.Equal(t, 3, list.Back().Value)
		require.True(t, checkAddresses(list))
	})
	t.Run("Move from tail", func(t *testing.T) {
		list := NewFilledList([]interface{}{1, 2, 3})
		elementToMove := list.Back()

		list.MoveToFront(elementToMove)

		require.Equal(t, 3, list.Front().Value)
		require.Equal(t, 1, list.Front().Next.Value)
		require.Equal(t, 2, list.Back().Value)
		require.True(t, checkAddresses(list))
	})
}
