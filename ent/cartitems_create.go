// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"hotsauceshop/ent/cartitems"
	"hotsauceshop/ent/inventory"
	"hotsauceshop/ent/user"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// CartItemsCreate is the builder for creating a CartItems entity.
type CartItemsCreate struct {
	config
	mutation *CartItemsMutation
	hooks    []Hook
}

// SetQuantity sets the "quantity" field.
func (cic *CartItemsCreate) SetQuantity(i int8) *CartItemsCreate {
	cic.mutation.SetQuantity(i)
	return cic
}

// SetNillableQuantity sets the "quantity" field if the given value is not nil.
func (cic *CartItemsCreate) SetNillableQuantity(i *int8) *CartItemsCreate {
	if i != nil {
		cic.SetQuantity(*i)
	}
	return cic
}

// SetCreatedAt sets the "createdAt" field.
func (cic *CartItemsCreate) SetCreatedAt(t time.Time) *CartItemsCreate {
	cic.mutation.SetCreatedAt(t)
	return cic
}

// SetNillableCreatedAt sets the "createdAt" field if the given value is not nil.
func (cic *CartItemsCreate) SetNillableCreatedAt(t *time.Time) *CartItemsCreate {
	if t != nil {
		cic.SetCreatedAt(*t)
	}
	return cic
}

// SetUpdatedAt sets the "updatedAt" field.
func (cic *CartItemsCreate) SetUpdatedAt(t time.Time) *CartItemsCreate {
	cic.mutation.SetUpdatedAt(t)
	return cic
}

// SetNillableUpdatedAt sets the "updatedAt" field if the given value is not nil.
func (cic *CartItemsCreate) SetNillableUpdatedAt(t *time.Time) *CartItemsCreate {
	if t != nil {
		cic.SetUpdatedAt(*t)
	}
	return cic
}

// AddUserIDs adds the "user" edge to the User entity by IDs.
func (cic *CartItemsCreate) AddUserIDs(ids ...int) *CartItemsCreate {
	cic.mutation.AddUserIDs(ids...)
	return cic
}

// AddUser adds the "user" edges to the User entity.
func (cic *CartItemsCreate) AddUser(u ...*User) *CartItemsCreate {
	ids := make([]int, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return cic.AddUserIDs(ids...)
}

// AddInventoryIDs adds the "inventory" edge to the Inventory entity by IDs.
func (cic *CartItemsCreate) AddInventoryIDs(ids ...int) *CartItemsCreate {
	cic.mutation.AddInventoryIDs(ids...)
	return cic
}

// AddInventory adds the "inventory" edges to the Inventory entity.
func (cic *CartItemsCreate) AddInventory(i ...*Inventory) *CartItemsCreate {
	ids := make([]int, len(i))
	for j := range i {
		ids[j] = i[j].ID
	}
	return cic.AddInventoryIDs(ids...)
}

// Mutation returns the CartItemsMutation object of the builder.
func (cic *CartItemsCreate) Mutation() *CartItemsMutation {
	return cic.mutation
}

// Save creates the CartItems in the database.
func (cic *CartItemsCreate) Save(ctx context.Context) (*CartItems, error) {
	cic.defaults()
	return withHooks(ctx, cic.sqlSave, cic.mutation, cic.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (cic *CartItemsCreate) SaveX(ctx context.Context) *CartItems {
	v, err := cic.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (cic *CartItemsCreate) Exec(ctx context.Context) error {
	_, err := cic.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cic *CartItemsCreate) ExecX(ctx context.Context) {
	if err := cic.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (cic *CartItemsCreate) defaults() {
	if _, ok := cic.mutation.Quantity(); !ok {
		v := cartitems.DefaultQuantity
		cic.mutation.SetQuantity(v)
	}
	if _, ok := cic.mutation.CreatedAt(); !ok {
		v := cartitems.DefaultCreatedAt()
		cic.mutation.SetCreatedAt(v)
	}
	if _, ok := cic.mutation.UpdatedAt(); !ok {
		v := cartitems.DefaultUpdatedAt()
		cic.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (cic *CartItemsCreate) check() error {
	if _, ok := cic.mutation.Quantity(); !ok {
		return &ValidationError{Name: "quantity", err: errors.New(`ent: missing required field "CartItems.quantity"`)}
	}
	if _, ok := cic.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "createdAt", err: errors.New(`ent: missing required field "CartItems.createdAt"`)}
	}
	return nil
}

func (cic *CartItemsCreate) sqlSave(ctx context.Context) (*CartItems, error) {
	if err := cic.check(); err != nil {
		return nil, err
	}
	_node, _spec := cic.createSpec()
	if err := sqlgraph.CreateNode(ctx, cic.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	cic.mutation.id = &_node.ID
	cic.mutation.done = true
	return _node, nil
}

func (cic *CartItemsCreate) createSpec() (*CartItems, *sqlgraph.CreateSpec) {
	var (
		_node = &CartItems{config: cic.config}
		_spec = sqlgraph.NewCreateSpec(cartitems.Table, sqlgraph.NewFieldSpec(cartitems.FieldID, field.TypeInt))
	)
	if value, ok := cic.mutation.Quantity(); ok {
		_spec.SetField(cartitems.FieldQuantity, field.TypeInt8, value)
		_node.Quantity = value
	}
	if value, ok := cic.mutation.CreatedAt(); ok {
		_spec.SetField(cartitems.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := cic.mutation.UpdatedAt(); ok {
		_spec.SetField(cartitems.FieldUpdatedAt, field.TypeTime, value)
		_node.UpdatedAt = &value
	}
	if nodes := cic.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   cartitems.UserTable,
			Columns: []string{cartitems.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := cic.mutation.InventoryIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   cartitems.InventoryTable,
			Columns: cartitems.InventoryPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(inventory.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// CartItemsCreateBulk is the builder for creating many CartItems entities in bulk.
type CartItemsCreateBulk struct {
	config
	err      error
	builders []*CartItemsCreate
}

// Save creates the CartItems entities in the database.
func (cicb *CartItemsCreateBulk) Save(ctx context.Context) ([]*CartItems, error) {
	if cicb.err != nil {
		return nil, cicb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(cicb.builders))
	nodes := make([]*CartItems, len(cicb.builders))
	mutators := make([]Mutator, len(cicb.builders))
	for i := range cicb.builders {
		func(i int, root context.Context) {
			builder := cicb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*CartItemsMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, cicb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, cicb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, cicb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (cicb *CartItemsCreateBulk) SaveX(ctx context.Context) []*CartItems {
	v, err := cicb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (cicb *CartItemsCreateBulk) Exec(ctx context.Context) error {
	_, err := cicb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cicb *CartItemsCreateBulk) ExecX(ctx context.Context) {
	if err := cicb.Exec(ctx); err != nil {
		panic(err)
	}
}
