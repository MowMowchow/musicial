import React, { useEffect } from "react";
import { useDispatch } from "react-redux";
import { AppDispatch } from "../../redux/store";
import {
  Box,
  Button,
  Grommet,
  Heading,
  ResponsiveContext,
  Text,
} from "grommet";
import { fetchLoginRoute } from "../../utils/services/musicial-api";
import { deepMerge } from "grommet/utils";
import { css } from "styled-components";
import { gColours } from "../../App";
import Lottie from "react-lottie";
import loginAnimation from "../../utils/svgs/wave.json";
import { useCookies } from "react-cookie";

type LoginPageProps = {};

const LoginPage = ({}: LoginPageProps) => {
  const dispatch: AppDispatch = useDispatch();
  const [cookies, setCookie] = useCookies(["userId", "accessToken"]);

  useEffect(() => {
    dispatch({
      type: "app/setNavBarVisible",
      payload: false,
    });
  }, []);

  useEffect(() => {
    console.log("Login Polling Page | COOKIES: ", cookies);
  }, [cookies]);

  const LoginPollingPageBoxBorder = deepMerge({
    box: {
      extend: ({}) => css`
        border-radius: 1rem;
        box-shadow: 0px 0.12rem 0.25rem rgba(0, 0, 0, 0.2);
      `,
    },
    anchor: {
      color: gColours.eerieBlack,
      hover: {
        extend: ({}) => css`
          color: ${gColours.darkSalmon};
          text-decoration: none;
        `,
      },
    },
  });

  const lottieDefaultOptions = {
    loop: true,
    autoplay: true,
    animationData: loginAnimation,
    renderer: "svg",
  };

  return (
    <Grommet>
      <ResponsiveContext.Consumer>
        {(size) => (
          <Box direction="row" justify="center" height="100vh">
            <Box direction="column" justify="center">
              <Grommet theme={LoginPollingPageBoxBorder}>
                <Box
                  // height={size === "small" ? "17rem" : "medium"}
                  width={size === "small" ? "17rem" : "medium"}
                  align="center"
                  direction="column"
                  pad={{ horizontal: "small" }}
                  // just
                >
                  <Grommet>
                    <Box margin="small">
                      <Heading margin="0">Hey!</Heading>
                    </Box>
                    <Box>
                      <Box margin="small">
                        <Text size={size === "small" ? "small" : "medium"}>
                          Pressing the button below will update our databse with
                          your public Spotify information
                        </Text>{" "}
                      </Box>
                      <Box margin="small">
                        <Button
                          label="press to sign in with Spotify"
                          size={size === "small" ? "small" : "medium"}
                          color={gColours.byzantium}
                          onClick={() =>
                            (window.location.href = fetchLoginRoute())
                          }
                        ></Button>
                      </Box>
                      <Box margin="small">
                        <Lottie
                          options={lottieDefaultOptions}
                          height={size === "small" ? "7rem" : "12rem"}
                          width={size === "small" ? "7rem" : "12rem"}
                        ></Lottie>
                      </Box>
                    </Box>
                  </Grommet>
                </Box>
              </Grommet>
            </Box>
          </Box>
        )}
      </ResponsiveContext.Consumer>
    </Grommet>
  );
};

export default LoginPage;
