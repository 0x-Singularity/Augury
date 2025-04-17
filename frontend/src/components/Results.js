import React from "react";

const LOOKUP_LINKS = (ioc) => ({
  PDNS: `/view?source=pdns&ioc=${encodeURIComponent(ioc)}`,
  LDAP: `/view?source=ldap&ioc=${encodeURIComponent(ioc)}`,
  GeoIP: `/view?source=geo&ioc=${encodeURIComponent(ioc)}`,
  Binary: `/view?source=binary&ioc=${encodeURIComponent(ioc)}`,
  Shodan: `https://www.shodan.io/search?query=${encodeURIComponent(ioc)}`,
  VPN: `http://localhost:8080/api/vpn?ioc=${encodeURIComponent(ioc)}`, // Updated to link directly to the raw JSON
  Censys: `https://search.censys.io/search?resource=hosts&q=${encodeURIComponent(ioc)}`,
  Spur: `https://spur.us/search?q=${encodeURIComponent(ioc)}`,
  IP2Proxy: `https://www.ip2proxy.com/demo/${encodeURIComponent(ioc)}`,
  BGP: `https://bgpview.io/ip/${encodeURIComponent(ioc)}`,
  OIL: `/view?source=oil&ioc=${encodeURIComponent(ioc)}`,
  CBR: `/view?source=cbr&ioc=${encodeURIComponent(ioc)}`,
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
            <th>Binary</th>
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
            const binary = entry.binary?.data || [];
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
                  <a href={lookups.PDNS} target="_blank" rel="noopener noreferrer">PDNS</a><br />
                  <a href={lookups.LDAP} target="_blank" rel="noopener noreferrer">LDAP</a><br />
                  <a href={lookups.GeoIP} target="_blank" rel="noopener noreferrer">GeoIP</a><br />
                  <a href={lookups.VPN} target="_blank" rel="noopener noreferrer">VPN</a><br />
                  <a href={lookups.Shodan} target="_blank" rel="noopener noreferrer">Shodan</a> | 
                  <a href={lookups.Censys} target="_blank" rel="noopener noreferrer">Censys</a><br />
                  <a href={lookups.Spur} target="_blank" rel="noopener noreferrer">Spur</a> | 
                  <a href={lookups.IP2Proxy} target="_blank" rel="noopener noreferrer">IP2Proxy</a> | 
                  <a href={lookups.BGP} target="_blank" rel="noopener noreferrer">BGP View</a>
                </td>
                <td>
                  {binary.length > 0 ? (
                    <a href={lookups.Binary} target="_blank" rel="noopener noreferrer">Binary</a>
                  ) : (
                    <em>None</em>
                  )}
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
                          href={LOOKUP_LINKS(ioc).OIL}
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
