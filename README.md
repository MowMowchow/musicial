# Musicial

### \*Currently in the process of requesting an extension from Spotify to avoid having to manually whitelist users to use the application

<br>

Musicial is a full-stack application that allows users to:

- ### Sign in with spotify
  ![Home Page Screenshot](https://i.imgur.com/WdCcrHX.png)
- ### Query for:

  - ### A user:

    Find and connect with other users! Click "Add User" to connect with them; allowing them to appear when querying the following data (below)

    ![Example User Query](https://i.imgur.com/wi6IfYA.png)

  - ### A song:

    Connected users with the song in any of their public playlists will appear along with links to all playlists containing the song

    ![Example Song Query](https://i.imgur.com/y0f8OV0.png)

  - ### An artist:

    Connected users that have a song by the artist in any of their public playlists will appear, showing each song and a link to the playlist(s) it belongs to

    ![Example Artist Query](https://i.imgur.com/HNMCMOD.png)

  - ### Album (COMING SOON):
    Connected users that have a song from a specific album in any of their public playlists will appear, showing each song and a link to the playlist(s) it belongs to

- ### \*\* Note \*\*:
  When querying for a song, artist or album, if the name does not appear in the suggestion dropdown, then none of your connections listen to the queried item

## How it was built:

- ### Technologies

  - Frontend: React, Redux, Typescrpt, Grommet UI library
  - Backend: AWS (Lambda, APIGateWay, S3), Golang
  - Database:
    - User Network: Neo4j
    - Cache: Redis
  - Other Tools/Services:

    - AWS (Route53, Cloudformation, Cloudwatch, Codebuild, CodePipeline)

  - ### Architecture diagram
    ![Architecture diagram](https://i.imgur.com/sd8qZbC.png)

## Future Plans:

- Allow users to query for albums
- Allow users to organize connections into groups
- Allow users to remove connections
- Update UI and provide more information for queried items (users, songs, arists, albums)
