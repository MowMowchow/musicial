import React, { useEffect, useState } from "react";

import { MapDispatchToProps, useSelector } from "react-redux";
import { getUserId, getUserAccessToken } from "../../redux/reducers/auth";

import { UserWithQueryName } from "../../utils/models/musicial-api";
import { AddUserAsConnection } from "../../utils/services/musicial-api";
import {
  Anchor,
  Avatar,
  Box,
  Button,
  Card,
  CardBody,
  CardFooter,
  CardHeader,
  Grommet,
  ResponsiveContext,
  Text,
} from "grommet";
import { gColours } from "../../App";
import { css } from "styled-components";
import userImageDnePlayer from "./../../utils/svgs/userImageDnePlayer.svg";
import { deepMerge } from "grommet/utils";

type UserCardProps = {
  cardInfo: UserWithQueryName;
  newlyConnectedUsers: Set<string>;
  setNewlyConnectedUsers: React.Dispatch<React.SetStateAction<Set<string>>>;
};

const UserCards = ({
  cardInfo,
  newlyConnectedUsers,
  setNewlyConnectedUsers,
}: UserCardProps) => {
  const { userId, accessToken } = useSelector((state: any) => {
    return {
      userId: getUserId(state),
      accessToken: getUserAccessToken(state),
    };
  });
  const searchUserCard = deepMerge({
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
    <Grommet theme={searchUserCard}>
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
                <Box pad={{ horizontal: "medium" }}></Box>
              </Box>
            </CardBody>
            <CardFooter
              pad={{ horizontal: "small", vertical: "xsmall" }}
              background="light-2"
            >
              <Button
                color={gColours.byzantium}
                size="small"
                label={
                  <Text size="small" color={gColours.eerieBlack}>
                    Add User
                  </Text>
                }
                onClick={() => {
                  AddUserAsConnection(userId, accessToken, cardInfo.userId);
                  setNewlyConnectedUsers(
                    (prevState) =>
                      new Set(newlyConnectedUsers.add(cardInfo.userId))
                  );
                }}
              />
            </CardFooter>
          </Card>
        )}
      </ResponsiveContext.Consumer>
    </Grommet>
  );
};

export default UserCards;
