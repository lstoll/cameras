# Cameras

A simple app for interacting with multiple DLink cameras. And probably others. Provides a web interface which renders images from a static list of cameras, and auto refreshes them.

I run this on Heroku, with the camera's being port forwarded.

## Configuration

Set the following env vars:

* `COOKIE_SECRET` - Secret for cookies. A random string
* `MARTINI_ENV=production` enforce SSL etc.
* `WEB_PASSCODE` - Passcode for logging in
* `CAMERAS` - List of cameras, in format `Name;http://user:pass@internetAddress:28001/image.jpg[,..]`
