import React from "react";

const LOOKUP_LINKS = (ioc) => ({
  PDNS: `insertlink.com`,
  Shodan: `insertlink.com`,
  Censys: `insertlink.com`,
  Spur: `insertlink.com`,
  IP2Proxy: `insertlink.com`,
  BGP: `insertlink.com`,
  OIL: `insertlink.com`,
});

function Results({ results }) {
  if (!results || !results.data) return <p>No results found.</p>;

  const { data } = results;

  return (
    <div className="results-table-wrapper">
      <table className="results-table">
        <thead>
          <tr>
            <th>IOC</th>
            <th>Look Ups</th>
            <th>Asset</th>
            <th>Security Log</th>
            <th>Hash</th>
            <th>Netflow</th>
            <th>Query Log</th>
          </tr>
        </thead>
        <tbody>
          {Object.entries(data).map(([ioc, entry]) => {
            const lookups = LOOKUP_LINKS(ioc);
            const assets = entry.asset?.data || [];
            const securityLogs = entry.coxsight?.data || [];
            const netflows = entry.netflow?.data || [];
            const logs = entry.query_log || [];

            return (
              <tr key={ioc}>
                <td>
                  <a href={`#`} style={{ color: "#fa8b8b", textDecoration: "underline" }}>{ioc}</a>
                </td>
                <td>
                <a href={lookups.PDNS}>PDNS</a><br />
                <a href={lookups.Shodan}>Shodan</a> | <a href={lookups.Censys}>Censys</a><br />
                <a href={lookups.Spur}>Spur</a> | <a href={lookups.IP2Proxy}>IP2Proxy</a> | <a href={lookups.BGP}>BGP View</a>
                </td>

                <td>
                  {assets.length > 0 ? (
                    assets.map((a, i) => (
                      <div key={i}>
                        Host: {a.host?.name}<br />
                        IP: {a.host?.ip}<br />
                        Owner: {a.platform?.owner?.full_name}
                      </div>
                    ))
                  ) : (
                    <em>None</em>
                  )}
                </td>

                <td>
                  {securityLogs.length > 0 ? (
                    securityLogs.map((log, i) => (
                      <div key={i}>
                        User: {log.user?.full_name || log.user?.email}<br />
                        Host: {log.host?.name}<br />
                        Time: {log.timestamp}{"   "}
                        <a
                          href={`/oil?ioc=${encodeURIComponent(ioc)}`}
                          target="_blank"
                          rel="noopener noreferrer"
                        >
                          View OIL
                        </a>
                      </div>
                    ))
                  ) : (
                    <em>None</em>
                  )}
                </td>

                <td>N/A</td>

                <td>
                  {netflows.length > 0 ? (
                    netflows.map((flow, i) => (
                      <div key={i}>
                        Src: {flow.source?.ip}:{flow.source?.port}<br />
                        Dst: {flow.destination?.ip}:{flow.destination?.port}<br />
                        Proto: {flow.network?.transport}
                      </div>
                    ))
                  ) : (
                    <em>None</em>
                  )}
                </td>

                <td>
                  {logs.length > 0 ? (
                    logs.map((log, i) => (
                      <div key={i}>
                        <span style={{ fontSize: "0.85rem", color: "#aaa" }}>•</span>{" "}
                        <strong>{log.last_lookup}</strong>: <strong>{log.user_name}</strong> queried <code>{log.result_count}</code>
                      </div>
                    ))
                  ) : (
                    <em>None</em>
                  )}
                </td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
}

export default Results;
