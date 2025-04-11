import React, { useState } from "react";
import Home from "./components/Home";
import "./App.css";

function App() {
  const [theme, setTheme] = useState("dark");

  const toggleTheme = () => {
    setTheme((prev) => (prev === "dark" ? "light" : "dark"));
  };

  return (
    <div className={`App ${theme}`}>
      <div className="theme-toggle">
        <button onClick={toggleTheme}>
          Switch to {theme === "dark" ? "Light" : "Dark"} Mode
        </button>
      </div>

      <div className="content">
        <h1 style={{ marginBottom: "10px" }}>AUGURY</h1>
        <h2 style={{ marginBottom: "20px" }}>IOC Intelligence</h2>
        <Home theme={theme} />
      </div>
    </div>
  );
}

export default App;
