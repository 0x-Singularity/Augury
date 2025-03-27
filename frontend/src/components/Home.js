import React, { useState } from "react";
import Results from "./Results"; 

<img src="/search.svg" alt="Search" />


function Home() {
  const [query, setQuery] = useState("");
  const [results, setResults] = useState(null);
  const [loading, setLoading] = useState(false);

  const handleSearch = async (e) => {
    e.preventDefault();
    if (!query) return;
  
    setLoading(true);
  
    try {
      const response = await fetch("http://localhost:8080/api/ioc/extract", {
        method: "POST",
        headers: {
          "Content-Type": "text/plain",
          "X-User-Name": "demo_user", // You can set this dynamically in future
        },
        body: query,
      });
  
      if (!response.ok) throw new Error("Failed to fetch data");
  
      const data = await response.json();
      setResults(data);
    } catch (error) {
      console.error("Error fetching query:", error);
      setResults({ error: "Failed to fetch results" });
    }
  
    setLoading(false);
  };
  

  return (
    <div className="container" style={{ marginTop: "40px"}}>
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

      {loading && <p>Loading...</p>}

      {results && (
        <div style={{ marginTop: "20px" }}>
          {results && !results.error && (
  <div style={{ marginTop: "20px" }}>
    <h3>Results:</h3>
    <Results results={results} />
  </div>
)}

{results?.error && (
  <p style={{ color: "red", marginTop: "20px" }}>{results.error}</p>
)}

        </div>
      )}
    </div>
  );
}

export default Home;
