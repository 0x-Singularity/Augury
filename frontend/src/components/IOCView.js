import { useEffect, useState } from "react";
import IOCTable from "./IOCTable";

function UniversalView() {
  const [results, setResults] = useState(null);
  const [error, setError] = useState(null);

  // Extract source and IOC from the query parameters
  const params = new URLSearchParams(window.location.search);
  const source = params.get("source"); // e.g., "pdns", "ldap", "geo"
  const ioc = params.get("ioc");

  useEffect(() => {
    const fetchData = async () => {
      try {
        if (!source || !ioc) {
          setError("Source and IOC parameters are required.");
          return;
        }

        const res = await fetch(`http://localhost:8080/api/ioc/${source}?ioc=${encodeURIComponent(ioc)}`);
        if (!res.ok) throw new Error(`Request failed: ${res.status}`);
        const data = await res.json();
        console.log(`${source.toUpperCase()} response:`, data); // Debugging
        setResults(data);
      } catch (err) {
        console.error(`Error fetching ${source} data:`, err);
        setError(`Failed to load ${source.toUpperCase()} data.`);
      }
    };

    fetchData();
  }, [source, ioc]);

  if (error) return <p style={{ color: "red" }}>{error}</p>;

  return (
    <div className="container">
      <h2>🔍 {source?.toUpperCase()} Results for: {ioc}</h2>
      <IOCTable results={results} />
    </div>
  );
}

export default UniversalView;