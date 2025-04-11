import { useEffect, useState } from "react";
import OilTable from "./OilTable";

function OilView() {
  const [results, setResults] = useState(null);
  const ioc = new URLSearchParams(window.location.search).get("ioc");

useEffect(() => {
  const fetchOil = async () => {
    try {
      const res = await fetch(`http://localhost:8080/api/ioc/oil?ioc=${encodeURIComponent(ioc)}`);
      if (!res.ok) throw new Error(`Request failed: ${res.status}`);
      const data = await res.json();
      console.log("OIL response:", data); // 🔍 inspect this!
      setResults(data);
    } catch (err) {
      console.error("Error fetching OIL data:", err);
      setResults({ error: "Failed to load OIL data" });
    }
  };

  if (ioc) fetchOil();
}, [ioc]);
  return (
    <div className="container">
      <h2>🛢️ OIL Results for: {ioc}</h2>
      <OilTable results={results} />
    </div>
  );
}

export default OilView;
