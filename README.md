# altstash

**altstash** is wallet app for GNOME, GTK 4, and Libadwaita optimized for a mobile user-experience.

**altstash** works with:
**Taler**,
and more.

## Other Files

* The **user guide** for **altstash** is at: [GUIDE.md](GUIDE.md)
* THe **developer guide** for **altstash** is at: [HACKING.md](HACKING.md)

## Build

### Requirements

- Go >= 1.21
- A C compiler (e.g., GCC) — required for CGo
- GTK 4 >= 4.10 (development headers)
- Libadwaita >= 1.4 (development headers)
- GLib 2.0 (development headers)
- GObject Introspection (development headers)

On Fedora:
```bash
sudo dnf install gcc gtk4-devel libadwaita-devel glib2-devel gobject-introspection-devel
```

On Debian/Ubuntu:
```bash
sudo apt install gcc libgtk-4-dev libadwaita-1-dev libglib2.0-dev libgirepository1.0-dev
```

### Development Build

```bash
go build
```

### Run

```bash
./altstash
```

### Flatpak

```bash
flatpak install --user org.gnome.{Platform,Sdk}//47
flatpak install --user org.freedesktop.Sdk.Extension.golang//24.08
flatpak-builder --user --force-clean --install build build-aux/flatpak/link.reiver.altstash.json
flatpak run link.reiver.altstash
```

## Author

Software **altstash** was written by [Charles Iliya Krempeaux](http://reiver.link)
