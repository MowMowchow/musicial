import { configureStore } from "@reduxjs/toolkit";

import { initialAppState } from "./reducers/web-app";
import { appReducer } from "./reducers/web-app";

import { initialAuthState } from "./reducers/auth";
import { authReducer } from "./reducers/auth";

const initialState = {
  appReducer: initialAppState,
  authReducer: initialAuthState,
};

const reducer = {
  appReducer: appReducer,
  authReducer: authReducer,
};

export const musicialStore = configureStore({
  reducer: reducer,
  devTools: true,
  preloadedState: initialState,
});

export type RootState = ReturnType<typeof musicialStore.getState>;
export type AppDispatch = typeof musicialStore.dispatch;
