package sets

import "sort"

var (
	exists = struct{}{}
)

type String struct {
	items map[string]struct{}
}

func NewString(args ...string) *String {
	set := &String{
		items: make(map[string]struct{}),
	}
	set.Add(args...)
	return set
}

func (s *String) Add(elements ...string) {
	for _, el := range elements {
		s.items[el] = exists
	}
}

func (s *String) Delete(element string) {
	delete(s.items, element)
}

// AsSlice returns the set as a slice ordered alphabetically
func (s *String) AsSlice() []string {
	sl := make([]string, 0, len(s.items))
	for k := range s.items {
		sl = append(sl, k)
	}
	sort.Strings(sl)
	return sl
}

func (s *String) Len() int {
	return len(s.items)
}
