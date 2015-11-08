# Cameras

A simple app for interacting with multiple DLink cameras. And probably others. Provides a web interface which renders images from a static list of cameras, and auto refreshes them.

I run this on Heroku, with the camera's being port forwarded.

## Configuration

Set the following env vars:

* `COOKIE_SECRET` - Secret for cookies. A random string
* `MARTINI_ENV=production` enforce SSL etc.
* `WEB_PASSCODE` - Passcode for logging in
* `CAMERAS` - List of cameras, in format `Name;http://user:pass@internetAddress:28001/image.jpg[,..]`

## Screenshots

Look, I can web.

![Desktop Screenshot](https://cdn.lstoll.net/screen/68747470733a2f2f63646e2e6c73746f6c6c2e6e65742f73637265656e2f43616d657261735f323031352d31312d30385f31362d33362d33342e706e67_878_2015-11-09_10-15-43.png)

![Mobile Screenshot](https://cdn.lstoll.net/screen/Screenshot_20151108-163601.png_2015-11-08_16-38-26.png)
