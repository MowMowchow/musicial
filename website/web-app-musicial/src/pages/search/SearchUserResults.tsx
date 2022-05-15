import React, { useEffect, useState } from "react";
import { useDispatch } from "react-redux";
import { AppDispatch } from "../../redux/store";
import { useSelector } from "react-redux";

import { getUserId, getUserAccessToken } from "../../redux/reducers/auth";
import { useNavigate } from "react-router-dom";

import { Box, Grommet, ResponsiveContext } from "grommet";
import { Trie } from "../../utils/helpers/trie";
import {
  Neo4jAllUserWrapper,
  UserWithQueryName,
} from "../../utils/models/musicial-api";
import SearchUserCard from "../../components/searchUserCard/SearchUserCard";
import { deepMerge } from "grommet/utils";
import { fetchAllExistingUsers } from "../../utils/services/musicial-api";

type SearchUserResultsProps = {
  searchInput: string;
};

var suggestionTrie = new Trie();

const SearchUserResults = ({ searchInput }: SearchUserResultsProps) => {
  const dispatch: AppDispatch = useDispatch();
  const { userId, accessToken } = useSelector((state: any) => {
    return {
      userId: getUserId(state),
      accessToken: getUserAccessToken(state),
    };
  });
  const navigate = useNavigate();
  const [allUsers, setAllUsers] = useState<Neo4jAllUserWrapper>();
  const [suggestedUsers, setSuggestedUsers] = useState<UserWithQueryName[]>();
  const [newlyConnectedUsers, setNewlyConnectedUsers] = useState<Set<string>>(
    new Set<string>()
  );
  const updateUserResults = (searchQuery: string) => {
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
      console.log("SUGGESTED USER LIST: ", suggestedList);
      console.log("SUGGESTED TRIE LIST: ", suggestionTrie.suggestionResults);
    }
  };

  useEffect(() => {
    const fetchData = async () => {
      if (userId !== "") {
        const data: Neo4jAllUserWrapper = await fetchAllExistingUsers(userId);
        suggestionTrie.build(Array.from(data.queryNameToUserId.keys()));
        setAllUsers(data);
      }
    };
    fetchData();
  }, []);

  useEffect(() => {
    if (allUsers !== undefined) {
      updateUserResults(searchInput);
    }
  }, [searchInput]);

  const searchUserResultTheme = deepMerge({
    global: {
      control: {
        border: {
          radius: "23px",
        },
      },
    },
  });

  return (
    <Grommet theme={searchUserResultTheme}>
      <ResponsiveContext.Consumer>
        {(size) => (
          <Box fill="horizontal">
            <Box direction="column" justify="center" fill="horizontal">
              <Box
                direction="row"
                justify={size === "small" ? "center" : "start"}
                fill="horizontal"
                wrap={true}
              >
                {suggestedUsers !== undefined ? (
                  suggestedUsers!.map((user) => {
                    if (!newlyConnectedUsers.has(user.userId)) {
                      return (
                        <SearchUserCard
                          cardInfo={user}
                          newlyConnectedUsers={newlyConnectedUsers}
                          setNewlyConnectedUsers={setNewlyConnectedUsers}
                        />
                      );
                    }
                  })
                ) : (
                  <></>
                )}
              </Box>
            </Box>
          </Box>
        )}
      </ResponsiveContext.Consumer>
    </Grommet>
  );
};

export default SearchUserResults;
