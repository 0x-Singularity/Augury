import React from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import Home from "./components/Home";
import IOCView from "./components/IOCView";
function App() {
  return (
    <Router>
      <div style={{ fontFamily: "Poppins, sans-serif", padding: "200px", textAlign: "center" }}>
        <h1 style={{ marginBottom: "10px" }}>AUGURY</h1>
        <h2 style={{ marginBottom: "20px" }}>IOC Intelligence</h2>

        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/view" element={<IOCView/>} /> {/* Universal View */}
        </Routes>
      </div>
    </Router>
  );
}

export default App;
