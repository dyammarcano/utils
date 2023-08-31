package orderedmap

import (
	"container/list"
)

type KeyError struct {
	MissingKey any
}

type Pair struct {
	Key     any
	Value   any
	element *list.Element
}

type OrderedMap struct {
	pairs map[any]*Pair
	list  *list.List
}

var _ error = &KeyError{}

func (e *KeyError) Error() string {
	return "KeyError: Key not found"
}

func NewOrderedMap() *OrderedMap {
	return &OrderedMap{
		pairs: make(map[any]*Pair),
		list:  list.New(),
	}
}

func (om *OrderedMap) GetWithPair(key any) (*Pair, bool) {
	pair, ok := om.pairs[key]
	return pair, ok
}

func (om *OrderedMap) Get(key any) (any, bool) {
	pair, ok := om.GetWithPair(key)
	if ok {
		return pair.Value, true
	}
	return nil, false
}

func (om *OrderedMap) Load(key any) (any, bool) {
	return om.Get(key)
}

func (om *OrderedMap) GetPair(key any) *Pair {
	pair, _ := om.GetWithPair(key)
	return pair
}

func (om *OrderedMap) Set(key any, value any) (any, bool) {
	if pair, present := om.GetWithPair(key); present {
		oldValue := pair.Value
		pair.Value = value
		return oldValue, true
	}

	pair := &Pair{
		Key:   key,
		Value: value,
	}
	pair.element = om.list.PushBack(pair)
	om.pairs[key] = pair

	return nil, false
}

func (om *OrderedMap) Store(key any, value any) (any, bool) {
	return om.Set(key, value)
}

func (om *OrderedMap) Delete(key any) (any, bool) {
	if pair, present := om.GetWithPair(key); present {
		om.list.Remove(pair.element)
		delete(om.pairs, key)
		return pair.Value, true
	}
	return nil, false
}

func (om *OrderedMap) Len() int {
	return len(om.pairs)
}

func (om *OrderedMap) Oldest() *Pair {
	return listElementToPair(om.list.Front())
}

func (om *OrderedMap) Newest() *Pair {
	return listElementToPair(om.list.Back())
}

func (p *Pair) Next() *Pair {
	return listElementToPair(p.element.Next())
}

func (p *Pair) Prev() *Pair {
	return listElementToPair(p.element.Prev())
}

func listElementToPair(element *list.Element) *Pair {
	if element == nil {
		return nil
	}
	return element.Value.(*Pair)
}

func (om *OrderedMap) getElements(keys ...any) ([]*list.Element, error) {
	elements := make([]*list.Element, len(keys))
	for i, k := range keys {
		pair, present := om.GetWithPair(k)
		if !present {
			return nil, &KeyError{k}
		}
		elements[i] = pair.element
	}
	return elements, nil
}

// MoveAfter and MoveBefore functions factored to use helper function.
func (om *OrderedMap) MoveAfter(key, markKey any) error {
	return om.movePair(key, markKey, func(keyElement, markElement *list.Element) {
		om.list.MoveAfter(keyElement, markElement)
	})
}

func (om *OrderedMap) MoveBefore(key, markKey any) error {
	return om.movePair(key, markKey, func(keyElement, markElement *list.Element) {
		om.list.MoveBefore(keyElement, markElement)
	})
}

func (om *OrderedMap) movePair(key, markKey any, moveFunc func(keyElement, markElement *list.Element)) error {
	keyPair, keyOk := om.GetWithPair(key)
	markPair, markOk := om.GetWithPair(markKey)

	if !keyOk || !markOk {
		return &KeyError{}
	}

	moveFunc(keyPair.element, markPair.element)

	return nil
}

// MoveToBack gets error checked version.
func (om *OrderedMap) MoveToBack(key any) error {
	pair, present := om.GetWithPair(key)
	if !present {
		return &KeyError{key}
	}
	om.list.MoveToBack(pair.element)
	return nil
}

// MoveToFront gets error checked version.
func (om *OrderedMap) MoveToFront(key any) error {
	pair, present := om.GetWithPair(key)
	if !present {
		return &KeyError{key}
	}
	om.list.MoveToFront(pair.element)
	return nil
}
