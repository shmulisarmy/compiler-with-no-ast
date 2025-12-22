// Package option provides a Rust-like Option[T] for Go.
package option

// Option represents either Some(T) or None.
type Option[T any] struct {
	v  T
	ok bool
}

// Some creates an Option with a value.
func Some[T any](v T) Option[T] {
	return Option[T]{v: v, ok: true}
}

// None creates an empty Option.
func None[T any]() Option[T] {
	var zero T
	return Option[T]{v: zero, ok: false}
}

// IsSome returns true if the option has a value.
func (o Option[T]) IsSome() bool {
	return o.ok
}

// IsNone returns true if the option is empty.
func (o Option[T]) IsNone() bool {
	return !o.ok
}

// Unwrap returns the value or panics if None.
func (o Option[T]) Unwrap() T {
	if !o.ok {
		panic("called Unwrap on None")
	}
	return o.v
}

// Expect returns the value or panics with msg.
func (o Option[T]) Expect(msg string) T {
	if !o.ok {
		panic(msg)
	}
	return o.v
}

// UnwrapOr returns the value or def.
func (o Option[T]) UnwrapOr(def T) T {
	if !o.ok {
		return def
	}
	return o.v
}

// UnwrapOrElse returns the value or computes it.
func (o Option[T]) UnwrapOrElse(f func() T) T {
	if !o.ok {
		return f()
	}
	return o.v
}
