import {
  Anchor,
  Box,
  Button,
  Grid,
  grommet,
  Header,
  Heading,
  Main,
  Nav,
  ResponsiveContext,
  TextInput,
} from "grommet";
import React, { useEffect, useState } from "react";
import { Link } from "react-router-dom";

import { useSelector } from "react-redux";
import { getNavBarVisible } from "../../redux/reducers/web-app";
import { getUserId, getUserAccessToken } from "../../redux/reducers/auth";
import { gColours } from "../../App";

type NavBarProps = {};

const NavBar = ({}: NavBarProps) => {
  const { isNavBarVisible } = useSelector((state: any) => {
    return {
      isNavBarVisible: getNavBarVisible(state),
    };
  });

  useEffect(() => {}, [isNavBarVisible]);

  return (
    <ResponsiveContext.Consumer>
      {(size) =>
        isNavBarVisible ? (
          <Box
            // width={"100%"}
            direction={size === "small" ? "column" : "row"}
            justify={size === "small" ? "around" : "between"}
            align={size === "small" ? "center" : ""}
            responsive={true}
            overflow="visible"
            flex={"grow"}
            margin="medium"
          >
            <Box
              direction="column"
              justify="center"
              margin={
                size === "small"
                  ? { left: "0", top: "0" }
                  : { left: "large", top: "0" }
              }
            >
              <Heading size="4.2rem" margin="0" color={gColours.byzantium}>
                Musicial
              </Heading>
            </Box>
            <Box direction="column" justify="center">
              <Nav direction="row" margin="medium" gap="small">
                <Heading
                  size="1.68rem"
                  margin="small"
                  color={gColours.byzantium}
                >
                  <Link
                    to="/"
                    style={{ textDecoration: "none", color: "inherit" }}
                  >
                    Home
                  </Link>
                </Heading>
                <Heading
                  size="1.68rem"
                  margin="small"
                  color={gColours.byzantium}
                >
                  <Link
                    to="/search"
                    style={{ textDecoration: "none", color: "inherit" }}
                  >
                    Search
                  </Link>
                </Heading>
                {/* <Heading size="1.68rem" margin="small" color={gColours.raisinBlack}>
                  <Link
                    to="/profile"
                    style={{ textDecoration: "none", color: "inherit" }}
                  >
                    Profile
                  </Link>
                </Heading> */}
              </Nav>
            </Box>
          </Box>
        ) : (
          <></>
        )
      }
    </ResponsiveContext.Consumer>
  );
};

export default NavBar;
