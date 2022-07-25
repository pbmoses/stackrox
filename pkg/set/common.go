package set

// TODO: remove these shortcuts to comonly used types, made to avoid huge diff.

// StringSet is a set of strings.
type StringSet = Set[string]

// FrozenStringSet is a frozen set of strings.
type FrozenStringSet = FrozenSet[string]

// NewStringSet creates an initialized set of strings.
func NewStringSet(initial ...string) Set[string] {
	return NewSet(initial...)
}

// NewFrozenStringSet creates an initialized frozen set of strings.
func NewFrozenStringSet(initial ...string) FrozenSet[string] {
	return NewFrozenSet(initial...)
}
