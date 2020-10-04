package cidrreverse

import (
	"sort"
)

type CIDRv4Set []CIDRv4

func (s *CIDRv4Set) Plus(c CIDRv4) {
	insertIndex := len(*s)

	// find if any exists cidr contains this cidr
	for i, ec := range *s {
		if ec.Contains(c) {
			// the added cidr is already in the set, nothing to do
			return
		}
		if ec.IP >= c.IP {
			insertIndex = i
			break
		}
	}

	removeCount := 0

	// find if any exists subnet covered by the added cidr
	for i := insertIndex; i < len(*s); i++ {
		if c.Contains((*s)[i]) {
			removeCount++
		} else {
			break
		}
	}

	s.delete(insertIndex, removeCount)

	s.insert(insertIndex, c)

	// check mergeable
	for {
		cur := (*s)[insertIndex]
		var ci1, ci2 int
		if cur.IP&IPv4(1<<(32-cur.Mask)) == 0 {
			ci1 = insertIndex
			ci2 = insertIndex + 1
			if ci2 >= len(*s) {
				break
			}
		} else {
			ci1 = insertIndex - 1
			ci2 = insertIndex
			if ci1 < 0 {
				break
			}
		}
		c1 := (*s)[ci1]
		c2 := (*s)[ci2]
		if ok, res := c1.Merge(c2); ok {
			s.delete(ci1, 2)
			s.insert(ci1, res)
			insertIndex = ci1
		} else {
			break
		}
	}
}

func (s *CIDRv4Set) Minus(c CIDRv4) {
	containsIndex := -1

	for i, ec := range *s {
		if ec.Contains(c) {
			containsIndex = i
		}
	}

	if containsIndex == -1 {
		// parent not found, try to check and clear subnet
		insertIndex := s.findInsertIndex(c.IP)
		removeCount := 0

		for i := insertIndex; i < len(*s); i++ {
			if c.Contains((*s)[i]) {
				removeCount++
			} else {
				break
			}
		}

		s.delete(insertIndex, removeCount)

		return
	}

	// parent found, split the parent
	container := (*s)[containsIndex]
	s.delete(containsIndex, 1)

	if c == container {
		// the minus cidr has been removed
		return
	}

	// generate split cidr
	var remains []CIDRv4

	for m := container.Mask + 1; m <= c.Mask; m++ {
		nc := CIDRv4{
			IP:   c.IP^(1 << (32 - m)),
			Mask: m,
		}
		nc.Standardize()
		remains = append(remains, nc)
	}

	sort.Slice(remains, func(i, j int) bool {
		return remains[i].IP < remains[j].IP
	})

	s.insert(containsIndex, remains...)
}

func (s *CIDRv4Set) Empty() bool {
	return len(*s) == 0
}

func (s *CIDRv4Set) Clear() {
	*s = nil
}

func (s *CIDRv4Set) insert(index int, data ...CIDRv4) {
	if len(data) == 0 {
		return
	}
	*s = append((*s)[:index], append(data, (*s)[index:]...)...)
}

func (s *CIDRv4Set) delete(index int, count int) {
	if count == 0 {
		return
	}
	*s = append((*s)[:index], (*s)[index+count:]...)
}

func (s *CIDRv4Set) findInsertIndex(ip IPv4) int {
	return sort.Search(len(*s), func(i int) bool {
		return (*s)[i].IP >= ip
	})
}
