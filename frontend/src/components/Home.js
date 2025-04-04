import React, { useState } from "react";
import Results from "./Results";
import "./Home.css";

function Home() {
  const [tabs, setTabs] = useState([]);
  const [activeTabId, setActiveTabId] = useState(null);
  const [query, setQuery] = useState("");
  const [loading, setLoading] = useState(false);

  const handleNewTab = () => {
    const newId = Date.now().toString();
    const newTab = { id: newId, label: "Untitled", results: null };
    setTabs([...tabs, newTab]);
    setActiveTabId(newId);
    setQuery(""); // Reset input
  };

  const handleSearch = async (e) => {
    e.preventDefault();
    if (!query || !activeTabId) return;

    setLoading(true);

    try {
      const response = await fetch("http://localhost:8080/api/ioc/extract", {
        method: "POST",
        headers: {
          "Content-Type": "text/plain",
          "X-User-Name": "demo_user",
        },
        body: query,
      });

      if (!response.ok) throw new Error("Failed to fetch data");

      const data = await response.json();

      // Update tab with results + label
      setTabs((prevTabs) =>
        prevTabs.map((tab) =>
          tab.id === activeTabId ? { ...tab, results: data, label: query } : tab
        )
      );
    } catch (error) {
      console.error("Error fetching query:", error);
      setTabs((prevTabs) =>
        prevTabs.map((tab) =>
          tab.id === activeTabId
            ? { ...tab, results: { error: "Failed to fetch results" } }
            : tab
        )
      );
    }

    setLoading(false);
  };

  return (
    <div className="container" style={{ marginTop: "40px" }}>
      {/* Shared search bar */}
      <div className="search-box">
        <form onSubmit={handleSearch}>
          <input
            type="text"
            placeholder="Enter search query..."
            value={query}
            onChange={(e) => setQuery(e.target.value)}
          />
          <button type="submit">
            <img src="/search.svg" alt="Search" />
          </button>
        </form>
      </div>

      {/* Tab buttons */}
      <div className="tab-buttons" style={{ marginTop: "20px" }}>
        {tabs.map((tab) => (
          <button
            key={tab.id}
            className={`tab-button ${tab.id === activeTabId ? "active" : ""}`}
            onClick={() => {
              setActiveTabId(tab.id);
              setQuery(tab.label === "Untitled" ? "" : tab.label); // repopulate input
            }}
          >
            {tab.label}
          </button>
        ))}
        <button className="tab-button new" onClick={handleNewTab}>
          ➕ New Tab
        </button>
      </div>

      {/* Results section */}
      <div className="tab-content">
        {tabs.map(
          (tab) =>
            tab.id === activeTabId && (
              <div key={tab.id} className="tab-panel">
                {loading && <p>Loading...</p>}

                {tab.results && !tab.results.error && (
                  <div style={{ marginTop: "20px" }}>
                    <h3>Results:</h3>
                    <Results results={tab.results} />
                  </div>
                )}

                {tab.results?.error && (
                  <p style={{ color: "red", marginTop: "20px" }}>{tab.results.error}</p>
                )}
              </div>
            )
        )}
      </div>
    </div>
  );
}

export default Home;