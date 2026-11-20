# Rook

Rook is ready-made [Corbie](https://github.com/zummed/corbie) machines for a bot that answers people on Telegram, Slack, GitLab issues and Jira issues, with agents, toolkits and approvals. Rook is a set of packs, each a directory of this repository. `corbie add` writes a pack into a project's `.packs/`, Corbie runs its machines, and the project's own files say where each message goes, who someone is, and what each agent can do.

| Pack | What it holds |
|---|---|
| `bot` | the conversation, the lookup, the approval and the question it asks |
| `slack`, `telegram`, `gitlab`, `jira` | a service's listener and calls |
| `shared` | scripts the others use; `corbie add` writes it with them |

## Start

Install [Corbie](https://github.com/zummed/corbie). Make a project with the pack `bot` and a service's pack, naming the bot, and read its page:

```sh
corbie init --ceiling all
mkdir my-bot && cd my-bot
corbie add github.com/zummed/corbie-packs/bot github.com/zummed/corbie-packs/slack --set BOT=grouse
less .packs/bot/pages/rook-project.txt
```

`corbie add` writes two pages of Rook's into the project:

| Page | What it gives |
|---|---|
| `.packs/bot/pages/rook-project.txt` | making a project, what to run after a change, and the common changes |
| `.packs/bot/pages/rook-design.txt` | the runs that answer a message, what each machine takes and gives, and who may do what |

A pack's starting files are under its `starter/` directory, and `corbie add` with no pack named adds a project's packs again from where they came.

`corbie run submit .` in the project records it in Corbie as one tree, `/my-bot`, and starts a router for each service. A router's listener hears a message, and the router sends it to the conversation for that person and thread, which looks the person up and calls the agent with the toolkits their scopes allow. Each run names the run whose send or call started it, so a message can be followed to its agent, tools and approvals.

A new project has no agent and no toolkit, and Rook needs no model: an agent names its own, or is a machine with none. To have a model add them, start Claude Code or another model in the project and say what the bot should do. The project's `.mcp.json` gives it Corbie's manual and `corbie test`, and its `AGENTS.md` tells any model how the project is changed.

## License

Copyright 2026 Paul Young. Licensed under the Apache License, Version 2.0: see `LICENSE` and `NOTICE`.

Rook was written with Claude Code, designed and directed by Paul Young.
