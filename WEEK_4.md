# Lala's Quest Week 4

A roguelike exploration game featuring a little cat, monsters and mutations.

During this week I worked on refactoring the input handling, to make it easier to add new keyboard shortcuts. This has also the advantage, that I can disentangle the sdl events from the actual input handling in the program. Currently only keyboard input is handled like this, but in the future I will also handle mouse input in the same way. In addition I have to think about a system on how I can have different input key bindings/behaviors depending on the state of the game (main menu, options menu, inventory modal open, for example), but this is something I have to think of when I actually have additional game states.
