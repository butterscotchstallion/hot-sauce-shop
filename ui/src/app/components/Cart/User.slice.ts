import {createSlice, Slice} from "@reduxjs/toolkit";

interface IInitialUserState {
    isSignedIn: boolean;
}

const initialState: IInitialUserState = {
    isSignedIn: false
}

export const userSlice: Slice = createSlice({
    name: 'user',
    initialState,
    reducers: {
        signedIn: (state, action) => {
            state.user.isSignedIn = action.payload;
        }
    }
})