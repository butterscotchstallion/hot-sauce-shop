import {createSlice, Slice} from "@reduxjs/toolkit";
import {IUser} from "./IUser.ts";
import {IUserRole} from "./IUserRole.ts";

interface IInitialUserState {
    isSignedIn: boolean;
    user: IUser | null;
    roles: IUserRole[];
}

const initialState: IInitialUserState = {
    isSignedIn: false,
    user: null,
    roles: []
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
        },
        setUserRoles: (state, action) => {
            state.roles = action.payload;
        }
    }
})

export const {setSignedIn, setUser, setSignedOut, setUserRoles} = userSlice.actions;
export default userSlice.reducer;