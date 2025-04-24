import React from "react";
import "./OilTable.css";

// ────────────────────────────────────────────────────────────────────────────────
// helpers
// ────────────────────────────────────────────────────────────────────────────────

/** Convert snake_case → "Sentence Case" for labels */
const toLabel = (key) =>
  key
    .replace(/_/g, " ")
    .replace(/\b\w/g, (c) => c.toUpperCase());

/**
 * One field line.  When "link" is given, render the value as an <a> tag.
 */
function FieldRow({ label, value, link }) {
  if (
    value === null ||
    value === undefined ||
    (Array.isArray(value) && value.length === 0) ||
    (typeof value === "string" && value.trim() === "")
  )
    return null;

  const printable = Array.isArray(value) ? value.join(", ") : value.toString();
  const content = link ? (
    <a href={link}>{printable}</a>
  ) : (
    <span>{printable}</span>
  );

  return (
    <div className="field-row">
      <strong>{label}: </strong>
      {content}
    </div>
  );
}

// ────────────────────────────────────────────────────────────────────────────────
// main component
// ────────────────────────────────────────────────────────────────────────────────

export default function IOCTable({ results }) {
  if (!results || !results.data) return <p>No results found.</p>;

  return (
    <div className="results-container">
      {Object.entries(results.data).map(([source, structures]) => (
        <div key={source} className="ioc-card">
          <h3 className="ioc-header">🔹 Source: {source}</h3>

          {Object.entries(structures).map(([structureType, entries]) => (
            <div key={structureType} className="oil-block">
              <h4 className="source-label">📦 Type: {structureType}</h4>

              {entries.map((entry, idx) => {
                const record = entry[structureType] ?? {};

                return (
                  <div key={idx} className="structured-entry">
                    {Object.entries(record).map(([k, v]) => {
                      // Special‑case: in a CBR "process" record, make host_name clickable.
                      if (
                        structureType === "process" &&
                        (k === "host_name" || k === "hostname") &&
                        typeof v === "string"
                      ) {
                        return (
                          <FieldRow
                            key={k}
                            label={toLabel(k)}
                            value={v}
                            link={`/view?source=host&ioc=${encodeURIComponent(v)}`}
                          />
                        );
                      }

                      return (
                        <FieldRow key={k} label={toLabel(k)} value={v} />
                      );
                    })}
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

