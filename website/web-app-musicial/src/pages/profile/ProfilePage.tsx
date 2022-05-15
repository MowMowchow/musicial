import { useEffect, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";

import { useDispatch } from "react-redux";
import { AppDispatch } from "../../redux/store";
import { useSelector } from "react-redux";
import { getUserId, getUserAccessToken } from "../../redux/reducers/auth";

import { Box, Main, Paragraph } from "grommet";

import { Trie } from "../../utils/helpers/trie";
import {
  Neo4jAllUserWrapper,
  UserWithQueryName,
} from "../../utils/models/musicial-api";
import { fetchAllExistingUsers } from "../../utils/services/musicial-api";

import { useCookies } from "react-cookie";

type LoginPageProps = {};

var suggestionTrie = new Trie();

const ProfilePage = ({}: LoginPageProps) => {
  const dispatch: AppDispatch = useDispatch();
  const { userId, accessToken } = useSelector((state: any) => {
    return {
      userId: getUserId(state),
      accessToken: getUserAccessToken(state),
    };
  });
  const navigate = useNavigate();
  const [searchInput, setSearchInput] = useState("");
  const [allUsers, setAllUsers] = useState<Neo4jAllUserWrapper>();
  const [suggestedUsers, setSuggestedUsers] = useState<UserWithQueryName[]>([]);
  const [newlyConnectedUsers, setNewlyConnectedUsers] = useState<Set<string>>(
    new Set<string>()
  );
  const [cookies, setCookie] = useCookies(["userId", "accessToken"]);

  useEffect(() => {
    console.log("Profile Page | COOKIES: ", cookies);
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

  const updateUserSuggestions = (searchQuery: string) => {
    if (allUsers!.users.size > 0) {
      suggestionTrie.getSuggestions(
        searchQuery.replaceAll(" ", "").toLocaleLowerCase()
      );
      let suggestedList: UserWithQueryName[] = [];
      let fetchedUsers: string[];
      let fetchedUser: UserWithQueryName;
      suggestionTrie.suggestionResults.forEach((queryDisplayName) => {
        fetchedUsers = allUsers!.queryNameToUserId.get(queryDisplayName)!;
        if (fetchedUsers !== undefined) {
          fetchedUsers.forEach((currUserId) => {
            fetchedUser = allUsers!.users.get(currUserId)!;
            if (
              fetchedUser.userId !== userId &&
              !newlyConnectedUsers.has(userId)
            ) {
              suggestedList.push(fetchedUser);
            }
          });
        }
      });
      setSuggestedUsers(suggestedList);
    }
  };

  useEffect(() => {
    const fetchData = async () => {
      const data: Neo4jAllUserWrapper = await fetchAllExistingUsers(userId);
      suggestionTrie.build(Array.from(data.queryNameToUserId.keys()));
      setAllUsers(data);
    };
    fetchData();
  }, []);

  useEffect(() => {
    if (allUsers !== undefined) {
      updateUserSuggestions(searchInput);
    }
  }, [searchInput, allUsers, newlyConnectedUsers]);

  return (
    <Main>
      <Box>
        <Paragraph>UserId: {userId}</Paragraph>
        <br />
        <Paragraph>AccessToken: {accessToken}</Paragraph>
      </Box>
      <br />
      <Box></Box>
    </Main>
  );
};

export default ProfilePage;
