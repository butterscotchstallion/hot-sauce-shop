// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"hotsauceshop/ent/inventory"
	"hotsauceshop/ent/predicate"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// InventoryUpdate is the builder for updating Inventory entities.
type InventoryUpdate struct {
	config
	hooks    []Hook
	mutation *InventoryMutation
}

// Where appends a list predicates to the InventoryUpdate builder.
func (iu *InventoryUpdate) Where(ps ...predicate.Inventory) *InventoryUpdate {
	iu.mutation.Where(ps...)
	return iu
}

// SetName sets the "name" field.
func (iu *InventoryUpdate) SetName(s string) *InventoryUpdate {
	iu.mutation.SetName(s)
	return iu
}

// SetNillableName sets the "name" field if the given value is not nil.
func (iu *InventoryUpdate) SetNillableName(s *string) *InventoryUpdate {
	if s != nil {
		iu.SetName(*s)
	}
	return iu
}

// SetDescription sets the "description" field.
func (iu *InventoryUpdate) SetDescription(s string) *InventoryUpdate {
	iu.mutation.SetDescription(s)
	return iu
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (iu *InventoryUpdate) SetNillableDescription(s *string) *InventoryUpdate {
	if s != nil {
		iu.SetDescription(*s)
	}
	return iu
}

// SetShortDescription sets the "shortDescription" field.
func (iu *InventoryUpdate) SetShortDescription(s string) *InventoryUpdate {
	iu.mutation.SetShortDescription(s)
	return iu
}

// SetNillableShortDescription sets the "shortDescription" field if the given value is not nil.
func (iu *InventoryUpdate) SetNillableShortDescription(s *string) *InventoryUpdate {
	if s != nil {
		iu.SetShortDescription(*s)
	}
	return iu
}

// SetSlug sets the "slug" field.
func (iu *InventoryUpdate) SetSlug(s string) *InventoryUpdate {
	iu.mutation.SetSlug(s)
	return iu
}

// SetNillableSlug sets the "slug" field if the given value is not nil.
func (iu *InventoryUpdate) SetNillableSlug(s *string) *InventoryUpdate {
	if s != nil {
		iu.SetSlug(*s)
	}
	return iu
}

// SetPrice sets the "price" field.
func (iu *InventoryUpdate) SetPrice(f float32) *InventoryUpdate {
	iu.mutation.ResetPrice()
	iu.mutation.SetPrice(f)
	return iu
}

// SetNillablePrice sets the "price" field if the given value is not nil.
func (iu *InventoryUpdate) SetNillablePrice(f *float32) *InventoryUpdate {
	if f != nil {
		iu.SetPrice(*f)
	}
	return iu
}

// AddPrice adds f to the "price" field.
func (iu *InventoryUpdate) AddPrice(f float32) *InventoryUpdate {
	iu.mutation.AddPrice(f)
	return iu
}

// SetCreatedAt sets the "createdAt" field.
func (iu *InventoryUpdate) SetCreatedAt(t time.Time) *InventoryUpdate {
	iu.mutation.SetCreatedAt(t)
	return iu
}

// SetNillableCreatedAt sets the "createdAt" field if the given value is not nil.
func (iu *InventoryUpdate) SetNillableCreatedAt(t *time.Time) *InventoryUpdate {
	if t != nil {
		iu.SetCreatedAt(*t)
	}
	return iu
}

// SetUpdatedAt sets the "updatedAt" field.
func (iu *InventoryUpdate) SetUpdatedAt(t time.Time) *InventoryUpdate {
	iu.mutation.SetUpdatedAt(t)
	return iu
}

// ClearUpdatedAt clears the value of the "updatedAt" field.
func (iu *InventoryUpdate) ClearUpdatedAt() *InventoryUpdate {
	iu.mutation.ClearUpdatedAt()
	return iu
}

// Mutation returns the InventoryMutation object of the builder.
func (iu *InventoryUpdate) Mutation() *InventoryMutation {
	return iu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (iu *InventoryUpdate) Save(ctx context.Context) (int, error) {
	iu.defaults()
	return withHooks(ctx, iu.sqlSave, iu.mutation, iu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (iu *InventoryUpdate) SaveX(ctx context.Context) int {
	affected, err := iu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (iu *InventoryUpdate) Exec(ctx context.Context) error {
	_, err := iu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (iu *InventoryUpdate) ExecX(ctx context.Context) {
	if err := iu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (iu *InventoryUpdate) defaults() {
	if _, ok := iu.mutation.UpdatedAt(); !ok && !iu.mutation.UpdatedAtCleared() {
		v := inventory.UpdateDefaultUpdatedAt()
		iu.mutation.SetUpdatedAt(v)
	}
}

func (iu *InventoryUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(inventory.Table, inventory.Columns, sqlgraph.NewFieldSpec(inventory.FieldID, field.TypeInt))
	if ps := iu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := iu.mutation.Name(); ok {
		_spec.SetField(inventory.FieldName, field.TypeString, value)
	}
	if value, ok := iu.mutation.Description(); ok {
		_spec.SetField(inventory.FieldDescription, field.TypeString, value)
	}
	if value, ok := iu.mutation.ShortDescription(); ok {
		_spec.SetField(inventory.FieldShortDescription, field.TypeString, value)
	}
	if value, ok := iu.mutation.Slug(); ok {
		_spec.SetField(inventory.FieldSlug, field.TypeString, value)
	}
	if value, ok := iu.mutation.Price(); ok {
		_spec.SetField(inventory.FieldPrice, field.TypeFloat32, value)
	}
	if value, ok := iu.mutation.AddedPrice(); ok {
		_spec.AddField(inventory.FieldPrice, field.TypeFloat32, value)
	}
	if value, ok := iu.mutation.CreatedAt(); ok {
		_spec.SetField(inventory.FieldCreatedAt, field.TypeTime, value)
	}
	if value, ok := iu.mutation.UpdatedAt(); ok {
		_spec.SetField(inventory.FieldUpdatedAt, field.TypeTime, value)
	}
	if iu.mutation.UpdatedAtCleared() {
		_spec.ClearField(inventory.FieldUpdatedAt, field.TypeTime)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, iu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{inventory.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	iu.mutation.done = true
	return n, nil
}

// InventoryUpdateOne is the builder for updating a single Inventory entity.
type InventoryUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *InventoryMutation
}

// SetName sets the "name" field.
func (iuo *InventoryUpdateOne) SetName(s string) *InventoryUpdateOne {
	iuo.mutation.SetName(s)
	return iuo
}

// SetNillableName sets the "name" field if the given value is not nil.
func (iuo *InventoryUpdateOne) SetNillableName(s *string) *InventoryUpdateOne {
	if s != nil {
		iuo.SetName(*s)
	}
	return iuo
}

// SetDescription sets the "description" field.
func (iuo *InventoryUpdateOne) SetDescription(s string) *InventoryUpdateOne {
	iuo.mutation.SetDescription(s)
	return iuo
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (iuo *InventoryUpdateOne) SetNillableDescription(s *string) *InventoryUpdateOne {
	if s != nil {
		iuo.SetDescription(*s)
	}
	return iuo
}

// SetShortDescription sets the "shortDescription" field.
func (iuo *InventoryUpdateOne) SetShortDescription(s string) *InventoryUpdateOne {
	iuo.mutation.SetShortDescription(s)
	return iuo
}

// SetNillableShortDescription sets the "shortDescription" field if the given value is not nil.
func (iuo *InventoryUpdateOne) SetNillableShortDescription(s *string) *InventoryUpdateOne {
	if s != nil {
		iuo.SetShortDescription(*s)
	}
	return iuo
}

// SetSlug sets the "slug" field.
func (iuo *InventoryUpdateOne) SetSlug(s string) *InventoryUpdateOne {
	iuo.mutation.SetSlug(s)
	return iuo
}

// SetNillableSlug sets the "slug" field if the given value is not nil.
func (iuo *InventoryUpdateOne) SetNillableSlug(s *string) *InventoryUpdateOne {
	if s != nil {
		iuo.SetSlug(*s)
	}
	return iuo
}

// SetPrice sets the "price" field.
func (iuo *InventoryUpdateOne) SetPrice(f float32) *InventoryUpdateOne {
	iuo.mutation.ResetPrice()
	iuo.mutation.SetPrice(f)
	return iuo
}

// SetNillablePrice sets the "price" field if the given value is not nil.
func (iuo *InventoryUpdateOne) SetNillablePrice(f *float32) *InventoryUpdateOne {
	if f != nil {
		iuo.SetPrice(*f)
	}
	return iuo
}

// AddPrice adds f to the "price" field.
func (iuo *InventoryUpdateOne) AddPrice(f float32) *InventoryUpdateOne {
	iuo.mutation.AddPrice(f)
	return iuo
}

// SetCreatedAt sets the "createdAt" field.
func (iuo *InventoryUpdateOne) SetCreatedAt(t time.Time) *InventoryUpdateOne {
	iuo.mutation.SetCreatedAt(t)
	return iuo
}

// SetNillableCreatedAt sets the "createdAt" field if the given value is not nil.
func (iuo *InventoryUpdateOne) SetNillableCreatedAt(t *time.Time) *InventoryUpdateOne {
	if t != nil {
		iuo.SetCreatedAt(*t)
	}
	return iuo
}

// SetUpdatedAt sets the "updatedAt" field.
func (iuo *InventoryUpdateOne) SetUpdatedAt(t time.Time) *InventoryUpdateOne {
	iuo.mutation.SetUpdatedAt(t)
	return iuo
}

// ClearUpdatedAt clears the value of the "updatedAt" field.
func (iuo *InventoryUpdateOne) ClearUpdatedAt() *InventoryUpdateOne {
	iuo.mutation.ClearUpdatedAt()
	return iuo
}

// Mutation returns the InventoryMutation object of the builder.
func (iuo *InventoryUpdateOne) Mutation() *InventoryMutation {
	return iuo.mutation
}

// Where appends a list predicates to the InventoryUpdate builder.
func (iuo *InventoryUpdateOne) Where(ps ...predicate.Inventory) *InventoryUpdateOne {
	iuo.mutation.Where(ps...)
	return iuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (iuo *InventoryUpdateOne) Select(field string, fields ...string) *InventoryUpdateOne {
	iuo.fields = append([]string{field}, fields...)
	return iuo
}

// Save executes the query and returns the updated Inventory entity.
func (iuo *InventoryUpdateOne) Save(ctx context.Context) (*Inventory, error) {
	iuo.defaults()
	return withHooks(ctx, iuo.sqlSave, iuo.mutation, iuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (iuo *InventoryUpdateOne) SaveX(ctx context.Context) *Inventory {
	node, err := iuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (iuo *InventoryUpdateOne) Exec(ctx context.Context) error {
	_, err := iuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (iuo *InventoryUpdateOne) ExecX(ctx context.Context) {
	if err := iuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (iuo *InventoryUpdateOne) defaults() {
	if _, ok := iuo.mutation.UpdatedAt(); !ok && !iuo.mutation.UpdatedAtCleared() {
		v := inventory.UpdateDefaultUpdatedAt()
		iuo.mutation.SetUpdatedAt(v)
	}
}

func (iuo *InventoryUpdateOne) sqlSave(ctx context.Context) (_node *Inventory, err error) {
	_spec := sqlgraph.NewUpdateSpec(inventory.Table, inventory.Columns, sqlgraph.NewFieldSpec(inventory.FieldID, field.TypeInt))
	id, ok := iuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Inventory.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := iuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, inventory.FieldID)
		for _, f := range fields {
			if !inventory.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != inventory.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := iuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := iuo.mutation.Name(); ok {
		_spec.SetField(inventory.FieldName, field.TypeString, value)
	}
	if value, ok := iuo.mutation.Description(); ok {
		_spec.SetField(inventory.FieldDescription, field.TypeString, value)
	}
	if value, ok := iuo.mutation.ShortDescription(); ok {
		_spec.SetField(inventory.FieldShortDescription, field.TypeString, value)
	}
	if value, ok := iuo.mutation.Slug(); ok {
		_spec.SetField(inventory.FieldSlug, field.TypeString, value)
	}
	if value, ok := iuo.mutation.Price(); ok {
		_spec.SetField(inventory.FieldPrice, field.TypeFloat32, value)
	}
	if value, ok := iuo.mutation.AddedPrice(); ok {
		_spec.AddField(inventory.FieldPrice, field.TypeFloat32, value)
	}
	if value, ok := iuo.mutation.CreatedAt(); ok {
		_spec.SetField(inventory.FieldCreatedAt, field.TypeTime, value)
	}
	if value, ok := iuo.mutation.UpdatedAt(); ok {
		_spec.SetField(inventory.FieldUpdatedAt, field.TypeTime, value)
	}
	if iuo.mutation.UpdatedAtCleared() {
		_spec.ClearField(inventory.FieldUpdatedAt, field.TypeTime)
	}
	_node = &Inventory{config: iuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, iuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{inventory.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	iuo.mutation.done = true
	return _node, nil
}
