// The person of input in people/scopes.yaml, as { name, email, scopes }. A row's
// key is an email, matched ignoring case, or a service and the person's ID
// there, as slack:U0123456789, matched as written. The row of the ID is taken
// before the email's. A person the table does not hold gives no scopes and
// the name the service gave, and the log says which keys were looked for. A
// table that is empty is a table with no rows. A row that is not { name,
// scopes } fails every lookup, naming the row, so that a mistake in the table
// is seen at once.
const table = require("../../../people/scopes.yaml") || {}
if (typeof table !== "object" || Array.isArray(table)) throw new Error("people/scopes.yaml is not a mapping of each person to { name, scopes }")
const isRow = (row) => row !== null && typeof row === "object" && !Array.isArray(row) &&
  Object.keys(row).every((k) => k === "name" || k === "scopes") &&
  (row.name === undefined || typeof row.name === "string") &&
  (row.scopes === undefined || Array.isArray(row.scopes) && row.scopes.every((s) => typeof s === "string"))
for (const [key, row] of Object.entries(table)) {
  if (!isRow(row)) throw new Error("people/scopes.yaml: the row of " + key + " is not { name: <text>, scopes: [<scope>, ...] }")
}
const id = input.service && input.user ? input.service + ":" + input.user : ""
const email = input.email.toLowerCase()
const found = id && id in table ? id : email ? Object.keys(table).find((e) => e.toLowerCase() === email) : undefined
if (found === undefined) {
  log("people/scopes.yaml holds no row for " + [id, input.email].filter(Boolean).join(" or ") + ": no scopes")
  return { name: input.name, email: input.email, scopes: [] }
}
const row = table[found]
return { name: row.name || input.name, email: found === id ? input.email : found, scopes: row.scopes || [] }
