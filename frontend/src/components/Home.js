import React, { useState } from "react";
import Results from "./Results";
import "./Home.css";

function Home() {
  const [tabs, setTabs] = useState([]);
  const [activeTabId, setActiveTabId] = useState(null);
  const [query, setQuery] = useState("");
  const [loading, setLoading] = useState(false);

  const extractFirstIOC = (data) => {
    const iocs = Object.keys(data?.data || {});
    return iocs.length > 0 ? iocs[0] : "Results";
  };

  const handleSearch = async (e) => {
    e.preventDefault();
    if (!query) return;

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
      const baseLabel = extractFirstIOC(data);

      // Check if this label already exists
      const existingTab = tabs.find((tab) => tab.query === query);
      if (existingTab) {
        setTabs((prevTabs) =>
          prevTabs.map((tab) =>
            tab.id === existingTab.id ? { ...tab, results: data } : tab
          )
        );
        setActiveTabId(existingTab.id);
      } else {
        // Create a new tab with a incrementing label
        const labelCount = tabs.filter((tab) => tab.label.startsWith(baseLabel)).length;
        const label = labelCount > 0 ? `${baseLabel} (${labelCount + 1})` : baseLabel;
        // Create a new tab with a unique ID 
        const newId = Date.now().toString();
        const newTab = { id: newId, label, results: data, query };
        setTabs([...tabs, newTab]);
        setActiveTabId(newId);
      }
    } catch (error) {
      console.error("Error fetching query:", error);
      const errorTab = { error: "Failed to fetch results" };

      const fallbackId = Date.now().toString();
      const fallbackTab = { id: fallbackId, label: "Error", results: errorTab };
      setTabs([...tabs, fallbackTab]);
      setActiveTabId(fallbackId);
    }

    setLoading(false);
  };

  return (
    <div className="container" style={{ marginTop: "40px" }}>
      {/* Search Bar */}
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
          <div key={tab.id} style={{ display: "inline-flex", alignItems: "center" }}>
            <button
              className={`tab-button ${tab.id === activeTabId ? "active" : ""}`}
              onClick={() => {
                setActiveTabId(tab.id);
                setQuery(tab.label === "Untitled" ? "" : tab.label);
              }}
            >
              {tab.label}
            </button>
            <button
              onClick={() => {
                setTabs((prevTabs) => {
                  const updated = prevTabs.filter((t) => t.id !== tab.id);

                  // 🔥 if active tab is closed, activate a neighbor or none
                  if (tab.id === activeTabId) {
                    const next = updated[0]?.id || null;
                    setActiveTabId(next);
                    setQuery(next ? updated[0].label : "");
                  }

                  return updated;
                });
              }}
              style={{
                marginLeft: "-15px",
                border: "none",
                background: "transparent",
                cursor: "pointer",
                fontWeight: "bold",
                color: "#FF0000",
              }}
              title="Close tab"
            >
              ×
            </button>
          </div>
        ))}
      </div>

      {/* Results Section */}
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
                  <p style={{ color: "red", marginTop: "20px" }}>
                    {tab.results.error}
                  </p>
                )}
              </div>
            )
        )}
      </div>
    </div>
  );
}

export default Home;
