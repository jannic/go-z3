package z3

import "unsafe"

// #include <stdlib.h>
// #include "go-z3.h"
import "C"

// Optimize is a single optimize instance tied to a specific Context within Z3.
//
// It is created via the NewOptimize methods on Context. When a optimize is
// no longer needed, the Close method must be called. This will remove the
// optimize from the context and no more APIs on Optimize may be called
// thereafter.
//
// Freeing the context (Context.Close) will NOT automatically close associated
// optimizers. They must be managed separately.
type Optimize struct {
	rawCtx    C.Z3_context
	rawOptimize C.Z3_optimize
}

// NewSolver creates a new solver.
func (c *Context) NewOptimize() *Optimize {
	rawOptimize := C.Z3_mk_optimize(c.raw)
	C.Z3_optimize_inc_ref(c.raw, rawOptimize)

	return &Optimize{
		rawOptimize: rawOptimize,
		rawCtx:    c.raw,
	}
}

// Close frees the memory associated with this.
func (s *Optimize) Close() error {
	C.Z3_optimize_dec_ref(s.rawCtx, s.rawOptimize)
	return nil
}

// Assert asserts a constraint onto the Solver.
//
// Maps to: Z3_solver_assert
func (s *Optimize) Assert(a *AST) {
	C.Z3_optimize_assert(s.rawCtx, s.rawOptimize, a.rawAST)
}

// Assert asserts a constraint onto the Solver.
//
// Maps to: Z3_solver_assert
func (s *Optimize) AssertSoft(a *AST, weight string, id *Symbol) {
	cweight := C.CString(weight)
	// TODO: make sure this is correct
	defer C.free(unsafe.Pointer(cweight))
	C.Z3_optimize_assert_soft(s.rawCtx, s.rawOptimize, a.rawAST, cweight, id.rawSymbol)
}

// Check checks if the currently set formula is consistent.
//
// Maps to: Z3_solver_check
func (s *Optimize) Check() LBool {
	return LBool(C.Z3_optimize_check(s.rawCtx, s.rawOptimize, 0, nil))
}

// Model returns the last model from a Check.
//
// Maps to: Z3_solver_get_model
func (s *Optimize) Model() *Model {
	m := &Model{
		rawCtx:   s.rawCtx,
		rawModel: C.Z3_optimize_get_model(s.rawCtx, s.rawOptimize),
	}
	m.IncRef()
	return m
}
