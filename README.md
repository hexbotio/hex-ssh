# Hex Plugin - SSH

Hex Plugin which executes commands on SSH accessible servers.

```
{
  "rule": "example ssh rule",
  "match": "disk size",
  "actions": [
    {
      "type": "hex-ssh",
      "command": "df -h",
      "config": {
        "server": "127.0.0.1",
        "port": "22",
        "login": "hexbot",
        "pass": "${HEXBOT_SSH_PASSWORD}",
        "retry": "2"
      }
    }
  ]
}
```
