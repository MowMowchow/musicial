import React, { useEffect, useState } from "react";
import { Link } from "react-router-dom";

import { MapDispatchToProps, useSelector } from "react-redux";
import { getUserId, getUserAccessToken } from "../../redux/reducers/auth";

import {
  Neo4jAllUserWrapper,
  ArtistResult,
  ArtistResultUser,
  UserWithQueryName,
} from "../../utils/models/musicial-api";
import { AddUserAsConnection } from "../../utils/services/musicial-api";
import {
  Accordion,
  AccordionPanel,
  Anchor,
  Avatar,
  Box,
  Card,
  Grommet,
  ResponsiveContext,
  Text,
} from "grommet";
import { deepMerge } from "grommet/utils";
import { css } from "styled-components";
import userImageDnePlayer from "./../../utils/svgs/userImageDnePlayer.svg";
import { gColours } from "../../App";
type UserCardProps = { cardInfo: ArtistResultUser; selectedArtist: string };

const SearchArtistCard = ({ cardInfo, selectedArtist }: UserCardProps) => {
  const { userId, accessToken } = useSelector((state: any) => {
    return {
      userId: getUserId(state),
      accessToken: getUserAccessToken(state),
    };
  });
  const [accordianExpanded, setAccordianExpanded] = useState(false);

  const noAccordianBorder = deepMerge({
    accordion: {
      icons: {
        color: gColours.eerieBlack,
      },
      border: {
        style: "hidden",
        color: gColours.fieryRose,
      },
    },
  });

  const searchArtistUserAccordianTheme = deepMerge({
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
          color: ${gColours.byzantium};
          text-decoration: none;
        `,
      },
    },
  });

  const searchArtistTrackBoxBorder = deepMerge({
    box: {
      extend: ({}) => css`
        // border-radius: 1rem;
        // box-shadow: 0px 0.12rem 0.25rem rgba(0, 0, 0, 0.2);
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

  const searchArtistTrackBoxTheme = deepMerge({
    box: {
      extend: ({}) => css`
        border-radius: 1rem;
      `,
    },
    avatar: {
      extend: ({}) => css`
        border-radius: 0.3rem;
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
    accordion: {
      icons: {
        color: gColours.eerieBlack,
      },
      border: {
        style: "hidden",
      },
    },
  });

  useEffect(() => {
    console.log("ACCORDIAN EXPANDED: ", accordianExpanded);
  }, [accordianExpanded]);

  return (
    <Grommet theme={searchArtistUserAccordianTheme}>
      <ResponsiveContext.Consumer>
        {(size) => (
          <Accordion
            width={accordianExpanded ? "100%" : "40vw"}
            animate={true}
            key={cardInfo.userId}
            background={gColours.cardWhite}
            margin={{ bottom: "small" }}
            activeIndex={accordianExpanded ? 0 : 1}
            onActive={() => setAccordianExpanded((prevState) => !prevState)}
          >
            <Box>
              <Grommet theme={noAccordianBorder}>
                <AccordionPanel
                  label={
                    <Box
                      direction="row"
                      pad={size === "small" ? "xsmall" : "small"}
                    >
                      {cardInfo.userImage !== "" ? (
                        <Box justify="center">
                          <Avatar
                            src={cardInfo.userImage}
                            size={size === "small" ? "medium" : "xlarge"}
                          />
                        </Box>
                      ) : (
                        <Box justify="center">
                          <Avatar
                            src={userImageDnePlayer}
                            size={size === "small" ? "medium" : "xlarge"}
                          />
                        </Box>
                      )}
                      <Box justify="center" margin={{ horizontal: "medium" }}>
                        <Grommet theme={searchArtistUserAccordianTheme}>
                          <Anchor
                            size={size === "small" ? "small" : "large"}
                            label={cardInfo.userDisplayName}
                          />
                        </Grommet>
                      </Box>
                    </Box>
                  }
                >
                  <Box
                    direction="row"
                    justify="around"
                    wrap={true}
                    pad={{ horizontal: "xsmall" }}
                  >
                    {Array.from(cardInfo.tracks.values()).map((track) => (
                      <Grommet theme={searchArtistTrackBoxBorder}>
                        <Box
                          margin={"xsmall"}
                          width={size === "small" ? "10rem" : "small"}
                          flex="grow"
                        >
                          <Grommet theme={searchArtistTrackBoxTheme}>
                            <Card direction="row" justify="center" flex="grow">
                              <Box
                                direction="column"
                                justify="start"
                                fill="horizontal"
                                height={
                                  size === "small"
                                    ? { min: "10rem" }
                                    : { min: "12rem" }
                                }
                              >
                                <Box
                                  pad={size === "small" ? "medium" : "small"}
                                  direction="row"
                                  justify="center"
                                >
                                  {track.trackImage !== "" ? (
                                    <Box justify="center">
                                      <Avatar
                                        src={track.trackImage}
                                        size="4.5rem"
                                      />
                                    </Box>
                                  ) : (
                                    <Box justify="center">
                                      <Avatar
                                        src={userImageDnePlayer}
                                        size="4.5rem"
                                      />
                                    </Box>
                                  )}
                                </Box>
                                <Box
                                  direction="column"
                                  justify="around"
                                  fill="horizontal"
                                  pad={
                                    size === "small"
                                      ? {
                                          horizontal: "small",
                                          bottom: "xsmall",
                                        }
                                      : { horizontal: "small", bottom: "small" }
                                  }
                                >
                                  <Box direction="column" justify="around">
                                    <Text>
                                      <Text
                                        size={
                                          size === "small" ? "xsmall" : "small"
                                        }
                                        color={gColours.byzantium}
                                      >
                                        Song:{" "}
                                      </Text>
                                      <Text
                                        size={
                                          size === "small" ? "xsmall" : "small"
                                        }
                                        color={gColours.eerieBlack}
                                      >
                                        {track.trackName}
                                      </Text>
                                    </Text>
                                    <Box direction="row">
                                      <Box>
                                        <Text
                                          size={
                                            size === "small"
                                              ? "xsmall"
                                              : "small"
                                          }
                                          color={gColours.byzantium}
                                        >
                                          Included in:
                                        </Text>
                                        <Text
                                          size={
                                            size === "small"
                                              ? "xsmall"
                                              : "small"
                                          }
                                          color={gColours.byzantium}
                                        >
                                          {track.playlists.map(
                                            (playlist, index) => (
                                              <>
                                                <Anchor
                                                  href={playlist.playlistLink}
                                                  label={playlist.playlistName}
                                                  color={gColours.eerieBlack}
                                                  weight="normal"
                                                  size={
                                                    size === "small"
                                                      ? "xsmall"
                                                      : "small"
                                                  }
                                                />
                                                <Text
                                                  size={
                                                    size === "small"
                                                      ? "xsmall"
                                                      : "small"
                                                  }
                                                  color={gColours.eerieBlack}
                                                >
                                                  {index <
                                                  track.playlists.length - 1
                                                    ? ", "
                                                    : ""}
                                                </Text>
                                              </>
                                            )
                                          )}
                                        </Text>
                                      </Box>
                                    </Box>
                                  </Box>
                                </Box>
                              </Box>
                            </Card>
                          </Grommet>
                        </Box>
                      </Grommet>
                    ))}
                  </Box>
                  <Box margin="xsmall" />
                </AccordionPanel>
              </Grommet>
            </Box>
          </Accordion>
        )}
      </ResponsiveContext.Consumer>
    </Grommet>
  );
};

export default SearchArtistCard;
