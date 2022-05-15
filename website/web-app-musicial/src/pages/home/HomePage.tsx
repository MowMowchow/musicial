import { useEffect } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";

import { useDispatch } from "react-redux";
import { AppDispatch } from "../../redux/store";
import { useSelector } from "react-redux";
import { getUserId, getUserAccessToken } from "../../redux/reducers/auth";

import {
  Anchor,
  Box,
  Button,
  Grid,
  Grommet,
  grommet,
  Header,
  Heading,
  Image,
  Main,
  ResponsiveContext,
  Stack,
  Text,
} from "grommet";
import { useCookies } from "react-cookie";
import { gColours } from "../../App";
import dancingPeople from "../../utils/svgs/dancingPeople.svg";
import network from "../../utils/svgs/network.svg";
import spotifyLogo from "../../utils/svgs/spotifyLogo.svg";
import boomBox from "../../utils/svgs/boombox.svg";
import combined from "../../utils/svgs/combined.svg";
import { deepMerge } from "grommet/utils";
import { css } from "styled-components";
import { fetchLoginRoute } from "../../utils/services/musicial-api";

type LoginPageProps = {};

const HomePage = ({}: LoginPageProps) => {
  const dispatch: AppDispatch = useDispatch();
  const { userId, accessToken } = useSelector((state: any) => {
    return {
      userId: getUserId(state),
      accessToken: getUserAccessToken(state),
    };
  });
  const [searchParams, setSearchParams] = useSearchParams();
  const [cookies, setCookie] = useCookies(["userId", "accessToken"]);
  const navigate = useNavigate();

  useEffect(() => {
    console.log("Home Page | COOKIES: ", cookies);
  }, [cookies]);

  useEffect(() => {
    if (
      cookies?.userId === "" ||
      cookies?.accessToken === "" ||
      Object.keys(cookies).length < 2
    ) {
      dispatch({
        type: "app/setNavBarVisible",
        payload: false,
      });
      navigate("/login");
    } else {
      dispatch({
        type: "auth/setUserId",
        payload: cookies.userId,
      });
      dispatch({
        type: "auth/setToken",
        payload: cookies.accessToken,
      });
    }
  }, []);

  const HomePageLoginButton = deepMerge({
    button: {
      extend: ({}) => css`
        border-radius: 7rem;
        border-width: 0.23rem;
        box-shadow: 0px 0.12rem 0.15rem rgba(0, 0, 0, 0.2);
      `,
    },
  });

  const HomePageImageWrapper = deepMerge({
    box: {
      extend: ({}) => css`
        min-width: 336px;
      `,
    },
  });

  return (
    <Grommet theme={HomePageLoginButton}>
      {" "}
      <ResponsiveContext.Consumer>
        {(size) => (
          <Box direction="column" justify="around" fill="horizontal">
            <Box
              direction={size === "small" ? "column-reverse" : "row"}
              justify="around"
              fill="horizontal"
              overflow="scroll"
            >
              <Box
                gridArea="homeLeft"
                direction="column"
                justify="around"
                flex="grow"
                margin={size === "small" ? { top: "3rem" } : { left: "medium" }}
              >
                <Box
                  direction="column"
                  justify="start"
                  align="center"
                  flex="grow"
                >
                  <Box
                    direction="column"
                    justify="start"
                    fill="horizontal"
                    flex="shrink"
                    margin={size === "small" ? "0" : { left: "4rem" }}
                  >
                    <Box
                      justify="start"
                      pad={size === "small" ? "0" : { left: "2rem" }}
                      margin={{ vertical: "medium" }}
                    >
                      <Heading
                        size="large"
                        margin="0"
                        color={gColours.raisinBlack}
                        alignSelf={size === "small" ? "center" : "start"}
                      >
                        Connect
                      </Heading>
                    </Box>
                    <Box
                      justify="start"
                      pad={size === "small" ? "0" : { left: "3.5rem" }}
                      margin={{ vertical: "medium" }}
                    >
                      <Heading
                        size="large"
                        margin="0"
                        color={gColours.raisinBlack}
                        alignSelf={size === "small" ? "center" : "start"}
                      >
                        With
                      </Heading>
                    </Box>
                    <Box
                      justify="start"
                      pad={size === "small" ? "0" : { left: "5rem" }}
                      margin={{ vertical: "medium" }}
                    >
                      <Heading
                        size="large"
                        margin="0"
                        color={gColours.raisinBlack}
                        alignSelf={size === "small" ? "center" : "start"}
                      >
                        Friends
                      </Heading>
                    </Box>
                  </Box>
                  <Box margin={{ top: "medium" }}>
                    <Button
                      label={
                        <Box
                          direction="row"
                          justify="start"
                          align="center"
                          margin="xxsmall"
                        >
                          <Text size="large">Connect with spotify</Text>
                          <Box width="5rem" margin={{ horizontal: "medium" }}>
                            <Image src={spotifyLogo} fill={true} />
                          </Box>
                        </Box>
                      }
                      size={size === "small" ? "small" : "medium"}
                      color={gColours.byzantium}
                      onClick={() => (window.location.href = fetchLoginRoute())}
                    />
                  </Box>
                </Box>
              </Box>
              <Box
                gridArea="homeRight"
                margin={
                  size === "small" ? { top: "2rem" } : { horizontal: "xlarge" }
                }
              >
                <Box direction="column" justify="center" fill={true}>
                  <Box
                    fill="horizontal"
                    direction="row"
                    justify="center"
                    pad={size === "small" ? "0" : { right: "small" }}
                  >
                    <Grommet theme={HomePageImageWrapper}>
                      <Box pad="0" align="center">
                        <Image
                          src={combined}
                          width={size === "small" ? "300rem" : "xlarge"}
                        />
                      </Box>
                    </Grommet>
                  </Box>
                </Box>
              </Box>
            </Box>
          </Box>
        )}
      </ResponsiveContext.Consumer>
    </Grommet>
  );
};

export default HomePage;
