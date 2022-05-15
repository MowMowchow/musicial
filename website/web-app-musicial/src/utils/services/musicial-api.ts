import {
  Neo4jAllUserWrapper,
  UserWithQueryName,
  Neo4jAllTrackWrapper,
  TrackResult,
  TrackResultUser,
  Neo4jAllArtistWrapper,
  ArtistResult,
  ArtistResultUser,
  ArtistResultTrack,
} from "../models/musicial-api";

const baseApiUrl = "https://api.musicial.net/api/";

export const setCustomTimeout = async (seconds: number) => {
  return new Promise((resolve) => setTimeout(resolve, 1000 * seconds));
};

export const fetchLoginRoute = () => {
  return baseApiUrl + "login";
};

export const fetchAllExistingUsers = async (
  userId: string
): Promise<Neo4jAllUserWrapper> => {
  const response = await fetch(baseApiUrl + "fetch/users/" + userId + "/all");
  if (!response.ok) {
    console.error("Error fetching all users | ", response.status);
    throw Error("Error fetching all users");
  }
  let data = (await response.json()) as Neo4jAllUserWrapper;
  console.log("USERS DATA: ", data);
  data.users = new Map<string, UserWithQueryName>(Object.entries(data.users));
  data.queryNameToUserId = new Map<string, string[]>(
    Object.entries(data.queryNameToUserId)
  );
  return data;
};

export const AddUserAsConnection = async (
  userId: string,
  accessToken: string,
  connectId: string
): Promise<Boolean> => {
  // if (userId)
  const response = await fetch(baseApiUrl + "user/add/" + connectId, {
    method: "PUT",
    body: JSON.stringify({
      userId: userId,
      accessToken: accessToken,
    }),
  });
  if (!response.ok) {
    console.error("Error connecting to user: ", connectId);
    return false;
  }
  return true;
};

export const fetchTracks = async (
  userId: string,
  accessToken: string
): Promise<Neo4jAllTrackWrapper> => {
  const response = await fetch(baseApiUrl + "fetch/tracks/" + userId, {
    method: "POST",
    body: JSON.stringify({
      userId: userId,
      accessToken: accessToken,
    }),
  });
  if (!response.ok) {
    console.error("Error fetching all tracks for user | ", response.status);
    throw Error("Error fetching all tracks for user");
  }

  let data: Neo4jAllTrackWrapper = await response.json();

  data.queryNameToTrackId = new Map<string, string[]>(
    Object.entries(data.queryNameToTrackId)
  );
  data.tracks = new Map<string, TrackResult>(Object.entries(data.tracks));
  data.tracks.forEach((value, key) => {
    let tempTrackResult = data.tracks.get(key)!;
    tempTrackResult.users = new Map<string, TrackResultUser>(
      Object.entries(value.users)
    );
  });
  return data;
};

export const fetchArtists = async (
  userId: string,
  accessToken: string
): Promise<Neo4jAllArtistWrapper> => {
  const response = await fetch(baseApiUrl + "fetch/artists/" + userId, {
    method: "POST",
    body: JSON.stringify({
      userId: userId,
      accessToken: accessToken,
    }),
  });
  if (!response.ok) {
    console.error("Error fetching all artists for user | ", response.status);
    throw Error("Error fetching all artists for user");
  }

  let data: Neo4jAllArtistWrapper = await response.json();

  data.queryNameToArtistId = new Map<string, string[]>(
    Object.entries(data.queryNameToArtistId)
  );
  data.artists = new Map<string, ArtistResult>(Object.entries(data.artists));
  data.artists.forEach((value, key) => {
    let tempArtistResult = data.artists.get(key)!;
    tempArtistResult.users = new Map<string, ArtistResultUser>(
      Object.entries(value.users)
    );
    tempArtistResult.users.forEach((user) => {
      let tempUser = tempArtistResult.users.get(user.userId)!.tracks;
      tempArtistResult.users.get(user.userId)!.tracks = new Map<
        string,
        ArtistResultTrack
      >(Object.entries(tempUser));
    });
  });
  return data;
};

export const pollLoginProcess = async (userId: string): Promise<boolean> => {
  let maxTime: number = 300; // seconds
  let currentTime: number = 0; // seconds
  let timeoutId: NodeJS.Timeout;
  let loginSuccess: boolean = false;
  let data;

  while (currentTime < maxTime) {
    console.log(
      "ABOUT TO QUERY PROCESS STATUS FOR USER: ",
      userId,
      " | currentTime: ",
      currentTime
    );
    const response = await fetch(baseApiUrl + "login/pollProcess/" + userId, {
      method: "GET",
    });
    if (!response.ok) {
      console.error(
        "Error fetching update process status for user | ",
        response.status
      );
      throw Error("Error fetching update process status for user");
    }

    data = await response.json();
    console.log(
      "PROCESS STATUS FOR USER: ",
      userId,
      " | currentTime: ",
      currentTime,
      " | data: ",
      data
    );

    if (data?.isProcessing === 1) {
      console.log("waiting 5 seconds before querying login process again");
      await setCustomTimeout(5);
      currentTime += 5;
    } else if (data?.isProcessing === 0) {
      loginSuccess = true;
      break;
    } else {
      console.log(
        "'isProcessing' field not found in response body | typeof data?.isProcessing: ",
        typeof data?.isProcessing
      );
      loginSuccess = false;
      break;
    }
  }
  // clearTimeout(timeoutId!);

  console.log(
    "user data processing done | pollLoginProcess returning | data: ",
    data
  );
  return loginSuccess;
};
