#### A tray monitor for Concourse-CI

I faced with the problem: `ccmenu` can't show concourse pipelines statuses because it can't proceed with github auth.
This project was a quick research is it hard to write some tool that can show statuses for Concourse CI in OSX tray.

How to use:
 1. connect in teams that you want with the `fly` tool
 2. run this tool

After running, tool will read `~/.flyrc` to get all necessary information.
