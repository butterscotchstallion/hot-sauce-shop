import {createSlice, Slice} from "@reduxjs/toolkit";
import {ICart} from "./ICart.ts";
import {getIdQuantityMap} from "./CartService.ts";

export interface IIDQuantityMap {
    [key: number]: number;
}

export interface IInitialCartState {
    items: ICart[];
    idQuantityMap: IIDQuantityMap;
    cartSubtotal: number;
}

const initialState: IInitialCartState = {
    items: [],
    idQuantityMap: {},
    cartSubtotal: 0
}

export const cartSlice: Slice = createSlice({
    name: 'cart',
    initialState,
    reducers: {
        setCartSubtotal: (state, action) => {
            state.cartSubtotal = action.payload;
        },
        setCartItems: (state, action) => {
            state.items = action.payload;
        },
        setIdQuantityMap: (state, action) => {
            state.idQuantityMap = getIdQuantityMap(action.payload);
        },
        setCartItemQuantity: (state, action) => {
            // Update cart item quantity
            for (let j = 0; j < state.items.length; j++) {
                if (state.items[j].inventoryItemId === action.payload.id) {
                    state.items[j].quantity = action.payload.quantity;
                    break;
                }
            }
            // Update map
            for (const inventoryItemId in state.idQuantityMap) {
                if (Number(inventoryItemId) === action.payload.id) {
                    state.idQuantityMap[inventoryItemId] = action.payload.quantity;
                    break;
                }
            }
        },
        cartItemAdded: (state, action) => {
            if (typeof state.idQuantityMap[action.payload.id] === "undefined") {
                state.idQuantityMap[action.payload.id] = 1;
            } else {
                state.idQuantityMap[action.payload.id]++;
            }
            if (typeof state.idQuantityMap[action.payload.id] === "undefined") {
                state.items.push(action.payload);
            } else {
                for (let j = 0; j < state.items.length; j++) {
                    if (state.items[j].inventoryItemId === action.payload.id) {
                        state.items[j].quantity++;
                        console.info("Updated quantity of " + state.items[j].name + " to " + state.items[j].quantity);
                        break;
                    }
                }
            }
        },
        cartItemRemoved: (state, action) => {
            state.idQuantityMap[action.payload.id] = 0;
            state.items = state.items.filter((item: ICart) => {
                return item.inventoryItemId !== action.payload.id
            })
        }
    }
});

export const {
    cartItemAdded,
    cartItemRemoved,
    setCartItemQuantity,
    setIdQuantityMap,
    setCartItems,
    setCartSubtotal,
} = cartSlice.actions;
export default cartSlice.reducer;