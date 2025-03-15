import {createSlice, Slice} from "@reduxjs/toolkit";
import {ICart} from "./ICart.ts";

interface INameQuantityMap {
    [key: string]: number;
}

interface IInitialCartState {
    items: ICart[];
    nameQuantityMap: INameQuantityMap;
}

const initialState: IInitialCartState = {
    items: [],
    nameQuantityMap: {}
}

export const cartSlice: Slice = createSlice({
    name: 'cart',
    initialState,
    reducers: {
        setItems: (state, action) => {
            state.items = action.payload;
        },
        cartItemAdded: (state, action) => {
            if (typeof state.nameQuantityMap[action.payload.name] === "undefined") {
                state.nameQuantityMap[action.payload.name] = 1;
            } else {
                state.nameQuantityMap[action.payload.name]++;
            }
            if (typeof state.nameQuantityMap[action.payload.name] === "undefined") {
                state.items.push(action.payload);
            }
        },
        cartItemRemoved: (state, action) => {
            for (let j = 0; j < state.items.length; j++) {
                if (state.items[j].id === action.payload.id) {
                    delete state.items[j];
                    state.nameQuantityMap[action.payload.name] = 0;
                    break;
                }
            }
        }
    }
});

export const {cartItemAdded, cartItemRemoved} = cartSlice.actions;

export default cartSlice.reducer;