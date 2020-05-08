# ChainReaction

This repository contains backend API code for the game [Chain Reaction](https://brilliant.org/wiki/chain-reaction-game/) with online multiplayer mode.

endpoints:(both are GET endpoints)
1. /new

query parmas: players_count, dimension (both optional)

2. /join

query params: instance-id (not-optional)


# How to Build
First you will need golang package installed for your respective OS. Refer golang homepage on HowTo.
```
git clone https://github.com/darkLord19/ChainReaction
cd ChainReaction
go build
```
This will generate binary named chainreaction in the directory which you can run as normal unix/windows executable.

# Code overview
- api/ contains logic for endpoints
- simulate/ contains logic for chain reaction game, how to simulate gameboard etc
- datastore/ contains simple in memory mapping of current active games and all games
- game/ contains models required for game instance and other helpers

# Future ideas
1. Build a decent front end to play the game
2. Plan a proper websocket messanging mechanism
3. Refractor code to better fit golang model

# Contribution guidelines
- Anyone who is interested in contributing is welcome. There is no hard requirements. If you are making code better or implementing some new feature, you are most welcome to do so.
