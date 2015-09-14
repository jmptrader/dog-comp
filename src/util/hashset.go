package util

type HashSet struct {
	set map[interface{}]bool
}

func HashSet_new() *HashSet {
	s := new(HashSet)
	s.set = make(map[interface{}]bool)
	return s
}

func (this *HashSet) Add(k interface{}) bool {
	if _, ok := this.set[k]; ok {
		return false
	} else {
		this.set[k] = true
		return true
	}
}

func (this *HashSet) Remove(k interface{}) bool {
	if _, ok := this.set[k]; ok {
		delete(this.set, k)
		return true
	} else {
		return false
	}
}

func (this *HashSet) Contains(k interface{}) bool {
	return this.set[k]
}

func (this *HashSet) Size() int {
	return len(this.set)
}

func (this *HashSet) Clear() {
	this.set = make(map[interface{}]bool)
}
