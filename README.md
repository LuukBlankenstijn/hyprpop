# HyprPop

Toggle floating windows in Hyprland with keybindings. Any app or web application is supported.

## Features

- **Instant toggle**: Show/hide windows with custom keybinds
- **Position memory**: Windows remember their size and position between toggles
- **Multi-monitor aware**: Respects monitor scaling and boundaries
- **Special workspace hiding**: Uses Hyprland's special workspaces for clean workspace switching
- **Per-monitor positioning**: Windows adapt to the monitor they're moved to
- **Chromium app mode**: Creates borderless web app windows with isolated profiles

## Installation

### Arch Linux (AUR)

```bash
yay -S hyprpop
# or
paru -S hyprpop
```

### Build from source

```bash
git clone <repository-url>
cd hyprpop
go build -o hyprpop
sudo cp hyprpop /usr/local/bin/
```

## Configuration

Create `~/.config/hypr/hyprpop.yaml`:

```yaml
windows:
  - name: "youtube" # Unique window name
    type: "floating" # Currently only "floating" supported
    startCommand: "https://youtube.com" # Web URL to open or app to run
    keybind: "SUPER+Y" # Hyprland keybind format
    position: [0.1, 0.1] # [x, y] - percentages (0.0-1.0) or pixels
    size: [0.8, 0.8] # [width, height]

  - name: "spotify"
    type: "floating"
    startCommand: "https://open.spotify.com"
    keybind: "SUPER+S" # Modifiers: SUPER, SHIFT, CTRL, ALT
    position: [200, 100] # Absolute pixels
    size: [1200, 800]

  - name: "chatgpt"
    type: "floating"
    startCommand: "https://chat.openai.com"
    keybind: "SUPER+C"
    position: [0.2, 0.1] # Mix percentages and pixels
    size: [0.6, 600]

  - name: "signal"
    type: "floating"
    startCommand: "sigal-desktop"
    keybind: "SUPER+S"
    position: [0.2, 0.1] # Mix percentages and pixels
    size: [0.6, 600]
```

### Config Options

- **name**: Unique identifier (no duplicates allowed)
- **type**: Window type (`"floating"` only for now)
- **startCommand**: Any web URL or command to start an app - if startCommand is an url, creates a browser window with that url - otherwise, starts the app
- **keybind**: Hyprland format - `MOD+key` or `MOD+key`
- **position/size**: `[x, y]` as decimals (0.0-1.0) for percentages or integers for pixels

## Usage

1. Install and configure (see above)

2. Add to Hyprland config to auto-start:
   make sure to put it somewhere in the beginning because it takes some time to start up

```bash
# ~/.config/hypr/hyprland.conf
exec-once = hyprpop
```

Or start manually:

```bash
hyprpop &
```

## How it works

- **First press**: Shows window (creates if needed)
- **Window focused**: Hides to special workspace
- **Window unfocused**: Brings to focus
- **Different workspace**: Moves window to current workspace

Position/size values: use `0.0-1.0` for percentages or integers for pixels.

## Requirements

- Hyprland
- Chromium
- Go 1.24+ (for building)
