# ChainReaction

This repository contains backend API code for the game [Chain Reaction](https://brilliant.org/wiki/chain-reaction-game/) with online multiplayer mode.

endpoints:(both are GET endpoints)
1. /new

query parmas: players_count, dimension
on success:
```
{game_roomname: "gamblerripple"}
```

2. /games/:rommname/join

query params: username, color
on success:
```
{
    Success: "You have joined the game mothafucka"
    game_instance: {
        all_players: [{…}]
        current_turn: 0
        dimension: 4
        players_count: 2
        room_name: "gamblerripple"
    }
    user: {
        color: "#FF1744"
        username: "test"
    }
```

3. /games/:roomname/colors
on success returns:
```
{
    colors: [...]
}
```

4. /games/:roomname/play

query params: username
on success: doesn't return, establishes websocket connection

# How to Build
First you will need golang package installed for your respective OS. Refer golang homepage on HowTo.
```
git clone https://github.com/darkLord19/chainreaction
cd ChainReaction
go build
```
This will generate binary named chainreaction in the directory which you can run as normal unix/windows executable.

# Code overview
- api/ contains logic for endpoints
- simulate/ contains logic for chain reaction game, how to simulate gameboard etc
- datastore/ contains simple in memory mapping of current active games and all games
- game/ contains models required for game instance and other helpers
- models/ contains types used for game abstraction
- helpers/ contains helper functions for models
- utils/ contains utility function

# Future ideas
1. Build a decent front end to play the game
2. Plan a proper websocket messanging mechanism
3. Refractor code to better fit golang model
4. Add tests

# Contribution guidelines
- Anyone who is interested in contributing is welcome. There is no hard requirements. If you are making code better or implementing some new feature, you are most welcome to do so.
