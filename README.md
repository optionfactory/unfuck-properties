# unfuck-properties

A CLI utility to rescue `.properties` files after your lovely IDE breaks them. 

## Why
Historically, the official Java documentation for the java.util.Properties class explicitly mandated
that .properties files must be encoded in `ISO-8859-1`.

If you wanted to use characters outside of `Latin-1`, the official Java convention 
dictated that you had to use encoded Unicode escape sequences (e.g., writing `\u00FC` instead of `ü`).

IDEs are built to strictly follow Java specifications, they default their text editors to read and write 
.properties files using ISO-8859-1 to ensure compliance with the JVM.

With the release of Java 9, `ResourceBundle` was updated to read .properties files in `UTF-8` by default. 
However, `Properties` still defaults to `ISO-8859-1` for backward compatibility.

This created a massive headache for IDEs. Should they open a `.properties` file as `ISO-8859-1`
or `UTF-8`? Because they can't always guess what version of Java your project targets, they often 
fallback to the safest legal standard: `ISO-8859-1`.

## The Mangle (Mojibake)
When your IDE transparently "unfucks" this or breaks it, it is playing with raw bytes. Here is exactly what happens under the hood when a file gets ruined:
The dev writes `UTF-8`: You or a modern tool saves a file in `UTF-8`. You type the German character `ü`. In UTF-8, `ü` is a 2-byte sequence: `0xC3 0xBC`.
The IDE Opens it as `ISO-8859-1`: The IDE looks at those two bytes (`0xC3` and `0xBC`) but interprets them through the lens of ISO-8859-1, 
where every single byte is its own character.

`0xC3` in ISO-8859-1 is `Ã`

`0xBC` in ISO-8859-1 is `¼`

* The Display: Your IDE screen now literally shows `Ã¼` instead of `ü`.
* The Fatal Overwrite: If you don't notice this and hit save, or if the IDE auto-saves, it commits those interpreted characters back to disk. It writes out the UTF-8 bytes for `Ã` (`0xC3 0x83`) and `¼` (`0xC2 0xBC`).

Your original 2-byte character is now a bloated, corrupted 4-byte disaster. 

This tool essentially treats the corrupted file as a map of misinterpreted ISO bytes, and reverses them back into the proper UTF-8 byte stream.

## Installation

You can download the pre-compiled, statically linked binary for Linux directly from the [Releases](https://github.com/optionfactory/unfuck-properties/releases) page.

```bash
curl -sSL \
  https://github.com/optionfactory/unfuck-properties/releases/latest/download/unfuck-properties-linux-amd64 \
  sudo tee /usr/local/bin/unfuck-properties > /dev/null \
  && sudo chmod +x /usr/local/bin/unfuck-properties
```


## Usage
Pass the target file to the tool:

```bash
./unfuck-properties target-file.properties rewritten-file.properties
```

or just rewrite it in place:

```bash
./unfuck-properties -i target-file.properties
```

## Building from source

If you prefer to build it yourself, you just need Go and Make installed.

```bash
git clone https://github.com/optionfactory/unfuck-properties
cd unfuck-properties
make build
sudo make install
```
