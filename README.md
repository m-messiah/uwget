# uwget

[![GitHub release](https://img.shields.io/github/release/m-messiah/uwget.svg?style=flat-square)](https://github.com/m-messiah/uwget/releases)
[![Travis](https://img.shields.io/travis/m-messiah/uwget.svg?style=flat-square)](https://travis-ci.org/m-messiah/uwget)

Simple tool for get requests to uwsgi socket (by uwsgi protocol)

## Usage

```bash
uwget [options] uwsgi://host:port/path

Parameters:
  -host string
        HTTP_HOST
  -remote string
        remote addr (default "127.0.0.1")
  -expected-status string
        Fail if response status not equal
  -q    Disable output
```

## Examples

+   Get localhost default uwsgi root

    ```bash
    uwget uwsgi://127.0.0.1:3031
    ```
+   Pass HTTP_HOST header for admin

    ```bash
    uwget --host=example.com uwsgi://backend:3031/admin/
    ```
+   Use another remote address instead of localhost

    ```bash
    uwget -remote=8.8.8.8 uwsgi://localhost:3031/geo/
    ```

+   Ping localhost and fail if not 200

    ```bash
    uwget -q -expected-status=200 uwsgi://127.0.0.1:3031/ping
    ```
