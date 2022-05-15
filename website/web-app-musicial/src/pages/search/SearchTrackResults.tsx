import React, { useEffect, useState } from "react";
import { useDispatch } from "react-redux";
import { AppDispatch } from "../../redux/store";
import { useSelector } from "react-redux";

import { getUserId, getUserAccessToken } from "../../redux/reducers/auth";
import { useNavigate } from "react-router-dom";

import { Box, Grommet, ResponsiveContext } from "grommet";
import { fetchTracks } from "../../utils/services/musicial-api";
import { Trie } from "../../utils/helpers/trie";
import {
  Neo4jAllTrackWrapper,
  TrackResult,
  TrackResultUser,
  TrackSuggestion,
} from "../../utils/models/musicial-api";
import SearchTrackCard from "../../components/searchTrackCard/SearchTrackCard";
import { deepMerge } from "grommet/utils";

type SearchTrackResultsProps = {
  searchInput: string;
  setSuggestedList: React.Dispatch<React.SetStateAction<Array<any>>>;
  selectedTrack: string;
};

var suggestionTrie = new Trie();

const truncate = (str: string, n: number) => {
  return str.length > n ? str.substring(0, n - 1) : str;
};

const SearchTrackResults = ({
  searchInput,
  setSuggestedList,
  selectedTrack,
}: SearchTrackResultsProps) => {
  const dispatch: AppDispatch = useDispatch();
  const { userId, accessToken } = useSelector((state: any) => {
    return {
      userId: getUserId(state),
      accessToken: getUserAccessToken(state),
    };
  });
  const navigate = useNavigate();
  const [allTracks, setAllTracks] = useState<Neo4jAllTrackWrapper>();
  const [suggestedTracks, setSuggestedTracks] = useState<TrackSuggestion[]>([]);
  const [shownUsers, setShownUsers] = useState<TrackResultUser[]>();

  const updateTrackResults = (searchQuery: string) => {
    if (allTracks?.queryNameToTrackId.size! > 0) {
      suggestionTrie.getSuggestions(
        searchQuery.replaceAll(" ", "").toLocaleLowerCase()
      );
      let suggestedList: TrackSuggestion[] = [];
      let fetchedTracks: string[];
      let tempFetchedTrack: TrackResult, fetchedTrack: TrackSuggestion;
      suggestionTrie.suggestionResults.forEach((queryName) => {
        fetchedTracks = allTracks?.queryNameToTrackId.get(queryName)!;
        if (fetchedTracks !== undefined) {
          fetchedTracks.forEach((currTrackId) => {
            tempFetchedTrack = allTracks!.tracks.get(currTrackId)!;
            if (tempFetchedTrack !== undefined) {
              let fetchedTrack: TrackSuggestion = {
                trackId: tempFetchedTrack.trackId,
                trackName: tempFetchedTrack.trackName,
                trackImage: tempFetchedTrack.trackImage,
                trackArtists: tempFetchedTrack.artists
                  .map((artist) => artist.artistName)
                  .join(", "),
              };
              suggestedList.push(fetchedTrack);
            }
          });
        }
      });
      setSuggestedTracks(suggestedList);
      console.log("SUGGESTED TRACK LIST: ", suggestedList);
      console.log("SUGGESTED TRIE LIST: ", suggestionTrie.suggestionResults);
      setSuggestedList(
        suggestedTracks.map((track) => {
          return {
            label: track.trackName + " - " + track.trackArtists,
            value: track.trackId,
          };
        })
      );
    }
  };

  useEffect(() => {
    const fetchData = async () => {
      if (userId !== "") {
        const data: Neo4jAllTrackWrapper = await fetchTracks(
          userId,
          accessToken
        );
        suggestionTrie.build(Array.from(data.queryNameToTrackId.keys()));
        setAllTracks(data);
      }
    };
    fetchData();
  }, []);

  useEffect(() => {
    if (allTracks !== undefined) {
      updateTrackResults(searchInput);
    }
  }, [searchInput]);

  useEffect(() => {
    if (allTracks !== undefined) {
      const track = allTracks.tracks.get(selectedTrack)!;
      if (track !== undefined) {
        setShownUsers(Array.from(track.users.values()));
      }
    }
  }, [selectedTrack]);

  const searchTrackResultTheme = deepMerge({
    global: {
      control: {
        border: {
          radius: "23px",
        },
      },
    },
    // breakpoints: {
    //   medium: {
    //     value: 1280,
    //   },
    // },
  });

  return (
    <Grommet theme={searchTrackResultTheme}>
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
                {shownUsers !== undefined ? (
                  shownUsers!.map((user) => <SearchTrackCard cardInfo={user} />)
                ) : (
                  <> </>
                )}
              </Box>
            </Box>
          </Box>
        )}
      </ResponsiveContext.Consumer>
    </Grommet>
  );
};

export default SearchTrackResults;
