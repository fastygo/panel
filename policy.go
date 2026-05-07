package panel

// Principal is the minimal actor contract used by panel policies.
type Principal[C comparable] interface {
	Has(C) bool
}

// Policy can make contextual authorization decisions for panel operations.
type Policy[P Principal[C], C comparable] interface {
	Can(P, C) bool
}

// PolicyFunc adapts a function into a Policy.
type PolicyFunc[P Principal[C], C comparable] func(P, C) bool

// Can evaluates the policy function.
func (f PolicyFunc[P, C]) Can(principal P, capability C) bool {
	return f(principal, capability)
}
