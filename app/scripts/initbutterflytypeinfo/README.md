# Butterfly Type Information Initialization Tool

This tool is used to load butterfly classification data from Excel/CSV files and perform initial image resizing.

## Usage

```bash
go run app/scripts/initbutterflytypeinfo/main.go [flags]
```

## Flags

| Flag     | Shorthand | Description                                               | Default                                         |
| -------- | --------- | --------------------------------------------------------- | ----------------------------------------------- |
| `--mode` | `-m`      | Operation mode (see below)                                | `init`                                          |
| `--file` | `-f`      | Path to the Excel or CSV file containing type information | `./scripts/initbutterflytypeinfo/蝴蝶信息.xlsx` |

## Operation Modes

- `init`: Loads butterfly type information (Chinese name, Latin name, English name, and description) from the specified file into the database. It first checks if the data has already been initialized.
- `resize`: Iterates through all image types and resizes their associated images to 200x200 pixels with padding, then stores them in a separate collection.
- `display-one`: Retrieves a single image, resizes it, and displays it in a window (requires GUI support).
- `display-batch`: Retrieves the first 10 resized images from the database and displays them sequentially.

## Example

1. Initialize butterfly types from a different Excel file:

    ```bash
    go run app/scripts/initbutterflytypeinfo/main.go -m init -f ./data/new_butterfly_info.xlsx
    ```

2. Perform batch resizing of all images:
    ```bash
    go run app/scripts/initbutterflytypeinfo/main.go -m resize
    ```
