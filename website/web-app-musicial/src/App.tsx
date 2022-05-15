import { Box, Grommet, ResponsiveContext } from "grommet";
import React from "react";
import { css } from "styled-components";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import "./App.css";
import NavBar from "./components/navbar/NavBar";
import HomePage from "./pages/home/HomePage";
import LoginPage from "./pages/login/LoginPage";
import ProfilePage from "./pages/profile/ProfilePage";
import SearchPage from "./pages/search/SearchPage";
import LoginPollingPage from "./pages/loginPolling/LoginPollingPage";
import { deepMerge } from "grommet/utils";

const theme = {
  // plain: true,
  global: {
    font: {
      family: "Jost, sans-serif",
      // family: "Roboto",
      // size: "18px",
      // height: "20px",
    },
    colors: {
      ivory: "#FFFFF1",
      white: "#FAFAFA",
      cardWhite: "#fbfbfb",
      babyPowder: "#FFFFF8",
      byzantium: "#79025B",
      unSelectedByzantium: "#965287",
      turquoise: "#41D3BD",
      fieryRose: "#FF5964",
      eerieBlack: "#221F1C",
      darkSalmon: "#E49273",
      ming: "#006D77",
      kobi: "#D2A1B8",
      cultured: "#F5F5F5",
      palePink: "#E8C2CA",
      thistle: "#D1B3C4",
      glossyGrape: "#B392AC",
      oldLavender: "#735D78",
      raisinBlack: "#2C252D",
    },
  },
};

export const gColours = theme.global.colors;

function App() {
  return (
    <Grommet theme={theme} background={gColours.babyPowder}>
      <Box width={{ min: "100vh" }} height={{ min: "100vh" }} fill={true}>
        <Box>
          <ResponsiveContext.Consumer>
            {(size) => (
              <BrowserRouter>
                <NavBar />
                <Routes>
                  <Route path="/" element={<HomePage />} />
                  <Route path="/profile" element={<ProfilePage />} />
                  <Route path="/search" element={<SearchPage />} />
                  <Route path="/login" element={<LoginPage />} />
                  <Route path="/pollLogin" element={<LoginPollingPage />} />
                </Routes>
              </BrowserRouter>
            )}
          </ResponsiveContext.Consumer>
        </Box>
      </Box>
    </Grommet>
  );
}

export default App;
