## Use

`kure import <manager-name> [-e erase] [-p path]`

## Description

Import entries from other password managers. Format: CSV.

If an entry already exists it will be overwritten.

Delete the CSV used with the `erase` flag, the file will be deleted only if no errors were encountered.

Password managers supported:
- 1Password
- Bitwarden
- Keepass/X/XC
- Lastpass

## Flags

|  Name     | Shorthand |     Type      |    Default    |                   Description                     |
|-----------|-----------|---------------|---------------|---------------------------------------------------|
| erase     | e         | bool          | false         | Erase file on exit (only if there are no errors)  |
| path      | p         | string        | ""            | Source file path                                  |

### Examples

Import:
```
kure import <manager-name> -p path/to/file
```

Import and erase the file:
```
kure import <manager-name> -e -p path/to/file
```