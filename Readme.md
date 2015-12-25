# ypload

`ypload` is utility for uploading image files to Yandex.Fotki service.

## Usage

    ypload <imageFile>

## Installation

- If you have Go installed (install with `apt-get install golang` for Ubuntu/Debian, `brew install go` with Homebrew on OS X):
0. Make sure you have set `GOPATH` environment variable (to some existing folder, ~/go for example)
1. `go get install ypload`
2. If your `PATH` contains `GOPATH`, then just run as `ypload ...`, otherwise run as `$GOPATH/bin/ypload ...`

- If you don't have (and/or don't want) Go installed: grab binary in releases tab.

## Author

Ivan Zezyulya, ypload@zoid.cc

## License

`ypload` is available under the MIT license. See the LICENSE file for more info.
