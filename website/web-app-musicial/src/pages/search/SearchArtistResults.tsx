import React, { useEffect, useState } from "react";
import { useDispatch } from "react-redux";
import { AppDispatch } from "../../redux/store";
import { useSelector } from "react-redux";

import { getNavBarVisible } from "../../redux/reducers/web-app";
import { getUserId, getUserAccessToken } from "../../redux/reducers/auth";
import { useNavigate } from "react-router-dom";

import { Box, Grommet } from "grommet";
import { fetchArtists } from "../../utils/services/musicial-api";
import { Trie } from "../../utils/helpers/trie";
import {
  Neo4jAllArtistWrapper,
  ArtistResultUser,
  ArtistSuggestion,
} from "../../utils/models/musicial-api";
import { deepMerge } from "grommet/utils";
import SearchArtistCard from "../../components/searchArtistCard/SearchArtistCard";

type SearchArtistResultsProps = {
  searchInput: string;
  setSuggestedList: React.Dispatch<React.SetStateAction<Array<any>>>;
  selectedArtist: string;
};

var suggestionTrie = new Trie();

const SearchArtistResults = ({
  searchInput,
  setSuggestedList,
  selectedArtist,
}: SearchArtistResultsProps) => {
  const dispatch: AppDispatch = useDispatch();
  const { userId, accessToken } = useSelector((state: any) => {
    return {
      userId: getUserId(state),
      accessToken: getUserAccessToken(state),
    };
  });
  const navigate = useNavigate();
  const [allArtists, setAllArtists] = useState<Neo4jAllArtistWrapper>();
  const [suggestedArtists, setSuggestedArtists] = useState<ArtistSuggestion[]>(
    []
  );
  const [shownUsers, setShownUsers] = useState<ArtistResultUser[]>();

  const updateArtistResults = (searchQuery: string) => {
    if (allArtists?.queryNameToArtistId.size! > 0) {
      suggestionTrie.getSuggestions(
        searchQuery.replaceAll(" ", "").toLocaleLowerCase()
      );
      let suggestedList: ArtistSuggestion[] = [];
      let fetchedArtists: string[];
      let tempFetchedArtist: ArtistSuggestion, fetchedArtist: string;
      suggestionTrie.suggestionResults.forEach((queryName) => {
        fetchedArtists = allArtists?.queryNameToArtistId.get(queryName)!;
        if (fetchedArtists !== undefined) {
          fetchedArtists.forEach((currArtistId) => {
            tempFetchedArtist = allArtists!.artists.get(currArtistId)!;
            if (tempFetchedArtist !== undefined) {
              let fetchedArtist: ArtistSuggestion = {
                artistId: tempFetchedArtist.artistId,
                artistName: tempFetchedArtist.artistName,
              };
              suggestedList.push(fetchedArtist);
            }
          });
        }
      });
      setSuggestedArtists(suggestedList);
      setSuggestedList(suggestedList);
      setSuggestedList(
        suggestedArtists.map((artist) => {
          return {
            label: artist.artistName,
            value: artist.artistId,
          };
        })
      );
    }
  };

  useEffect(() => {
    const fetchData = async () => {
      if (userId !== "") {
        const data: Neo4jAllArtistWrapper = await fetchArtists(
          userId,
          accessToken
        );
        suggestionTrie.build(Array.from(data.queryNameToArtistId.keys()));
        setAllArtists(data);
      }
    };
    fetchData();
  }, []);

  useEffect(() => {
    if (allArtists !== undefined) {
      updateArtistResults(searchInput);
    }
  }, [searchInput]);

  useEffect(() => {
    if (allArtists !== undefined) {
      const artist = allArtists.artists.get(selectedArtist)!;
      if (artist !== undefined) {
        setShownUsers(Array.from(artist.users.values()));
      }
    }
  }, [selectedArtist]);

  const searchArtistResultTheme = deepMerge({
    global: {},
  });

  return (
    <Grommet theme={searchArtistResultTheme}>
      <Box>
        <Box direction="row" justify="center" fill="horizontal">
          <Box direction="column" justify="center" fill="horizontal">
            <Box direction="column" align="center" fill="horizontal">
              {shownUsers !== undefined ? (
                shownUsers!.map((user) => {
                  return (
                    <SearchArtistCard
                      cardInfo={user}
                      selectedArtist={selectedArtist}
                      key={user.userId}
                    />
                  );
                })
              ) : (
                <> </>
              )}
            </Box>
          </Box>
        </Box>
      </Box>
    </Grommet>
  );
};

export default SearchArtistResults;
