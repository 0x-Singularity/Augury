import React, { useState, useEffect } from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import Home from "./components/Home";
import Login from "./components/Login";
import IOCView from "./components/IOCView";
import "./App.css";

function App() {
  /**
   * ----------------------
   *  Local state & prefs
   * ----------------------
   */
  // user login
  const [username, setUsername] = useState("");
  const [showChangeButton, setShowChangeButton] = useState(false);

  // theme (dark | light)
  const [theme, setTheme] = useState("dark");

  /**
   * Hydrate state from localStorage once, on mount
   */
  useEffect(() => {
    setUsername(localStorage.getItem("username") ?? "");
    setTheme(localStorage.getItem("theme") ?? "dark");
  }, []);

  /**
   * Handlers
   */
  const handleUsernameChange = (newUsername) => {
    setUsername(newUsername);
    if (newUsername) {
      localStorage.setItem("username", newUsername);
    } else {
      localStorage.removeItem("username");
    }
  };

  const toggleChangeButton = () => setShowChangeButton((prev) => !prev);

  const toggleTheme = () => {
    setTheme((prev) => {
      const next = prev === "dark" ? "light" : "dark";
      localStorage.setItem("theme", next);
      return next;
    });
  };

  /**
   * ----------------------
   *        Render
   * ----------------------
   */
  return (
    <Router>
      <div className={`App ${theme}`} style={{ fontFamily: "Poppins, sans-serif", minHeight: "100vh" }}>
        {/* Theme toggle */}
        <div className="theme-toggle" style={{ position: "absolute", top: "10px", left: "20px" }}>
          <button onClick={toggleTheme}>
            Switch to {theme === "dark" ? "Light" : "Dark"} Mode
          </button>
        </div>

        {/* User display */}
        <div style={{ position: "absolute", top: "10px", right: "20px", textAlign: "right" }}>
          <span
            onClick={toggleChangeButton}
            style={{
              cursor: "pointer",
              color: showChangeButton ? "#89b4fa" : "white",
              textDecoration: showChangeButton ? "underline" : "none",
              transition: "color 0.3s ease, text-decoration 0.3s ease",
            }}
            onMouseEnter={(e) => !showChangeButton && (e.target.style.color = "#89b4fa")}
            onMouseLeave={(e) => !showChangeButton && (e.target.style.color = "white")}
          >
            {username ? (
              <>
                Logged in as: <strong>{username}</strong>
              </>
            ) : (
              "Not logged in"
            )}
          </span>
          {showChangeButton && (
            <button
              onClick={() => handleUsernameChange("")}
              className="user-button"
              style={{
                marginTop: "5px",
                padding: "5px 10px",
                fontSize: "0.9rem",
                cursor: "pointer",
                display: "block",
                marginLeft: "auto",
              }}
            >
              Change Username
            </button>
          )}
        </div>

        {/* Main content */}
        <div style={{ padding: "200px", textAlign: "center" }}>
          <h1 style={{ marginBottom: "10px" }}>AUGURY</h1>
          <h2 style={{ marginBottom: "20px" }}>IOC Intelligence</h2>

          {!username ? (
            <Login onLogin={handleUsernameChange} />
          ) : (
            <Routes>
              <Route path="/" element={<Home />} />
              <Route path="/view" element={<IOCView />} />
            </Routes>
          )}
        </div>
      </div>
    </Router>
  );
}

export default App;
