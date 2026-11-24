// The approval's script. Its do: is brief, for what the people asked about a
// held job are shown, or answer, for the approval's answer.
//
// They are shown what the agent said as it asked and the job's arguments, then
// who asked and what was asked. What the agent said, what was asked and an
// argument of more than one line are quoted, each line after "> ", so that no
// line of them reads as the bot's own. The job's name is in the question's
// first line.
const quoted = (text) => String(text).split("\n").map((line) => "> " + line).join("\n")
// A value of one line longer than this, such as a link or a file's content,
// is shown as its first and last 150 characters, with its length: whole, it
// can make the question longer than a service takes. Its end is shown, since
// what a command does can be at its end. A value of several lines, such as a
// diff, is shown whole.
const LONG = 300
const cut = (text, at) => text.slice(0, /[\uDC00-\uDFFF]/.test(text[at] || "") ? at - 1 : at)
const short = (value) => value.length > LONG && !value.includes("\n")
  ? cut(value, 150) + " … (" + value.length + " characters, the middle left out) … " + value.slice(cut(value, value.length - 150).length)
  : value
// flat gives each argument as [path, value], a nested key joined by dots.
const flat = (value, path) =>
  value && typeof value === "object" && !Array.isArray(value)
    ? Object.entries(value).flatMap(([key, inner]) => flat(inner, path ? path + "." + key : key))
    : [[path, short(typeof value === "string" ? value : JSON.stringify(value))]]
const shown = (args) => {
  const pairs = flat(args, "").map(([path, value]) => [path ? path + ":" : "", value])
  if (pairs.every(([, value]) => !value.includes("\n"))) return pairs.map((pair) => pair.join(" ").trim()).join(" · ")
  return pairs.map(([label, value]) => value.includes("\n") ? (label + "\n" + quoted(value)).trim() : (label + " " + value).trim()).join("\n")
}
// answer gives the first answer that counts in any place, given the question
// stage's output, an item for each place. With no answer, the answer is no.
if (input.do === "answer") {
  if (!input.asked.length) return { ok: false, note: "No one was asked: the approval's ask names no place." }
  const answer = input.asked.map((p) => p.output).find((o) => o && o.by && !o.elsewhere)
  if (answer) return { ok: answer.choice === "yes", note: answer.by + " said " + answer.choice + "." }
  if (input.asked.every((p) => p.error)) return { ok: false, note: "The question could not be asked: " + input.asked.map((p) => p.error).join("; ") }
  return { ok: false, note: "No one answered in time, so the answer is no." }
}

const tool = input.job.tool || {}
const blocks = []
const lines = []
if (tool.said) lines.push(quoted(tool.said))
const args = shown(tool.args || input.job.input || {})
if (args) lines.push(args)
if (lines.length) blocks.push(lines.join("\n"))
const request = []
if (input.asker) request.push("Asked by " + input.asker + ".")
if (input.request) request.push(quoted(input.request))
if (request.length) blocks.push(request.join("\n"))
return blocks.join("\n\n")
