# Security policy

AegisFlow holds write access to repositories, so we take reports seriously.

## Reporting a vulnerability

Please do **not** open a public issue. Use
[GitHub private vulnerability reporting](https://github.com/rj-karan/AegisFlow/security/advisories/new).

Include: affected component, steps to reproduce, impact, and any proof of concept.
You will get an acknowledgement within 72 hours and a fix or mitigation plan within 14 days
for confirmed issues.

## Scope

- Anything that lets one workspace read or change another workspace's data.
- Webhook forgery or replay.
- Secret leakage into logs, prompts, artifacts or generated workflows.
- A change that passes the gate-weakening detector but weakens CI security.
- Sandbox escape from the validator.

## Supported versions

Only the `main` branch and the latest release receive fixes.
