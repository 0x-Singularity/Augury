import React from "react";
import "./OilTable.css";

function IOCTable({ results }) {
  if (!results || !results.data) return <p>No results found.</p>;

  const { data } = results;

  const renderField = (label, value) => {
    if (!value) return null;
    return (
      <div className="field-row">
        <strong>{label}:</strong> <span>{value}</span>
      </div>
    );
  };

  return (
    <div className="results-container">
      {Object.entries(data).map(([source, structures]) => (
        <div key={source} className="ioc-card">
          <h3 className="ioc-header">🔹 Source: {source}</h3>

          {Object.entries(structures).map(([structureType, entries]) => (
            <div key={structureType} className="oil-block">
              <h4 className="source-label">📦 Type: {structureType}</h4>

              {entries.map((entry, idx) => {
                const oil = entry.oil || {};
                const client = entry.client || {};
                const asset = entry.asset || {};
                const binary = entry.binary || {};
                const geo = entry.geo || {};
                const ldap = entry.ldap || {};
                const pdns = entry.pdns || {};
                const process = entry.cbr || {};
                return (
                  <div key={idx} className="structured-entry">
                    {/* Grouped fields by type */}
                    {structureType === "oil" && (
                      <>
                        {renderField("Timestamp", oil.timestamp)}
                        {renderField("User Principal", oil.userPrincipalName)}
                        {renderField("Display Name", oil.displayName)}
                        {renderField("Client IP", oil.clientIp)}
                        {renderField("ASN Org", oil.clientAsOrg)}
                        {renderField("Event Type", oil.eventType)}
                        {renderField("Outcome", oil.outcome)}
                        {renderField("Message", oil.message)}
                      </>
                    )}

                    {structureType === "client" && (
                      <>
                        {renderField("Client IP", client.ip)}
                        {renderField("ASN", client.asn)}
                        {renderField("AS Org", client.as_org)}
                      </>
                    )}

                    {structureType === "asset" && (
                      <>
                        {renderField("Host Name", asset.name)}
                        {renderField("IP", asset.ip)}
                        {renderField("Platform", asset.platformName)}
                        {renderField("Platform Owner", asset.platformOwner)}
                        {renderField("Executive", asset.executive)}
                        {renderField("Stack", asset.stackName)}
                        {renderField("Stack Owner", asset.stackOwner)}
                        {renderField("Created", asset.created)}
                        {renderField("Updated", asset.updated)}
                      </>
                    )}

                    {structureType === "binary" && (
                      <>
                        {renderField("Filename", binary.filename)}
                        {renderField("Accessed", binary.accessed)}
                        {renderField("MD5", binary.md5)}
                        {renderField("SHA256", binary.sha256)}
                        {renderField("URL", binary.url)}
                        {binary.hosts?.length > 0 && (
                          <div className="field-row">
                            <strong>Hosts:</strong> {binary.hosts.join(", ")}
                          </div>
                        )}
                        {renderField("Code Signed", binary.codeSigned ? "Yes" : null)}
                      </>
                    )}

                    {structureType === "geo" && (
                      <>
                        {renderField("City", geo.city)}
                        {renderField("Country Code", geo.countryCode)}
                        {renderField("Country Name", geo.countryName)}
                        {renderField("Latitude", geo.latitude)}
                        {renderField("Longitude", geo.longitude)}
                        {renderField("AS Number", geo.asNumber)}
                        {renderField("AS Org", geo.asOrg)}
                      </>
                    )}

                    {structureType === "ldap" && (
                      <>
                        {renderField("Full Name", ldap.fullName)}
                        {renderField("Name", ldap.name)}
                        {renderField("Title", ldap.title)}
                        {renderField("Email", ldap.email)}
                        {renderField("Phone", ldap.phone)}
                        {renderField("Mobile", ldap.mobile)}
                        {renderField("Created", ldap.created)}
                        {renderField("Manager", ldap.manager)}
                        {renderField("Age", ldap.age)}
                        {renderField("Company", ldap.companyName)}
                      </>
                    )}

                    {structureType === "pdns" && pdns.answers?.length > 0 && (
                      <>
                        <div className="field-row"><strong>DNS Answers:</strong></div>
                        {pdns.answers.map((a, i) => (
                          <div key={i} className="field-row" style={{ marginLeft: "1rem" }}>
                            {renderField("Name", a.name)}
                            {renderField("Type", a.type)}
                            {renderField("Data", a.data)}
                            {renderField("Count", a.count)}
                            {renderField("Start", a.start)}
                            {renderField("End", a.end)}
                          </div>
                        ))}
                      </>
                    )}

                    {structureType === "process" && (
                      <>
                        {renderField("Command Line", process.command_line)}
                        {renderField("Entity ID", process.entity_id)}
                        {renderField("Executable", process.executable)}
                        {renderField("Name", process.name)}
                        {renderField("PID", process.pid)}
                        {renderField("Start Time", process.start)}
                        {renderField("Uptime", process.uptime)}
                    
                        {/* Parent Process Information */}
                        {process.parent && (
                          <>
                            <div className="field-row"><strong>Parent Process:</strong></div>
                            {renderField("Parent Name", process.parent.name)}
                            {renderField("Parent PID", process.parent.pid)}
                            {renderField("Parent Entity ID", process.parent.entity_id)}
                          </>
                        )}
                    
                        {/* User Information */}
                        {process.user && (
                          <>
                            <div className="field-row"><strong>User Information:</strong></div>
                            {renderField("User Name", process.user.name)}
                          </>
                        )}
                    
                        {/* Host Information */}
                        {process.host && (
                          <>
                            <div className="field-row"><strong>Host Information:</strong></div>
                            {renderField("Host Name", process.host.name)}
                            {renderField("Host Type", process.host.type)}
                            {process.host.ip?.length > 0 && (
                              <div className="field-row">
                                <strong>Host IPs:</strong> {process.host.ip.join(", ")}
                              </div>
                            )}
                            {renderField("Host OS", process.host.os?.family)}
                          </>
                        )}
                    
                        {/* Code Signature */}
                        {process.code_signature && (
                          <>
                            <div className="field-row"><strong>Code Signature:</strong></div>
                            {renderField("Code Signed", process.code_signature.exists ? "Yes" : "No")}
                          </>
                        )}
                    
                        {/* Labels */}
                        {process.labels && (
                          <>
                            <div className="field-row"><strong>Labels:</strong></div>
                            {renderField("URL", process.labels.url)}
                          </>
                        )}
                      </>
                    )}
                  </div>
                );
              })}
            </div>
          ))}
        </div>
      ))}
    </div>
  );
}

export default IOCTable;

