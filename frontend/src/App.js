import React, { useState, useEffect } from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import Home from "./components/Home";
import Login from "./components/Login";
import IOCView from "./components/IOCView";

function App() {
  const[username, setUsername] = useState("");
  const [showChangeButton, setShowChangeButton] = useState(false); //toggle the change username button

    useEffect(() => {
      // Load the username from localStorage when the app starts
      const storedUsername = localStorage.getItem("username");

      if (storedUsername) {
        setUsername(storedUsername);
      }
    }, []);

    const handleUsernameChange = (newUsername) => {
      setUsername(newUsername);
      localStorage.setItem("username", newUsername); // Save the new username to localStorage
    };

    const toggleChangeButton = () => {
      setShowChangeButton((prev) => !prev); // Toggle the visibility of the button
    };
  return (
  
    <Router>
      <div style={{ fontFamily: "Poppins, sans-serif", padding: "200px", textAlign: "center" }}>
      <div style={{ position: "absolute", top: "10px", right: "20px" }}>
      <span
            onClick={toggleChangeButton} // Toggle the dropdown on click
            style={{
              cursor: "pointer",
              color: showChangeButton ? "#89b4fa" : "white", // Change color based on button visibility
              textDecoration: showChangeButton ? "underline" : "none", // Underline when button is visible
              transition: "color 0.3s ease, text-decoration 0.3s ease", 
            }}
            onMouseEnter={(e) => {
              if (!showChangeButton) {
                e.target.style.color = "#89b4fa"; // Change color on hover
              }
            }}
            onMouseLeave={(e) => {
              if (!showChangeButton) {
                e.target.style.color = "white"; // Revert color on mouse leave
              }
            }}
          >
            Logged in as: <strong>{username}</strong>
          </span>
          {showChangeButton && ( // Conditionally render the "Change Username" button
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

        <h1 style={{ marginBottom: "10px" }}>AUGURY</h1>
        <h2 style={{ marginBottom: "20px" }}>IOC Intelligence</h2>

        {!username ? (
          <Login onLogin={handleUsernameChange} />
        ) : (
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/view" element={<IOCView/>} /> {/* Universal View */}
        </Routes>
        )}
      </div>
    </Router>
  );
}

export default App;
