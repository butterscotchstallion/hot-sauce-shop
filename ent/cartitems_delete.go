// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"hotsauceshop/ent/cartitems"
	"hotsauceshop/ent/predicate"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// CartItemsDelete is the builder for deleting a CartItems entity.
type CartItemsDelete struct {
	config
	hooks    []Hook
	mutation *CartItemsMutation
}

// Where appends a list predicates to the CartItemsDelete builder.
func (cid *CartItemsDelete) Where(ps ...predicate.CartItems) *CartItemsDelete {
	cid.mutation.Where(ps...)
	return cid
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (cid *CartItemsDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, cid.sqlExec, cid.mutation, cid.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (cid *CartItemsDelete) ExecX(ctx context.Context) int {
	n, err := cid.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (cid *CartItemsDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(cartitems.Table, sqlgraph.NewFieldSpec(cartitems.FieldID, field.TypeInt))
	if ps := cid.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, cid.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	cid.mutation.done = true
	return affected, err
}

// CartItemsDeleteOne is the builder for deleting a single CartItems entity.
type CartItemsDeleteOne struct {
	cid *CartItemsDelete
}

// Where appends a list predicates to the CartItemsDelete builder.
func (cido *CartItemsDeleteOne) Where(ps ...predicate.CartItems) *CartItemsDeleteOne {
	cido.cid.mutation.Where(ps...)
	return cido
}

// Exec executes the deletion query.
func (cido *CartItemsDeleteOne) Exec(ctx context.Context) error {
	n, err := cido.cid.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{cartitems.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (cido *CartItemsDeleteOne) ExecX(ctx context.Context) {
	if err := cido.Exec(ctx); err != nil {
		panic(err)
	}
}
