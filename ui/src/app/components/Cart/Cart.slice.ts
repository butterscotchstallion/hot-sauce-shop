import {createSlice, Slice} from "@reduxjs/toolkit";
import {ICart} from "./ICart.ts";
import {getIdQuantityMap} from "./CartService.ts";

export interface IIDQuantityMap {
    [key: number]: number;
}

export interface IInitialCartState {
    items: ICart[];
    idQuantityMap: IIDQuantityMap;
}

const initialState: IInitialCartState = {
    items: [],
    idQuantityMap: {}
}

export const cartSlice: Slice = createSlice({
    name: 'cart',
    initialState,
    reducers: {
        setCartItems: (state, action) => {
            state.items = action.payload;
        },
        setIdQuantityMap: (state, action) => {
            state.idQuantityMap = getIdQuantityMap(action.payload);
        },
        cartItemAdded: (state, action) => {
            if (typeof state.idQuantityMap[action.payload.id] === "undefined") {
                state.idQuantityMap[action.payload.id] = 1;
            } else {
                state.idQuantityMap[action.payload.id]++;
                console.log("incremented " + action.payload.name + " to " + state.idQuantityMap[action.payload.id]);
            }
            if (typeof state.idQuantityMap[action.payload.id] === "undefined") {
                state.items.push(action.payload);
                console.log("pushed new item: " + action.payload.name);
            } else {
                for (let j = 0; j < state.items.length; j++) {
                    if (state.items[j].inventoryItemId === action.payload.id) {
                        state.items[j].quantity++;
                        break;
                    }
                }
            }
        },
        cartItemRemoved: (state, action) => {
            for (let j = 0; j < state.items.length; j++) {
                if (state.items[j].id === action.payload.id) {
                    delete state.items[j];
                    state.idQuantityMap[action.payload.id] = 0;
                    break;
                }
            }
        }
    }
});

export const {cartItemAdded, cartItemRemoved, setIdQuantityMap, setCartItems} = cartSlice.actions;
export default cartSlice.reducer;