# Asteroids Game

A simple Asteroids-inspired game built using [Ebiten](https://ebiten.org/), a 2D game library for Go. Pilot your ship, shoot bullets, and destroy asteroids while avoiding collisions.

## Features
- Rotating and moving ship
- Shooting bullets to destroy asteroids
- Asteroids split into smaller pieces when hit
- Periodic generation of new asteroids
- Score tracking

## Prerequisites
- [Go](https://golang.org/doc/install) installed (1.18 or higher).

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/dhawal-pandya/Asteroids.git
   cd asteroids-game
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

## Running the Game

### On macOS and Linux
1. Open a terminal in the project directory.
2. Run the game:
   ```bash
   go run main.go
   ```

### On Windows
1. Open Command Prompt or PowerShell in the project directory.
2. Run the game:
   ```cmd
   go run main.go
   ```

## Building Executables
To build platform-specific executables:

### macOS
```bash
go build -o asteroids-game
```

### Linux
```bash
go build -o asteroids-game
```

### Windows
```cmd
go build -o asteroids-game.exe
```

The resulting executable can be run directly.

## Controls
- **Arrow Keys**: Rotate and move the ship
- **Spacebar**: Shoot bullets

## Dependencies
- [Ebiten](https://ebiten.org/): 2D game library for Go
- [golang.org/x/image](https://pkg.go.dev/golang.org/x/image): For basic fonts

Install dependencies using:
```bash
go mod tidy
```

## Contribution
Feel free to open issues or submit pull requests to enhance the game!

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
