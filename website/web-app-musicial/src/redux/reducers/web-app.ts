import { createAction, createReducer } from "@reduxjs/toolkit";
import { buildQueries } from "@testing-library/react";
import { RootState } from "../store";

export const initialAppState = {
  navBarVisible: true,
  userProcessed: false,
};

const setNavBarVisible = createAction<boolean>("app/setNavBarVisible");
const setUserProcessed = createAction<boolean>("app/setUserProcessed");

export const appReducer = createReducer(initialAppState, (builer) => {
  builer
    .addCase(setNavBarVisible, (state, action) => {
      state.navBarVisible = action.payload;
    })
    .addCase(setUserProcessed, (state, action) => {
      state.userProcessed = action.payload;
    });
});

export const getNavBarVisible = (state: RootState) => {
  return state.appReducer.navBarVisible;
};

export const getUserProcessed = (state: RootState) => {
  return state.appReducer.userProcessed;
};
