# Weather CLI tool written in Go

This is a simple command-line tool to fetch weather information using a weather API.

API used : [github.com/chubin/wttr.in](https://github.com/chubin/wttr.in)

- Build it using `make build`
- Run it using `make execute`
- Clean build files using `make clean`
- All three combined `make run`

Optionally, specify the loaction using a command-line argument `LOCATION`.

For example:
```sh
make run LOCATION=mumbai
```

