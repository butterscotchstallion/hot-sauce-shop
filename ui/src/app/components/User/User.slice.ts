import {createSlice, Slice} from "@reduxjs/toolkit";
import {IUser} from "./IUser.ts";

interface IInitialUserState {
    isSignedIn: boolean;
    user: IUser | null;
}

const initialState: IInitialUserState = {
    isSignedIn: false,
    user: null
}

export const userSlice: Slice = createSlice({
    name: 'user',
    initialState,
    reducers: {
        setSignedIn: (state, action) => {
            state.isSignedIn = action.payload;
        },
        setUser: (state, action) => {
            state.user = action.payload;
        },
        setSignedOut: (state, _) => {
            state.user = null;
            state.isSignedIn = false;
        }
    }
})

export const {setSignedIn, setUser, setSignedOut} = userSlice.actions;
export default userSlice.reducer;