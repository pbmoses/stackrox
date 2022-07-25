package set

// TODO: remove these shortcuts to comonly used types, made to avoid huge diff.
type StringSet = Set[string]
type FrozenStringSet = FrozenSet[string]

func NewStringSet(initial ...string) Set[string] {
	return NewSet(initial...)
}

func NewFrozenStringSet(initial ...string) FrozenSet[string] {
	return NewFrozenSet(initial...)
}
