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
        setSignedIn: (state, action) => {
            state.isSignedIn = action.payload;
        }
    }
})

export const {setSignedIn} = userSlice.actions;
export default userSlice.reducer;