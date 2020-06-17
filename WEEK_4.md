# Lili's Quest Week 4

A roguelike exploration game featuring a little cat, monsters and mutations.

During this week I worked on refactoring the input handling, to make it easier to add new keyboard shortcuts. This has also the advantage, that I can disentangle the sdl events from the actual input handling in the program. Currently only keyboard input is handled like this, but in the future I will also handle mouse input in the same way. In addition I have to think about a system on how I can have different input key bindings/behaviors depending on the state of the game (main menu, options menu, inventory modal open, for example), but this is something I have to think of when I actually have additional game states.

In addition I was not happy with the rendering. I wanted to have the ability to render directly onto a grid and therefore I added a console. For now there is only MatrixConsole which can work with square or rectangle fonts (e.g., 12x12 or 6x12) and form a grid. In the same manner as with libtcod you can then put chars onto those grids and customize their foreground and background color. The game map is using this now which makes rendering much simpler. Due to that refactoring/rewrite I also moved the rendering of the entities into the game map, which makes more sense to me than in the actual game logic.
