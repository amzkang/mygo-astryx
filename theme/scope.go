package theme

import "github.com/egoist/mygo/ui"

type scopeKey struct{}

// scopeStack keeps per-window frame-local *Theme overrides during a single build pass.
// Stored on c.Root() via ui.Local so no *ui.Context or ui.Element ever outlives a frame.
type scopeStack struct {
	stack []*Theme
}

func getScope(c *ui.Context) *scopeStack {
	return ui.Local(c.Root(), scopeKey{}, func() scopeStack {
		return scopeStack{}
	})
}

// WithTheme applies t to both mygo's native ui.Themed subtree and mygo-astryx's extended token stack.
// It is safe to nest multiple WithTheme calls in the same frame.
func WithTheme(c *ui.Context, t *Theme, fn func()) ui.Element {
	if t == nil {
		return ui.Column(c).Children(fn)
	}
	st := getScope(c)
	return ui.Themed(c, t.ToUI()).Children(func() {
		st.stack = append(st.stack, t)
		defer func() {
			st.stack = st.stack[:len(st.stack)-1]
		}()
		if fn != nil {
			fn()
		}
	})
}

// SetRoot sets t as the window's root theme for the remainder of the frame.
func SetRoot(c *ui.Context, t *Theme) {
	if t == nil {
		return
	}
	c.SetTheme(t.ToUI())
	st := getScope(c)
	st.stack = st.stack[:0]
	st.stack = append(st.stack, t)
}

// Of returns the active design system *Theme for the current build scope.
// If no WithTheme or SetRoot was called, it falls back to Neutral(c.Theme().Dark).
func Of(c *ui.Context) *Theme {
	st := getScope(c)
	if len(st.stack) > 0 {
		return st.stack[len(st.stack)-1]
	}
	return Neutral.For(c.Theme().Dark)
}
