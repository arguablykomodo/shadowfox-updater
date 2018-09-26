# Shadowfox Updater

This is a cross-platform installer/uninstaller/updater for [Shadowfox](https://github.com/overdodactyl/ShadowFox), a universal dark theme for Firefox.

## Installing

- For all platforms: go to the [latest release](https://github.com/SrKomodo/shadowfox-updater/releases/latest) and download the respective file for your OS
  - If you are in Linux or Mac, you will probably need to run `chmod +x [filename]` for the OS to register it as an executable
- On Arch Linux you can install the package shadowfox-updater from AUR
- On MacOS, you can install using Homebrew via `brew install srkomodo/tap/shadowfox-updater` and then running the program with `shadowfox`

## How to use

There are various ways to use Shadowfox Updater

### GUI Mode

If you run the file from the command line, it will show you a text-based UI. You can use `TAB` to move between the different options (and `SHIFT+TAB` on some terminals to move backwards), and you can use `ENTER` to toggle the checkboxes and press the buttons.

The "Profile to use" dropdown will show all available profiles in which to install to, and you can cycle through them with the arrow keys.

The "Auto-Generate UUIDs" checkbox, if toggled, will make the updater automatically populate the `internal_UUIDs.txt` file, which is used for styling of extensions. Generally you would toggle this unless you want to manage precisely which extensions get styled.

The "Set Firefox dark theme" checkbox, if toggled, will make the updater automatically enable Firefox's dark theme for it's UI and devtools. If you already have the dark theme enabled, you shouldn't toggle this one.

Then the "Install/Update Shadowfox", "Uninstall Shadowfox" and "Exit" buttons are pretty self explanatory.

#### Fallback

If the text-based UI fails to load (something that happens in some terminals), the program will load a more basic text-only prompt that has the same features as the usual GUI but without the fancy buttons and dropdowns.

### CLI Mode

If you run the file with one or more arguments, the updater will work as a command line tool, which can be useful for automated scripts and such. Instead of explaining how it works I'm just going to paste the result of the command `shadowfox-updater -h`.

```
Usage of shadowfox-updater:
  -generate-uuids
    	Wheter to automatically generate UUIDs or not
  -profile-index int
    	Index of profile to use
  -profile-name string
    	Name of profile to use, if not defined or not found will fallback to profile-index
  -set-dark-theme
    	Wheter to automatically set Firefox's dark theme
  -uninstall
    	Wheter to install or uninstall ShadowFox
```
