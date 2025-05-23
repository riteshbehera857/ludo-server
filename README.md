# game-server-go

# game-server-go

- [x] All the message classes, interfaces for a respective classes should come under that specific package/folder

- [x] All file names must match the package name of their respective folder.

- [ ] Remove the top-level game folder.

- [x] If a folder exists, it must contain only files, not sub-folders.

  Example:
  In the Board folder:
  Board/
  |_ Message (Folder) ❌  
   |_ Message.go ✔

- [x] Replace game.ended with game.end; game.ended should no longer exist.

- [x] Rename select_quadrant_message.go to select_quadrant_options_message.go.

<!-- - [ ] There will be no dice_rolling_message.go. (DOUBT) -->

- [x] Remove the game_manager.go file.

- [x] The mapping between boardId and BoardInstances should be moved to the ludo_game_service.go file.

- [x] When a game finishes, the status of the corresponding Board instance should be set to FINISHED.

- [x] The Board instance does not need to inform the Ludo game service about game.ended. Instead, implement a periodic check to determine the game’s status.

- [x] Each Quadrant will have an isOccupied: bool property.

- [x] When a Quadrant is occupied, set isOccupied: true. The Board will only send Quadrants where isOccupied is false.

- [x] Remove all occurrences of room and roomId. Rename them to board and boardId, respectively.

- [x] Move the SelectQuadrant logic to the board package.

- [x] The SelectQuadrant logic will:
      Receive the message.
      Locate the relevant Quadrant class.
      Call the appropriate method on that class.

- [x] Use ludo_games as the name for the database collection instead of games.

- [x] Eliminate all references to "game" terminology; use "board" instead throughout the project.
#   l u d o - s e r v e r  
 