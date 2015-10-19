
ESCAPE FROM PROXIMA 5
---------------------

Version 0.1

(For the in-game backstory, see story.txt.)

To build it:  ./rebuild.sh

THIS GAME IS A WORK IN PROGRESS; PROBABLY OF LITTLE ACTUAL PLAYABLE VALUE
OTHER THAN AS A NOVELTY.

Only a single game level is implemented so far.

I built this intending to show off the capabilities of my Tcell package for
Go (the library used for "terminal" graphics under the hood.)

I got kind of caught up in it, and now there is a pretty powerful game engine
under the hood.  There is full support for animated sprites, with collision
detection, layering, etc.  And I've made some pretty cool visual effects too.

It runs in a text window.  You want a reasonably big window (80x25 will work,
but its harder to play with smaller screens), and you'll want a fast display.

You'll also want a color terminal -- 256 color support in xterm or Terminal
is highly recommended.  Oh, and if you can, try to use a UTF-8 locale.  The
game will work in basic ASCII with black and white, and it isn't completely
garbage in that mode, but its *soo* much better with a more capable terminal.

Oh, don't try to run this over a 300 baud modem.  You won't be happy.

Local terminals work fine.  9600 baud is probably pushing it. I haven't
tested anything other than local connections; there are a lot of display
updates in the game, so you do want to have a reasonably fast display.
You can reduce the screen size if its too slow, that may help a little.

Much of the code here will probably get cleaned up and structured properly
for reuse in nice library form.  Probably sub-packages under Tcell, for both
the dynamic views and probably also the sprite management code and maybe even
the full game loop.
