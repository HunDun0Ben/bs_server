# JWT Token Generation Tool

This tool is used to generate a pair of Access Token and Refresh Token for testing purposes.

## Usage

```bash
go run app/scripts/jwtscr/generate_jwt_tokens.go [flags]
```

## Flags

| Flag                     | Shorthand | Description                              | Default      |
| ------------------------ | --------- | ---------------------------------------- | ------------ |
| `--expire`               | `-t`      | Access token expiration time in seconds  | 10           |
| `--refresh token expire` | `-f`      | Refresh token expiration time in seconds | 10           |
| `--username`             | `-u`      | Username for the token                   | "username"   |
| `--roles`                | `-r`      | User roles (comma separated)             | "admin,user" |

## Example

Generate tokens for user "admin" with "admin" role and 1-hour expiration:

```bash
go run app/scripts/jwtscr/generate_jwt_tokens.go -u admin -r admin -t 3600 -f 86400
```
