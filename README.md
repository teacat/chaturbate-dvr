# Chaturbate DVR

A tool to record **multiple** Chaturbate streams. Supports macOS, Windows, Linux, and Docker. Favicon from [Twemoji](https://github.com/twitter/twemoji).

![Image](https://github.com/user-attachments/assets/d71f0aaa-e821-4371-9f48-658a137b42b6)

![Image](https://github.com/user-attachments/assets/43ab0a07-0ece-40ba-9a0f-045ca0316638)

&nbsp;

# Getting Started

Go to the [📦 Releases page](https://github.com/teacat/chaturbate-dvr/releases) and download the appropriate binary. (e.g., `x64_windows_chaturbate-dvr.exe`)

&nbsp;

## 🌐 Launching the Web UI

```bash
# Windows
$ x64_windows_chaturbate-dvr.exe

# macOS / Linux
$ ./x64_linux_chaturbate-dvr
```

Then visit [`http://localhost:8080`](http://localhost:8080) in your browser.

&nbsp;

## 💻 Using as a CLI Tool

```bash
# Windows
$ x64_windows_chaturbate-dvr.exe -u CHANNEL_USERNAME

# macOS / Linux
$ ./x64_linux_chaturbate-dvr -u CHANNEL_USERNAME
```

This starts recording immediately. The Web UI will be disabled.

&nbsp;

## 🐳 Running with Docker

Pre-built image `yamiodymel/chaturbate-dvr` from [Docker Hub](https://hub.docker.com/r/yamiodymel/chaturbate-dvr):

```bash
# Run the container and save videos to ./videos
$ docker run -d \
    --name my-dvr \
    -p 8080:8080 \
    -v "./videos:/usr/src/app/videos" \
    -v "./conf:/usr/src/app/conf" \
    yamiodymel/chaturbate-dvr
```

...Or build your own image using the Dockerfile in this repository.

```bash
# Build the image
$ docker build -t chaturbate-dvr .

# Run the container and save videos to ./videos
$ docker run -d \
    --name my-dvr \
    -p 8080:8080 \
    -v "./videos:/usr/src/app/videos" \
    -v "./conf:/usr/src/app/conf" \
    chaturbate-dvr
```

...Or use [`docker-compose.yml`](https://github.com/teacat/chaturbate-dvr/blob/master/docker-compose.yml):

```bash
$ docker-compose up
```

Then visit [`http://localhost:8080`](http://localhost:8080) in your browser.

&nbsp;

# 🧾 Command-Line Options

Available options:

```
--username value, -u value  The username of the channel to record
--admin-username value      Username for web authentication (optional)
--admin-password value      Password for web authentication (optional)
--framerate value           Desired framerate (FPS) (default: 30)
--resolution value          Desired resolution (e.g., 1080 for 1080p) (default: 1080)
--pattern value             Template for naming recorded videos (default: "videos/{{.Username}}_{{.Year}}-{{.Month}}-{{.Day}}_{{.Hour}}-{{.Minute}}-{{.Second}}{{if .Sequence}}_{{.Sequence}}{{end}}")
--max-duration value        Split video into segments every N minutes ('0' to disable) (default: 0)
--max-filesize value        Split video into segments every N MB ('0' to disable) (default: 0)
--port value, -p value      Port for the web interface and API (default: "8080")
--interval value            Check if the channel is online every N minutes (default: 1)
--cookies value             Cookies to use in the request (format: key=value; key2=value2)
--user-agent value          Custom User-Agent for the request
--domain value              Chaturbate domain to use (default: "https://chaturbate.global/")
--help, -h                  show help
--version, -v               print the version
```

**Examples**:

```bash
# Record at 720p / 60fps
$ ./chaturbate-dvr -u yamiodymel -resolution 720 -framerate 60

# Split every 30 minutes
$ ./chaturbate-dvr -u yamiodymel -max-duration 30

# Split at 1024 MB
$ ./chaturbate-dvr -u yamiodymel -max-filesize 1024

# Custom filename format
$ ./chaturbate-dvr -u yamiodymel \
    -pattern "video/{{.Username}}/{{.Year}}-{{.Month}}-{{.Day}}_{{.Hour}}-{{.Minute}}-{{.Second}}_{{.Sequence}}"
```

_Note: In Web UI mode, these flags serve as default values for new channels._

&nbsp;

# 🍪 Cookies & User-Agent

You can set Cookies and User-Agent via the Web UI or command-line arguments.

![localhost_8080_ (4)](https://github.com/user-attachments/assets/cbd859a9-4255-404b-b6bf-fa89342f7258)

_Note: Use semicolons to separate multiple cookies, e.g., `key1=value1; key2=value2`._

&nbsp;

## ☁️ Bypass Cloudflare

1. Open [Chaturbate](https://chaturbate.com) in your browser and complete the Cloudflare check.

    (Keep refresh with F5 if the check doesn't appear)

2. **DevTools (F12)** → **Application** → **Cookies** → `https://chaturbate.com` → Copy the `cf_clearance` value

![sshot-2025-04-30-146](https://github.com/user-attachments/assets/69f4061b-29a2-48a7-ad57-0c86148805e2)

3. User-Agent can be found using [WhatIsMyBrowser](https://www.whatismybrowser.com/detect/what-is-my-user-agent/), now run with `-cookies` and `-user-agent`:

    ```bash
    $ ./chaturbate-dvr -u yamiodymel \
        -cookies "cf_clearance=PASTE_YOUR_CF_CLEARANCE_HERE" \
        -user-agent "PASTE_YOUR_USER_AGENT_HERE"
    ```

    Example:

    ```bash
    $ ./chaturbate-dvr -u yamiodymel \
        -cookies "cf_clearance=i975JyJSMZUuEj2kIqfaClPB2dLomx3.iYo6RO1IIRg-1746019135-1.2.1.1-2CX..." \
        -user-agent "Mozilla/5.0 (Windows NT 10.0; Win64; x64)..."
    ```

&nbsp;

## 🕵️ Record Private Shows

1. Login [Chaturbate](https://chaturbate.com) in your browser.

2. **DevTools (F12)** → **Application** → **Cookies** → `https://chaturbate.com` → Copy the `sessionid` value

3. Run with `-cookies`:

    ```bash
    $ ./chaturbate-dvr -u yamiodymel -cookies "sessionid=PASTE_YOUR_SESSIONID_HERE"
    ```

&nbsp;

# 📄 Filename Pattern

The format is based on [Go Template Syntax](https://pkg.go.dev/text/template), available variables are:

`{{.Username}}`, `{{.Year}}`, `{{.Month}}`, `{{.Day}}`, `{{.Hour}}`, `{{.Minute}}`, `{{.Second}}`, `{{.Sequence}}`

&nbsp;

Default it hides the sequence if it's zero.

```
Pattern: {{.Username}}_{{.Year}}-{{.Month}}-{{.Day}}_{{.Hour}}-{{.Minute}}-{{.Second}}{{if .Sequence}}_{{.Sequence}}{{end}}
 Output: yamiodymel_2024-01-02_13-45-00.ts    # Sequence won't be shown if it's zero.
 Output: yamiodymel_2024-01-02_13-45-00_1.ts
```

**👀 or... The sequence can be shown even if it's zero.**

```
Pattern: {{.Username}}_{{.Year}}-{{.Month}}-{{.Day}}_{{.Hour}}-{{.Minute}}-{{.Second}}_{{.Sequence}}
 Output: yamiodymel_2024-01-02_13-45-00_0.ts
 Output: yamiodymel_2024-01-02_13-45-00_1.ts
```

**📁 or... Folder per each channel.**

```
Pattern: video/{{.Username}}/{{.Year}}-{{.Month}}-{{.Day}}_{{.Hour}}-{{.Minute}}-{{.Second}}_{{.Sequence}}
 Output: video/yamiodymel/2024-01-02_13-45-00_0.ts
```

_Note: Files are saved in `.ts` format, and this is not configurable._

&nbsp;

# 🤔 Frequently Asked Questions

**Q: The program closes immediately on Windows.**

> Open it via **Command Prompt**, the error message should appear. If needed, [create an issue](https://github.com/teacat/chaturbate-dvr/issues).

&nbsp;

**Q: Error `listen tcp :8080: bind: An attempt was... by its access permissions`**

> The port `8080` is in use. Try another port with `-p 8123`, then visit [http://localhost:8123](http://localhost:8123).
>
> If that fails, run **Command Prompt** as Administrator and execute:
>
> ```bash
> $ net stop winnat
> $ net start winnat
> ```

&nbsp;

**Q: Error `A connection attempt failed... host has failed to respond`**

> Likely a network issue (e.g., VPN, firewall, or blocked by Chaturbate). This cannot be fixed by the program.

&nbsp;

**Q: Error `Channel was blocked by Cloudflare`**

> You've been temporarily blocked. See the [Cookies & User-Agent](#-cookies--user-agent) section to bypass.

&nbsp;

**Q: Is Proxy or SOCKS5 supported?**

> Yes. You can launch the program using the `HTTPS_PROXY` environment variable:
>
> ```bash
> $ HTTPS_PROXY="socks5://127.0.0.1:9050" ./chaturbate-dvr -u CHANNEL_USERNAME
> ```
