import {createSlice, Slice} from "@reduxjs/toolkit";
import {ICart} from "./ICart.ts";
import {getIdQuantityMap} from "./CartService.ts";

interface IIDQuantityMap {
    [key: number]: number;
}

interface IInitialCartState {
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
            console.log(state.idQuantityMap);
        },
        cartItemAdded: (state, action) => {
            if (typeof state.idQuantityMap[action.payload.id] === "undefined") {
                state.idQuantityMap[action.payload.id] = 1;
            } else {
                state.idQuantityMap[action.payload.id]++;
            }
            if (typeof state.idQuantityMap[action.payload.id] === "undefined") {
                state.items.push(action.payload);
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