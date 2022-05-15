import React, { useEffect, useState } from "react";
import { useDispatch } from "react-redux";
import { AppDispatch } from "../../redux/store";
import { useSelector } from "react-redux";

import { getNavBarVisible } from "../../redux/reducers/web-app";
import { getUserId, getUserAccessToken } from "../../redux/reducers/auth";
import { useNavigate } from "react-router-dom";

import {
  Box,
  Grommet,
  ResponsiveContext,
  Tab,
  Tabs,
  Text,
  TextInput,
} from "grommet";
import { gColours } from "../../App";
import SearchTrackResults from "./SearchTrackResults";
import { css } from "styled-components";
import { deepMerge } from "grommet/utils";
import SearchArtistResults from "./SearchArtistResults";
import SearchUserResults from "./SearchUserResults";
import { useCookies } from "react-cookie";
type SearchPageProps = {};

const SearchPage = ({}: SearchPageProps) => {
  const dispatch: AppDispatch = useDispatch();
  const { userId, accessToken } = useSelector((state: any) => {
    return {
      userId: getUserId(state),
      accessToken: getUserAccessToken(state),
    };
  });
  const [searchInput, setSearchInput] = useState("");
  const [suggestedList, setSuggestedList] = useState<Array<Object>>(
    new Array()
  );
  const [selectedItem, setSelectedItem] = useState<string>("");
  const [searchPage, setSearchPage] = useState<string>("tracks");
  const [cookies, setCookie] = useCookies(["userId", "accessToken"]);
  const navigate = useNavigate();

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

  useEffect(() => {
    console.log("Search Page | COOKIES: ", cookies);
  }, [cookies]);

  const updateSearchPage = (searchPage: string) => {
    setSearchPage(searchPage);
    setSuggestedList(new Array());
    setSelectedItem("");
    setSearchInput("");
  };

  const SearchPageTheme = deepMerge({
    global: {},
    tabs: {
      header: {
        size: "medium",
      },
    },
    tab: {
      border: undefined,
      active: {
        background: {
          color: gColours.byzantium,
          opacity: 1.0,
        },
        color: "#FFFFFF",
        extend: ({}) => css`
          transition: transform 0.3s;
          transform: scale(20);
        `,
      },
      hover: {
        background: gColours.byzantium,
        opacity: "weak",
        color: "#FFFFFF",
      },
      background: {
        color: gColours.byzantium,
        opacity: 0.7,
      },
      pad: {
        top: "small",
        right: "small",
        left: "small",
        bottom: "small",
      },
      extend: ({}) => css`
        border-radius: 16px 16px 0% 0%;
        box-shadow: 0px 1px 5px rgba(0, 0, 0, 0.4);
      `,
    },
    textInput: {
      extend: ({}) => css`
        background: #ffffff;
        color: ${gColours.eerieBlack};
      `,

      suggestions: {
        extend: ({}) => css`
          width: ${document.getElementById("textInputWrapperBox")
            ?.offsetWidth}px;
        `,
      },
    },
  });

  return (
    <Grommet theme={SearchPageTheme}>
      <ResponsiveContext.Consumer>
        {(size) => (
          <Box margin={{ vertical: "5rem", horizontal: "10.5rem" }}>
            <Box
              direction={size === "small" ? "column" : "row"}
              justify="start"
            >
              <Box margin="small" flex="grow">
                <Tabs alignSelf={size === "small" ? "center" : "start"}>
                  <Tab
                    title={<Text size="large">Songs</Text>}
                    onClick={() => updateSearchPage("tracks")}
                  ></Tab>
                  <Tab
                    title={<Text size="large">Artists</Text>}
                    onClick={() => updateSearchPage("artists")}
                  ></Tab>
                  {/* <Tab
                title={<Text size="large">Albums</Text>}
                onClick={() => updateSearchPage("albums")}
              ></Tab> */}
                  <Tab
                    title={<Text size="large">Users</Text>}
                    onClick={() => updateSearchPage("users")}
                  ></Tab>
                </Tabs>
              </Box>
              <Box
                margin="small"
                fill="horizontal"
                direction="column"
                justify="center"
                id="textInputWrapperBox"
              >
                <TextInput
                  value={searchInput}
                  dropHeight="medium"
                  placeholder="don't be shy :)"
                  onChange={(e) => setSearchInput(e.target.value)}
                  suggestions={suggestedList}
                  onSuggestionSelect={(suggestionObj) => {
                    console.log("CHOSE: ", suggestionObj.suggestion.value);
                    setSelectedItem(suggestionObj.suggestion.value);
                    setSearchInput(suggestionObj.suggestion.label);
                  }}
                ></TextInput>
              </Box>
            </Box>
            <Box margin={{ horizontal: "large", top: "medium" }}>
              {searchPage === "tracks" ? (
                <SearchTrackResults
                  searchInput={searchInput}
                  setSuggestedList={setSuggestedList}
                  selectedTrack={selectedItem}
                />
              ) : searchPage === "artists" ? (
                <SearchArtistResults
                  searchInput={searchInput}
                  setSuggestedList={setSuggestedList}
                  selectedArtist={selectedItem}
                />
              ) : (
                <SearchUserResults
                  searchInput={searchInput}
                ></SearchUserResults>
              )}
            </Box>
          </Box>
        )}
      </ResponsiveContext.Consumer>
    </Grommet>
  );
};

export default SearchPage;
