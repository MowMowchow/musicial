import { MapDispatchToProps, useSelector } from "react-redux";
import { getUserId, getUserAccessToken } from "../../redux/reducers/auth";

import { TrackResultUser } from "../../utils/models/musicial-api";
import {
  Anchor,
  Avatar,
  Box,
  Card,
  CardBody,
  CardFooter,
  CardHeader,
  Grommet,
  ResponsiveContext,
  Text,
} from "grommet";
import { deepMerge } from "grommet/utils";
import { css } from "styled-components";
import userImageDnePlayer from "./../../utils/svgs/userImageDnePlayer.svg";
import { gColours } from "../../App";
type UserCardProps = { cardInfo: TrackResultUser };

const SearchTrackCard = ({ cardInfo }: UserCardProps) => {
  const searchTrackUserCard = deepMerge({
    anchor: {
      hover: {
        extend: ({}) => css`
          color: ${gColours.ming};
          text-decoration: none;
        `,
      },
    },
  });

  return (
    <Grommet theme={searchTrackUserCard}>
      <ResponsiveContext.Consumer>
        {(size) => (
          <Card
            key={cardInfo.userId}
            width={size === "small" ? "13rem" : " 15rem"}
            background={gColours.cardWhite}
            flex="grow"
            margin={"small"}
          >
            <CardHeader pad="small"></CardHeader>
            <CardBody>
              <Box>
                {console.log("USERIMAGE: ", cardInfo.userImage)}
                <Box
                  direction="row"
                  justify="center"
                  pad={{ horizontal: "large" }}
                >
                  {cardInfo.userImage !== "" ? (
                    <Avatar src={cardInfo.userImage} size="2xl" />
                  ) : (
                    <Avatar src={userImageDnePlayer} size="2xl" />
                  )}
                </Box>
                <Box
                  pad={{ vertical: "small", horizontal: "medium" }}
                  margin={{ bottom: "0" }}
                >
                  <Anchor
                    size="1.6rem"
                    color={gColours.eerieBlack}
                    href={cardInfo.userLink}
                    margin="0"
                  >
                    {cardInfo.userDisplayName}
                  </Anchor>
                </Box>
                <Box
                  pad={{ horizontal: "medium" }}
                  margin={{ bottom: "medium" }}
                >
                  <Text size=".99rem" color={gColours.eerieBlack}>
                    This song is included in:
                  </Text>
                  <Text size=".99rem" color={gColours.eerieBlack}>
                    {cardInfo.playlists.map((playlist, index) => {
                      return (
                        <>
                          <Anchor
                            href={playlist.playlistLink}
                            label={playlist.playlistName}
                            color={gColours.eerieBlack}
                            weight="normal"
                          />
                          <Text color={gColours.eerieBlack}>
                            {index < cardInfo.playlists.length - 1 ? ", " : ""}
                          </Text>
                        </>
                      );
                    })}
                  </Text>
                </Box>
              </Box>
            </CardBody>
            <CardFooter pad={{ horizontal: "small" }}></CardFooter>
          </Card>
        )}
      </ResponsiveContext.Consumer>
    </Grommet>
  );
};

export default SearchTrackCard;
