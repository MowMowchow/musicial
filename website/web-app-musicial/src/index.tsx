import React from "react";
import { createRoot } from "react-dom/client";
import { Provider } from "react-redux";
import { musicialStore } from "./redux/store";
import App from "./App";
import "./index.css";
import { CookiesProvider } from "react-cookie";

const rootElement = document.getElementById("root") as HTMLElement;
const root = createRoot(rootElement);

root.render(
  <CookiesProvider>
    <React.StrictMode>
      <Provider store={musicialStore}>
        <App />
      </Provider>
    </React.StrictMode>
  </CookiesProvider>
);
