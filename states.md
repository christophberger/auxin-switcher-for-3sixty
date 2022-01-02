---
created: 2022-01-02
tags:
  - state_diagram
  - 3sixty
---

# 3sixty status transitions on music change

## State transition table

| Current 3sixty status   | Music starts           | Music stops        |
| ----------------------- | ---------------------- | ------------------ |
| off                     | on, aux-in             | *no change*        |
| on, aux-in              | *no change*            | dab mode, off      |
| on, switched to aux-in  | *n/a*                  | previous mode, off |
| on, anything but aux-in | on, switched to aux-in | *n/a*              |

## State-event-action-new-state table

| event        | current status          | action                       | new status              |
| ------------ | ----------------------- | ---------------------------- | ----------------------- |
| music starts | off                     | switch on & to aux-in        | on, aux-in              |
| music starts | on, anything but aux-in | switch to aux-in             | on, aux-in              |
| music stops  | on, aux-in              | power off                    | off                     |
| music stops  | on, switched to aux-in  | switch back to previous mode | on, anything but aux-in |

