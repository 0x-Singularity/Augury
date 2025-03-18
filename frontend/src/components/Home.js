import React, { useState } from "react";
<img src="/search.svg" alt="Search" />


function Home() {
  const [query, setQuery] = useState("");
  const [results, setResults] = useState(null);
  const [loading, setLoading] = useState(false);

  const handleSearch = async (e) => {
    e.preventDefault(); // Prevent form submission from reloading the page
    if (!query) return;
    
    setLoading(true);
    
    try {
      const response = await fetch(`http://localhost:8080/api/ioc/fakeula?ioc=${query}`);
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
    <div className="container">
      <h1>Search</h1>

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
          <h3>Results:</h3>
          <pre>{JSON.stringify(results, null, 2)}</pre>
        </div>
      )}
    </div>
  );
}

export default Home;
