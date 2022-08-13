package main

func dedup1(s []uint32) []uint32 {
	if len(s) < 2 {
		return s
	}
	tmp := make([]uint32, 0, len(s))

	for i := uint32(0); i < uint32(len(s)); i++ {
		// If current is not equal to next then store the current
		if s[i] != s[i+1] {
			tmp = append(tmp, s[i])
		}
	}

	// The last must be stored
	// Note that if it was repeated, the duplicates are NOT stored before
	tmp = append(tmp, s[len(s)-1])

	// Modify original slice
	s = nil
	s = append(s, tmp...)
	return s
}
