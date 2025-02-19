# Augury

Senior Capstone Project - IOC Intelligence program

Augury is an IOC (Indicators of Compromise) Intelligence application written in Go. It provides a simple web interface where users can paste their IOCs, which are then sent to a locally running Python API for further processing.
The results will be displayed in a table along with links to external vendors for additional intelligence.

## Project Structure

```bash
augury/
├── backend/
│   ├── main.go
│   ├── go.mod
│   └── templates/
│       └── index.html
└── frontend/
    └── static/
        ├── css/
        │   └── styles.css
        ├── js/
        │   └── scripts.js
        └── images/
            └── search.svg

```

- **backend**
  - `main.go`: The main Go application.
  - `templates/`: Contains `index.html` and any additional templates.
  - `go.mod`: Go module file.
- **frontend**
  - `static/css/styles.css`: Global stylesheet (includes Poppins font usage, search bar styling, etc.).
  - `static/js/scripts.js`: JavaScript for front-end interactions.
  - `static/images/search.svg`: Feather Icons search icon.

## Prerequisites

1. **Go** (version 1.18 or higher is recommended)
2. **Python** (for Count Fakeula)
3. **Gorilla Mux** for routing:
   ```bash
   go get github.com/gorilla/mux
   ```
4. **Air** (hot reloading webserver) - Repo: https://github.com/air-verse/air

Using Air for Hot Reloading
Air is a Go development tool that automatically rebuilds and restarts your application when you modify your Go files, saving you from manually restarting the server after every change.

## Install Air

```bash
go install github.com/cosmtrek/air@latest
```

## Initialize Air

In the backend directory, run:

```bash
air init
This will generate an .air.toml configuration file. You can edit this file to customize the watch paths, build options, and other settings.
```

Now run Air

```bash
air
```

Air will now watch your Go files and automatically rebuild and restart your server upon changes. Refresh your browser to see updates.
