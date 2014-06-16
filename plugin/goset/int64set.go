
package goset

import "sort"

type Int64Slice []int64

var _ *sort.Interface = (*sort.Interface)(nil)

func (p Int64Slice) Len() int           { return len(p) }
func (p Int64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Int64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type Int64Set struct {
	Map map[int64]struct{}
}

func NewInt64Set(values ...int64) *Int64Set{
	a := new(Int64Set)
	a.Map = make(map[int64]struct{},len(values))
	for _,v := range values{
		a.Map[v] = none
	}
	return a
}

func (a *Int64Set)Add(elements ...int64) *Int64Set{
	for _,e := range elements {
		a.Map[e] = none
	}
	return a
}

func (a *Int64Set)Remove(elements ...int64) *Int64Set{
	for _,e := range elements {
		delete(a.Map,e)
	}
	return a
}

func (a *Int64Set)Contains(elements ...int64)bool{
	for _,e := range elements {
		if _,ok := a.Map[e]; !ok {
			return false
		}
	}
	return true
}

func (a *Int64Set)ContainsSet(b *Int64Set) bool{
	for element, _ := range b.Map{
		if _, ok := a.Map[element]; !ok {
			return false
		}
	}
	return true
}

func (a *Int64Set)Size() int{
	return len(a.Map)
}

func (a *Int64Set)RemoveSet(b *Int64Set) *Int64Set{
	for element,_ := range b.Map{
		delete(a.Map,element)
	}
	return a
}

func (a *Int64Set)AddSet(b *Int64Set) *Int64Set{
	for element, _ := range b.Map{
		a.Map[element] = none
	}
	return a
}

func (a *Int64Set)InterSection(b *Int64Set) *Int64Set{
	for element, _ := range a.Map{
		if _,ok := b.Map[element]; !ok{
			delete(a.Map, element)
		}
	}
	return a
}

func (a *Int64Set)ToSortedSlice() []int64{
	slice := make([]int64, len(a.Map))
	i := 0
	for element, _ := range a.Map{
		slice[i] = element
		i++
	}
	sort.Sort(Int64Slice(slice))
	return slice
}
