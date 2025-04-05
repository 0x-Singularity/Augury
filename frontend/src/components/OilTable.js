import React from "react";


function OilTable({ results }) {
  if (!results || !results.data) return <p>No results found.</p>;

  const { data, queryLogs } = results;

  return (
    <div className="results-container">
      {Object.entries(data).map(([ioc, oils]) => (
        <div key={ioc} className="ioc-card">
          <h2 className="ioc-header">IOC: {ioc}</h2>

          {/* Query Logs */}
          <div className="query-logs">
            <h3>Query Logs</h3>
            <ul>
              {queryLogs[ioc]?.length > 0 ? (
                queryLogs[ioc].map((log) => (
                  <li key={log.log_id}>
                    <span className="timestamp">[{log.last_lookup}]</span> —{" "}
                    <strong>{log.user_name}</strong> queried {log.result_count} result(s)
                  </li>
                ))
              ) : (
                <li>No logs found</li>
              )}
            </ul>
          </div>

          {/* IOC -> OIL -> Source */}
          {Object.entries(oils).map(([oil, sources]) => (
            <div key={oil} className="oil-block">
              <h3 className="oil-header">🛢️ OIL: <strong>{oil}</strong></h3>
              
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
                      {entries.map((entry, idx) => (
                        <tr key={idx}>
                          <td>{entry.timestamp}</td>
                          <td>{entry.userPrincipalName || "-"}</td>
                          <td>{entry.displayName}</td>
                          <td>{entry.client?.ip || "-"}</td>
                          <td>{entry.client?.as_org || "-"}</td>
                        </tr>
                      ))}
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
