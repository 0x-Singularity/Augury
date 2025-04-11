import React from "react";


function OilTable({ results }) {
  if (!results || !results.data) return <p>No results found.</p>;

  const { data} = results;

  return (
    <div className="results-container">
      {Object.entries(data).map(([ioc, oils]) => (
        <div key={ioc} className="ioc-card">
          {/* IOC -> OIL -> Source */}
          {Object.entries(oils).map(([oil, sources]) => (
            <div key={oil} className="oil-block">
              <h3 className="oil-header">🛢️ OIL </h3>
              
              {Object.entries(sources).map(([source, entries]) => (
                <div key={source} className="source-table">
                  <h4 className="source-label">🔹 Source: {source}</h4>
                  <table>
                    <thead>
                      <tr>
                        <th>Timestamp</th>
                        <th>User</th>
                        <th>Display</th>
                        <th>Client IP</th>
                        <th>ASN Org</th>
                      </tr>
                    </thead>
                    <tbody>
                    {Array.isArray(entries) ? (
                      entries.map((entry, idx) => (
                        <tr key={idx}>
                          <td>{entry.timestamp}</td>
                          <td>{entry.userPrincipalName || "-"}</td>
                          <td>{entry.displayName}</td>
                          <td>{entry.client?.ip || "-"}</td>
                          <td>{entry.client?.as_org || "-"}</td>
                        </tr>
                      ))
                    ) : (
                      <tr>
                        <td colSpan={5}><em>No valid data</em></td>
                      </tr>
                    )}
                    </tbody>
                  </table>
                </div>
              ))}
            </div>
          ))}
        </div>
      ))}
    </div>
  );
}

export default OilTable;
