# Go Shufflify

"Go Shufflify" is a web application that integrates with the Spotify API to provide a personalized music queue.

## Features

- **User Authentication**: Authenticate users via Spotify and manage sessions.
- **Queue Management**: Automatically manages the music queue based on user preferences and active playlists.

## Project Structure

- `server.go`: Entry point for the application.
- `types/`: Contains all the data structures used across the application.
- `lib/`: Utility functions and common library code.
- `data/`: Functions for interacting with the Spotify API and the database.
- `routes/`: HTTP route handlers for various endpoints.

## Getting Started

### Prerequisites

- [Go](https://golang.org/dl/) (1.16 or later)
- Spotify Developer Account with Client ID and Client Secret

### Installation

1. Clone the repository:

```bash
git clone https://github.com/yourusername/go-shufflify.git
cd go-shufflify
```

2. Install the required Go packages:

```bash
go mod tidy
```

3. Set up your environment variables:

Create a `.env` file in the root directory with the following content:

```env
SPOTIFY_CLIENT_ID=your_spotify_client_id
SPOTIFY_CLIENT_SECRET=your_spotify_client_secret
SESSION_KEY=your_session_key # A really (!) random key, use 'uuid' or anything else generating random strings
PORT=3333
DB_FILE=./shufflify.db
```

4. Start the application

```bash
./start_server.sh
```

This will create the `shufflify.db` file and set up the necessary tables.

### Running the Application

Start the server:

```bash
./start_server.sh
```

Open your browser and navigate to `http://localhost:3333/`.

## Usage

- **Login**: Navigate to `/login` and log in with your Spotify account.
- **Queue Management**: The application will automatically manage your music queue based on your preferences when 'Shuffle' is enabled.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request or open an issue for any bugs or feature requests.

1. Fork the repository
2. Create a new branch (`git checkout -b feature-branch`)
3. Make your changes
4. Commit your changes (`git commit -am 'Add new feature'`)
5. Push to the branch (`git push origin feature-branch`)
6. Create a new Pull Request

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Acknowledgements

- [Spotify Web API](https://developer.spotify.com/documentation/web-api/)
- [gorilla/sessions](https://github.com/gorilla/sessions)
- [Mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)

---

Happy shuffling! 🎶
