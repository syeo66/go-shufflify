# Go Shufflify

Go Shufflify is a web application built with Go, which interacts with the Spotify API to manage user queues, playback state, and player information. The app is designed to help users shuffle their Spotify playlists and manage their music experience in a more interactive way.

## Features

- **User Authentication**: Secure user authentication using Spotify OAuth.
- **Queue Management**: View and manage the current queue of tracks.
- **Player State**: Display the current playing track, and player state including shuffle and repeat modes.
- **Database Support**: Stores user data securely in a SQLite database.
- **Session Management**: Maintains user sessions with Gorilla sessions.

## Getting Started

### Prerequisites

- [Go](https://golang.org/dl/)
- [SQLite](https://www.sqlite.org/download.html)
- Spotify Developer Account and an app setup on the [Spotify Developer Dashboard](https://developer.spotify.com/dashboard/)

### Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/syeo66/go-shufflify.git
    cd go-shufflify
    ```

2. Set up environment variables:

    Create a `.env` file in the root directory with the following content:

    ```env
    SPOTIFY_CLIENT_ID=your_spotify_client_id
    SPOTIFY_CLIENT_SECRET=your_spotify_client_secret
    SESSION_KEY=your_session_key
    PORT=3333
    DB_FILE=./shufflify.db
    ```

3. Install dependencies:

    ```sh
    go mod tidy
    ```

4. Run the application:

    ```sh
    go run server.go
    ```

5. Open your browser and navigate to `http://localhost:3333`.

### Usage

- **Login**: Click the login button to authenticate with your Spotify account.
- **Queue**: View the current queue of tracks.
- **Player**: View the current playing track and player state.

## Contributing

1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Make your changes.
4. Commit your changes (`git commit -am 'Add new feature'`).
5. Push to the branch (`git push origin feature-branch`).
6. Create a new Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgements

- [Spotify Web API](https://developer.spotify.com/documentation/web-api/)
- [Gorilla Sessions](https://github.com/gorilla/sessions)
- [Mattn Go-SQLite3](https://github.com/mattn/go-sqlite3)
