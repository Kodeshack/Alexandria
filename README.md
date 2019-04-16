# Alexandria [![Build Status](https://travis-ci.org/Kodeshack/Alexandria.svg?branch=master)](https://travis-ci.org/Kodeshack/Alexandria)
Minimialistic Wiki written in Go

## Config

Configuration for Alexandria is handled using environment variables. The following options are available and will be read at startup time:

- `ALEXANDRIA_TEMPLATE_DIR`: Relative path where Alexandria will look for templates. This should not be changed in most cases.
- `ALEXANDRIA_ASSET_DIR`: Relative path where Alexandria will look for static assets such as JavaScript and CSS files. This should not be changed in most cases.
- `ALEXANDRIA_HOST`: Host on which the HTTP server will listen. Default is `localhost`.
- `ALEXANDRIA_PORT`: Port on which the HTTP server will listen. Default is `:8080`.
    _NOTE_: Must start with a colon `:`.
- `ALEXANDRIA_BASE_URL`: Base URL which will be used to construct all the links and paths for static assets. Default is `http://localhost:8080/`.
    _NOTE_: Must include a trailing slash `/`.
