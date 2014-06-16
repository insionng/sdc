
package goset

type InterfaceSet struct {
	Map map[interface {}]struct{}
}


func NewInterfaceSet(values ... interface {}) *InterfaceSet {
	a := new(InterfaceSet)
	a.Map = make(map[interface{}]struct{}, len(values))
	for _,v := range values{
		a.Map[v] = none
	}
	return a
}

func (a *InterfaceSet) Add(elements ...interface{}) *InterfaceSet {
	for _, e := range elements {
		a.Map[e] = none
	}
	return a
}

func (a *InterfaceSet) Remove(elements ...interface{}) *InterfaceSet {
	for _, e := range elements {
		delete(a.Map, e)
	}
	return a
}

func (a *InterfaceSet) Contains(elements ...interface{}) bool {
	for _, e := range elements {
		if _, ok := a.Map[e]; !ok{
			return false
		}
	}
	return true
}

func (a *InterfaceSet)ContainsSet(b *InterfaceSet) bool{
	for element, _ := range b.Map{
		if _, ok := a.Map[element]; !ok {
			return false
		}
	}
	return true
}


func (a *InterfaceSet) Size() int {
	return len(a.Map)
}

func (a *InterfaceSet) RemoveSet(b *InterfaceSet) *InterfaceSet {
	for element, _ := range b.Map {
		delete(a.Map, element)
	}
	return a
}

func (a *InterfaceSet) AddSet(b *InterfaceSet) *InterfaceSet {
	for element, _ := range b.Map {
		a.Map[element] = none
	}
	return a
}

func (a *InterfaceSet) InterSection(b *InterfaceSet) *InterfaceSet {
	for element, _ := range a.Map {
		if _, ok := b.Map[element]; !ok {
			delete(a.Map, element)
		}
	}
	return a
}


