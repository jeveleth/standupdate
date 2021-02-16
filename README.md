# standupdate
This cli tool provides a users a way to give formatted standup updates in the standard presentation:

*Yesterday:*
*Today:*
*Blockers:*

To use, set the following envvars:

`export SLACK_TOKEN=yourUserToken`
`export CHANNEL_ID=yourStandupChannel`

Then run `./standupdate report` and follow the prompts.
