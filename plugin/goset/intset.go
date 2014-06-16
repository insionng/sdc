
package goset

import (
	"sort"
)

type IntSet struct {
	Map map[int]struct{}
}

func NewIntSet(values ...int) *IntSet{
	a := new(IntSet)
	a.Map = make(map[int]struct{},len(values))
	for _,v := range values{
		a.Map[v] = none
	}
	return a
}

func (a *IntSet)Add(elements ...int) *IntSet{
	for _,e := range elements {
		a.Map[e] = none
	}
	return a
}

func (a *IntSet)Remove(elements ...int) *IntSet{
	for _,e := range elements {
		delete(a.Map,e)
	}
	return a
}

func (a *IntSet)Contains(elements ...int)bool{
	for _,e := range elements {
		if _,ok := a.Map[e]; !ok {
			return false
		}
	}
	return true
}

func (a *IntSet)ContainsSet(b *IntSet) bool{
	for element, _ := range b.Map{
		if _, ok := a.Map[element]; !ok {
			return false
		}
	}
	return true
}

func (a *IntSet)Size() int{
	return len(a.Map)
}

func (a *IntSet)RemoveSet(b *IntSet) *IntSet{
	for element,_ := range b.Map{
		delete(a.Map,element)
	}
	return a
}

func (a *IntSet)AddSet(b *IntSet) *IntSet{
	for element, _ := range b.Map{
		a.Map[element] = none
	}
	return a
}

func (a *IntSet)InterSection(b *IntSet) *IntSet{
	for element, _ := range a.Map{
		if _,ok := b.Map[element]; !ok{
			delete(a.Map, element)
		}
	}
	return a
}

func (a *IntSet)ToSortedSlice() []int{
	slice := make([]int, len(a.Map))
	i := 0
	for element, _ := range a.Map{
		slice[i] = element
		i++
	}
	sort.Ints(slice)
	return slice
}
