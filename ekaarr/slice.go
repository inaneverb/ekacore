package ekaarr

func Reduce[T any, R any](
	in []T, out R, cb func(acc R, value T, index int, arr []T) R) R {

	for i, n := 0, len(in); i < n; i++ {
		out = cb(out, in[i], i, in)
	}

	return out
}

func Filter[T any](in []T, cb func(T) bool) []T {

	// Copy on write.
	// Meaning, there's no copy, if cb returns true for all elements.

	var i, n, t = 0, len(in), true
	for t = true; i < n && t; i++ {
		t = cb(in[i])
		// Don't forget i is incremented even when t set to false.
	}

	if !t {
		i--
	}
	if i == n || i == n-1 {
		return in[:i] // Nothing to filter.
	}

	// Maybe remained elements will be filtered all?
	// If so, there's no need to make a copy so early.

	var j = i
	for t = true; i < n && t; i++ {
		t = !cb(in[i])
	}

	if !t {
		i--
	}
	if i == n {
		return in[:j] // Leave only not filtered.
	}

	// Copy those which were not filtered.

	var out = make([]T, 0, n) // <-- Make a copy.
	for i := 0; i < j; i++ {
		out = append(out, in[i])
	}

	for ; i < n; i++ {
		if cb(in[i]) {
			out = append(out, in[i])
		}
	}

	return out
}

func Unique[T comparable](in []T) []T {
	return unique(in, false)
}

func Distinct[T comparable](in []T) []T {
	return unique(in, true)
}

func Remove[T comparable](in []T, remove ...T) []T {
	return Filter(in, func(v T) bool { return !ContainsAny(remove, v) })
}

func ContainsAny[T comparable](in []T, search ...T) bool {
	return contains(in, search, false)
}

func ContainsAll[T comparable](in []T, search ...T) bool {
	return contains(in, search, true)
}

////////////////////////////////////////////////////////////////////////////////
///// PRIVATE METHODS //////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func contains[T comparable](in, search []T, allRequired bool) bool {

	var found, stop bool

	for i, n := 0, len(search); i < n && !stop; i++ {
		found = false
		for j, m := 0, len(in); j < m && !found; j++ {
			found = search[i] == in[j]
		}
		stop = found != allRequired // DO NOT TRY TO SIMPLIFY THIS!!
	}

	return found
}

func unique[T comparable](in []T, includeOnce bool) []T {

	var n = len(in)

	for i := 0; i < n-1; i++ {
		var duplicate = false

		for j := i + 1; j < n; j++ {
			if in[i] == in[j] {
				in[j] = in[n-1]
				n--
				j--
				duplicate = true
			}
		}

		if !includeOnce && duplicate {
			in[i] = in[n-1]
			n--
			i--
		}
	}

	return in[:n]
}
