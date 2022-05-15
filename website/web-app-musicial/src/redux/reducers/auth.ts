import { createAction, createReducer } from "@reduxjs/toolkit";
import { RootState } from "../store";

export const initialAuthState = {
  userId: "",
  accessToken: "",
};

const setUserId = createAction<string>("auth/setUserId");
const setToken = createAction<string>("auth/setToken");

export const authReducer = createReducer(initialAuthState, (builer) => {
  builer
    .addCase(setUserId, (state, action) => {
      state.userId = action.payload;
    })
    .addCase(setToken, (state, action) => {
      state.accessToken = action.payload;
    });
});

export const getUserId = (state: RootState) => {
  return state.authReducer.userId;
};

export const getUserAccessToken = (state: RootState) => {
  return state.authReducer.accessToken;
};
