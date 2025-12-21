import {createSlice, Slice} from "@reduxjs/toolkit";
import {IUser} from "./types/IUser.ts";
import {IUserRole} from "./types/IUserRole.ts";

interface IInitialUserState {
    isSignedIn: boolean;
    user: IUser | null;
    roles: IUserRole[];
    level: number;
    experience: number;
}

const initialState: IInitialUserState = {
    isSignedIn: false,
    user: null,
    roles: [],
    level: 1,
    experience: 0,
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
        },
        setUserLevel: (state, action) => {
            state.level = action.payload;
        },
        setUserExperience: (state, action) => {
            state.experience = action.payload;
        }
    }
})

export const {setSignedIn, setUser, setSignedOut, setUserRoles, setUserExperience, setUserLevel} = userSlice.actions;
export default userSlice.reducer;