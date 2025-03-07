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

1. **Go** (I'm using version 1.23.2)
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
go install github.com/air-verse/air@latest
```

## Initialize Air

In the backend directory, run:

```bash
air init
```

This will generate an .air.toml configuration file. You can edit this file to customize the watch paths, build options, and other settings.

Now run Air

```bash
air
```

Air will now watch your Go files and automatically rebuild and restart your server upon changes. Refresh your browser to see updates.

# How to Launch Augury Locally

Follow these steps to set up and run Augury on your local machine:

1. **Install Go**  
   Make sure [Go](https://go.dev/dl/) is installed on your system.  
   _Verify by running:_

   ```bash
   go version
   ```

2. **Set up the Fakeula API**

- Download the Fakeula zip file from the teams channel
- Install its Python dependencies by running:

```bash
pip install -r requirements.txt
```

3. **Launch the Fakeula API**

- In the Fakeula API directory run:

```bash
python main.py -p 7000
```

This starts the API on port 7000

4. **Start the Augury Backend**

- Navigate to the backend directory of the main Augury project
- Run the following command:

```bash
go run ./main.go
```

This start the Augury backend server

5. **Test the Application**

- Open your browser and go to:

```bash
http:localhost:8080
```

Enter an IOC into the search bar and click the magnifying glass
You should recieve JSON as response
