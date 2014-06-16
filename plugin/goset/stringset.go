
package goset

import "sort"

var none = struct{}{}
type StringSet struct {
	Map map[string]struct{}
}


func NewStringSet(values ... string) *StringSet {
	a := new(StringSet)
	a.Map = make(map[string]struct{}, len(values))
	for _,v := range values{
		a.Map[v] = none
	}
	return a
}

func (a *StringSet) Add(elements ...string) *StringSet {
	for _, e := range elements {
		a.Map[e] = none
	}
	return a
}

func (a *StringSet) Remove(elements ...string) *StringSet {
	for _, e := range elements {
		delete(a.Map, e)
	}
	return a
}

func (a *StringSet) Contains(elements ...string) bool {
	for _, e := range elements {
		if _, ok := a.Map[e]; !ok{
			return false
		}
	}
	return true
}

func (a *StringSet)ContainsSet(b *StringSet) bool{
	for element, _ := range b.Map{
		if _, ok := a.Map[element]; !ok {
			return false
		}
	}
	return true
}


func (a *StringSet) Size() int {
	return len(a.Map)
}

func (a *StringSet) RemoveSet(b *StringSet) *StringSet {
	for element, _ := range b.Map {
		delete(a.Map, element)
	}
	return a
}

func (a *StringSet) AddSet(b *StringSet) *StringSet {
	for element, _ := range b.Map {
		a.Map[element] = none
	}
	return a
}

func (a *StringSet) InterSection(b *StringSet) *StringSet {
	for element, _ := range a.Map {
		if _, ok := b.Map[element]; !ok {
			delete(a.Map, element)
		}
	}
	return a
}

func (a *StringSet) ToSortedSlice() []string {
	slice := make([]string, len(a.Map))
	i := 0
	for element, _ := range a.Map {
		slice[i] = element
		i++
	}
	sort.Strings(slice)
	return slice
}

