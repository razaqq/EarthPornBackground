# EarthPornBackground

Fetches one of the top trending images of the day from reddit.com/r/EarthPorn 
and sets it as your desktop background on each login.
Image resolution has to be at least monitor resolution and aspect ratio has to be within 10% of the monitor aspect ratio.

### Installing
1. Move the binary anywhere 
2. Create a shortcut to the binary
3. Move the shortcut to `%APPDATA%\Microsoft\Windows\Start Menu\Programs\Startup`

### Compiling
`go build -ldflags -H=windowsgui`
