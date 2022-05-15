import { useEffect, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";

import { useDispatch } from "react-redux";
import { AppDispatch } from "../../redux/store";
import { useSelector } from "react-redux";
import { getNavBarVisible } from "../../redux/reducers/web-app";
import { getUserId, getUserAccessToken } from "../../redux/reducers/auth";
import { getUserProcessed } from "../../redux/reducers/web-app";

import { Box, Grommet, Heading, ResponsiveContext, Text } from "grommet";
import { deepMerge } from "grommet/utils";
import { css } from "styled-components";
import Lottie from "react-lottie";
import { gColours } from "../../App";
import { pollLoginProcess } from "../../utils/services/musicial-api";
import loadingAnimation from "../../utils/svgs/social-media.json";
import { useCookies } from "react-cookie";

type LoginPollingPageProps = {};

const LoginPollingPage = ({}: LoginPollingPageProps) => {
  const dispatch: AppDispatch = useDispatch();
  const { userId, accessToken, userProcessed } = useSelector((state: any) => {
    return {
      userId: getUserId(state),
      accessToken: getUserAccessToken(state),
      userProcessed: getUserProcessed(state),
    };
  });
  const [searchParams, setSearchParams] = useSearchParams();
  const [cookies, setCookie] = useCookies(["userId", "accessToken"]);
  const navigate = useNavigate();

  useEffect(() => {
    dispatch({
      type: "app/setNavBarVisible",
      payload: false,
    });
  }, []);

  useEffect(() => {
    let queryUserId = searchParams.get("userId");
    let queryAccessToken = searchParams.get("accessToken");
    if (queryUserId != null && queryAccessToken != null) {
      setCookie("userId", queryUserId, { path: "/" });
      setCookie("accessToken", queryAccessToken, { path: "/" });
      dispatch({
        type: "auth/setUserId",
        payload: queryUserId,
      });
      dispatch({
        type: "auth/setToken",
        payload: queryAccessToken,
      });
      setSearchParams("");
    } else {
      navigate("/login");
    }
  }, []);

  useEffect(() => {
    if (userId !== "") {
      if (userProcessed) {
        dispatch({
          type: "app/setNavBarVisible",
          payload: true,
        });
        navigate("/");
      } else {
        const pollFunc = async () => {
          let userProcessSuccess: boolean = await pollLoginProcess(userId);
          if (userProcessSuccess) {
            dispatch({
              type: "app/setUserProcessed",
              payload: true,
            });
          } else {
            console.log("user proccessing failed!");
          }
        };
        pollFunc();
      }
    }
  }, [userProcessed, userId]);

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
    animationData: loadingAnimation,
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
                  width={size === "small" ? "17rem" : "medium"}
                  direction="column"
                  align="center"
                >
                  <Grommet>
                    <Box margin="small">
                      <Heading margin="0">Hold On!</Heading>
                    </Box>
                    <Box margin="small">
                      <Text size={size === "small" ? "small" : "medium"}>
                        Please wait while we load your Spotify data...
                      </Text>
                    </Box>
                    <Box margin="small">
                      <Lottie
                        options={lottieDefaultOptions}
                        height={size === "small" ? "7rem" : "12rem"}
                        width={size === "small" ? "7rem" : "12rem"}
                      ></Lottie>
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

export default LoginPollingPage;
