# Corbie packs

The packs are ready-made Corbie machines: the listeners and calls of Telegram, Slack, GitLab and Jira, a lookup, an approval and the question it asks, and the starting files of a project, its routers and conversation among them. Each pack is one directory, and Corbie's command `corbie add` writes a pack into a project's `.packs/`. `corbie run submit` brings a project's runs into line with its files, and `corbie ui` shows them.

## Where to start

- `bot/`: the pack every project adds: the conversation's machines, the lookup, the approval, the question and what an agent includes. `bot/pages/` holds the packs' pages: `bot-design.txt` is the design, and `bot-project.txt` how a person uses it. `bot/AGENTS.md` is what a model is told. `README.md` is the front page on GitHub.
- `slack/`, `telegram/`, `gitlab/`, `jira/`: a service's pack, its listener and calls. `shared/`: the pack of scripts the others use, which `corbie add` writes with a pack whose files name `.packs/shared/`.
- `<pack>/starter/`: a pack's starting files, a project's first files, which hold no agent and no toolkit. `BOT` in a path or a text there is the bot's name, which `corbie add --set BOT=<name>` writes, and `<service>/starter/BOT/<service>/` holds a service's router and placed calls.
- `fixture/fixture.go`: what the tests use to make a project from the packs.
- `test/`: the tests. `test/testdata/` holds the agents, toolkit and people the tests add to a project, and `testdata/grouse/`, the worked example as a whole project without `.packs/`, which `grouse_test.go` runs.

## What holds

- **Two ideas**: machine and toolkit. The packs ship flows, statecharts and agents, and files an agent includes (`.packs/bot/agent.yaml`, `.packs/bot/filters/model.yaml`). `people/scopes.yaml` is the one table, and only the lookup reads it. A pack ships no file it reads as configuration. Add no other idea a person must learn.
- **The project's machines say what happens.** A router routes on a message's fields, and the conversation is the project's own file, with its stages written in it. No location is passed in a message's data.
- **A message is one mapping**, the same from every service, given in `bot/pages/bot-design.txt`. A listener only translates between its service and that mapping, and a call takes the service's own input or the shared `{ in, text }`.
- **A project's files hold choices; `.packs/` holds what the packs do.** `corbie add` never updates a project's file, so anything a later version may change goes in a pack, outside its `starter/`. The one exception: `corbie add` adds `.packs/` to `.corbieignore` when it lacks it. A pack's `starter/` holds this version's starting files, and `corbie add` names a project's files that differ from them.
- **Nothing personal lives here.** No real name, address, ID or host goes in this repository. A person's own machines live in their own directory.
- **The packs ship no program.** `corbie add` writes the packs into a project. Everything else is written in machines, which the person can change.
- **A person changes nothing in `.packs/`.** `corbie add` rewrites a pack's directory there.
- **A pack holds no secret.** Every service is reached by the service's machines, whose secrets Corbie keeps at `<bot>/<service>`.
- Corbie stays simple and generic, and gains nothing that serves these packs alone.

## What the code is for

In this order:

1. Simple. Drop anything that can be dropped.
2. One path to do each thing. Flag anything that adds a second.
3. Obvious logic. Prefer the plain way to the clever one.
4. Simple error handling. The layer that notices passes the error up untouched; one outermost place per message, run or command reports everything that went wrong together. Nothing in between catches, retries or rewords.
5. Data flowing through functions, where that is sensible.

And these hold everywhere:

- Nothing important lives in memory. Everything is in Corbie's store or the directory.
- Security rests on a few rules that can be said in a sentence. Do not add a check where there is no real boundary, and do not add fiddly mechanism for a special case.
- If a normal person would need a special config, package or key to set something up, the design is wrong.

## Writing

These rules hold for everything a reader sees: the README, the design, what the bot says and prints, and reports to the author.

Words:

- Use the precise technical word. The developer's word is the plain word: process, file, directory, S3 bucket, network connection, permission, environment variable, default, timeout, poll, encrypted, signed.
- Name the real thing: the field (`allow:`, `job:`), the command, the status (`denied`), the path.
- No metaphor and no soft verbs. Not "decides", "reaches", "touches", "meets", "holds on to", "is owed", "stands", "says so". Not "unless it says": write "by default".
- One word for one thing, the same in the docs, the code and what is printed. A word with two meanings is a fault, and so is a second word for the same thing. A word is changed everywhere or not at all.
- Corbie's words are the packs' words. These are settled:
  - **machine**: a YAML definition of type flow, statechart, agent or decision. Never a computer.
  - **job**: the `job:` block of a stage, a state or a tool, and the `before:` or `after:` of a hook. A job is a script, a machine, or a send. A machine is named by its file (`machine:`) or by its location (`call:`), and a send delivers an event and returns nothing. A runner executes a script; a machine named by its file runs as a child run, and one named by its location as a run of its own whose result is the job's output.
  - **hook**: a job a definition names under `hooks:`, which runs before or after other jobs.
  - **run**: one execution of a machine, which a client starts, lists and shows.
  - **directory**, never "folder". **computer**, never "machine".
  - **controller**, **runner**, **client**: the three kinds of process. **store**: the directory or bucket they share.
- The packs' own words are settled too:
  - **service**: Telegram, Slack, GitLab or Jira, and its directory `<bot>/<service>/`. Never "bridge".
  - **router**, **listener**, **conversation**, **lookup**, **approval**, **question** (a bot's `chat/ask`), **scope**.
  - **toolkit**, never "part"; **description**, never "does"; **listed file**, a file `.corbieignore` lists, never "piece".
- A page does not use words that only the code needs: pass, wake, mark, head, receipt, party. It says what happens instead.

Sentences:

- Few words. Every sentence states a fact the reader needs. If a sentence can be deleted with nothing lost, delete it. Say a thing once.
- Cut every word that carries no fact: "also", "that is", "which is", "in order to", a subject said again, a reason the reader does not need. After writing a paragraph, cut it by a third and check that no fact went.
- A paragraph is at most three sentences. Facts of the same kind go in a list, not in a paragraph.
- Choose the list by what its items are. A term and a few words about it (a field, a flag, a command, a path) is a row of two columns. Steps in order are a numbered list. Anything else is a dash list, and each item of it is a whole sentence with a subject and a verb: "When the script sends the handle to another host, Corbie refuses", never "Sent to another host: Corbie refuses". Do not set sentences in two columns.
- For example, this paragraph:

  "The bucket calls a function when a file is written to that function's inbox directory: the controller's function for the controller's inbox, and a runner's function for that runner's inbox. The function processes everything that is waiting and returns. A schedule also calls the controller's function once a minute."

  is written as:

  ```
    a file written to the controller's inbox   the controller's function
    a file written to a runner's inbox         that runner's function
    a schedule, once a minute                  the controller's function, for timers

  A function handles everything that is waiting, then exits.
  ```
- Say who does what: "The router reads the message, then sends it to the conversation."
- Whole sentences. No slogans, no two-part titles, no comma fragments. A title is a short name.
- A row in a list is one line. A fact added to a row is a few words in that line, not another sentence.
- Never claim more than is true. Check a statement against the code, or run it.
- No diagram drawn in text: no boxes, lines or arrows made of characters. They confuse the author every time. Say the same thing as a short table or a numbered list of one-line steps, and do not let the replacement grow longer than the diagram was.

## Working on the code

- Run Go through mise, which pins the version: `mise exec -- go …`. `mise run check` vets and tests.
- Pure Go, no cgo. Tests need no network, no container, no key and no model: they run a whole Corbie in the test's process with a stand-in model.
- One package per area, named as the area is.
- Until version 1, nothing is kept backwards compatible.

## Docs

Docs say how things are now. No history, no "previously", no alternatives. Something removed is simply gone.

A page of the packs is written as Corbie's are:

- It answers three things: what the thing is, why a reader wants it, and how to do it.
- Its first line is `Title: one sentence that says what the page is for.`, which is its summary in the list of pages.
- It is plain text: a heading is underlined with dashes, a row is indented two spaces, and an example is indented four.
- It gives the definition first, then a real example: the real files, the real command and the output it printed. `pages_test.go` checks every example file and the printed output.
- A mechanism is stated first, in three or four numbered one-line steps, before any example.
- A page of Corbie's is named as `corbie manual <page>`, and a page of the packs by its path in a project, as `.packs/bot/pages/bot-design.txt`. A page ends with a `Reference` list of pages.
- Do not add a page without asking the author. The aim is fewer and shorter pages.

## Commits

- The message starts with the area.
- Never push straight to `main`: push a branch.
