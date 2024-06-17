import { configureStore, createSlice } from "@reduxjs/toolkit";

// User
const initialUserState = {
    username: "",
    first_name: "",
    isLoggedIn: false
};
const userSlice = createSlice({
  name: "user",
  initialState: initialUserState,
  reducers: {
    login: (state, action) => {
      state.username = action.payload.username;
      state.first_name = action.payload.first_name;
      state.isLoggedIn = action.payload.isLoggedIn;
    },

    logout: (state) => {
      state.username = initialUserState.username;
      state.first_name = initialUserState.first_name;
      state.isLoggedIn = initialUserState.isLoggedIn;
    },
  }
});

export const { login, logout } = userSlice.actions;

export const store = configureStore({
  reducer: {
    user: userSlice.reducer,
  },
});