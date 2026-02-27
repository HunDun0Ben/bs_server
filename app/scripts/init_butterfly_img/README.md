# Butterfly Image Initialization Tool

This tool handles the initial processing and database insertion of butterfly images and their segmentations. It also supports feature extraction and clustering.

## Usage

```bash
go run app/scripts/init_butterfly_img/main.go [flags]
```

## Flags

| Flag           | Shorthand | Description                                                                     | Default                               |
| -------------- | --------- | ------------------------------------------------------------------------------- | ------------------------------------- |
| `--path`       | `-p`      | Base path for butterfly dataset (contains `images` and `segmentations` folders) | `/home/workspace/data/leedsbutterfly` |
| `--mode`       | `-m`      | Operation mode (see below)                                                      | `verify`                              |
| `--clusters`   | `-k`      | Number of clusters (k) for KMeans                                               | 1024                                  |
| `--iterations` | `-i`      | Maximum iterations for KMeans                                                   | 10                                    |

## Operation Modes

- `insert`: Scans the dataset path and inserts images with their corresponding masks into the database.
- `verify`: Checks if all images in the `images` folder have a corresponding mask in the `segmentations` folder.
- `display`: Retrieves a sample image from the database and displays it in a window (requires GUI support).
- `shift`: Extracts SIFT features from resized images and updates them in the database.
- `kmeans`: Performs KMeans clustering on all extracted features and generates a Bag-of-Words (BoW) training dataset (`data.csv`).

## Example

1. Verify the dataset:

    ```bash
    go run app/scripts/init_butterfly_img/main.go -m verify -p /path/to/data
    ```

2. Run KMeans with 512 clusters:
    ```bash
    go run app/scripts/init_butterfly_img/main.go -m kmeans -k 512
    ```
