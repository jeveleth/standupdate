# standupdate

This cli tool provides a users a way to give formatted standup updates to a Slack channel in the standard presentation:

    *Name:*
    *Yesterday:*
    *Today:*
    *Blockers:*

## usage

Install the [appropriate version](https://github.com/jeveleth/standupdate/releases) for your OS.

Set the following envvars:

    `export SLACK_TOKEN=yourUserToken`
    `export CHANNEL_ID=yourStandupChannel`

Then run `./standupdate report` and follow the prompts.

### TODOS

* Document functions
* Tests
* CI/CD
* Some DRYing up of structs
