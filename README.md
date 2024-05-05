# Chaturbate DVR

The program can records **multiple** Chaturbate streams, supports macOS, Windows, Linux, can be run on Docker.

For Chaturbate-**only**, private/ticket stream is **unsupported**.

※ **[DMCA WARNING](https://www.dmca.com/)**: Contents on Chaturbate are copyrighted, you should not copy, share, distribute the content.

&nbsp;

## Usage

Download executable from **[Release](https://github.com/teacat/chaturbate-dvr/releases)** page (e.g., `windows_chatubrate-dvr.exe`).

&nbsp;

**🌐 Start the program with the Web UI**

Visit [`http://localhost:8080`](http://localhost:8080) to use the Web UI.

```yaml
# Windows (or double-click `chaturbate-dvr.exe` to open)
$ chaturbate-dvr.exe

# macOS or Linux
$ chaturbate-dvr
```

&nbsp;

**💻 or... Run as a command-line tool**

Run the program with a channel name (`-u CHANNEL_USERNAME`) records the channel immediately, and the Web UI will be disabled.

```yaml
# Windows
$ chaturbate-dvr.exe -u CHANNEL_USERNAME

# macOS or Linux
$ chaturbate-dvr -u CHANNEL_USERNAME
```

&nbsp;

## Preview

![image_1](https://github.com/teacat/chaturbate-dvr/assets/7308718/c6d17ffe-eba7-4296-9315-f501489d85f3)
![image_2](https://github.com/teacat/chaturbate-dvr/assets/7308718/d02923e0-574d-4a15-a373-8b0599101e3f)

**or... Command-line tool**

```
$ ./chaturbate-dvr -u emillybrowm start

 ██████╗██╗  ██╗ █████╗ ████████╗██╗   ██╗██████╗ ██████╗  █████╗ ████████╗███████╗
██╔════╝██║  ██║██╔══██╗╚══██╔══╝██║   ██║██╔══██╗██╔══██╗██╔══██╗╚══██╔══╝██╔════╝
██║     ███████║███████║   ██║   ██║   ██║██████╔╝██████╔╝███████║   ██║   █████╗
██║     ██╔══██║██╔══██║   ██║   ██║   ██║██╔══██╗██╔══██╗██╔══██║   ██║   ██╔══╝
╚██████╗██║  ██║██║  ██║   ██║   ╚██████╔╝██║  ██║██████╔╝██║  ██║   ██║   ███████╗
 ╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝   ╚═╝    ╚═════╝ ╚═╝  ╚═╝╚═════╝ ╚═╝  ╚═╝   ╚═╝   ╚══════╝
██████╗ ██╗   ██╗██████╗
██╔══██╗██║   ██║██╔══██╗
██║  ██║██║   ██║██████╔╝
██║  ██║╚██╗ ██╔╝██╔══██╗
██████╔╝ ╚████╔╝ ██║  ██║
╚═════╝   ╚═══╝  ╚═╝  ╚═╝
[2024-01-24 00:11:54] [INFO] [emillybrowm] channel created
[2024-01-24 00:11:55] [INFO] [emillybrowm] channel is online, start fetching...
[2024-01-24 00:11:55] [INFO] [emillybrowm] the stream will be saved as videos/emillybrowm_2024-01-24_00-11-55.ts
[2024-01-24 00:11:55] [INFO] [emillybrowm] resolution 1080p is used
[2024-01-24 00:11:55] [INFO] [emillybrowm] framerate 30fps is used
[2024-01-24 00:11:57] [INFO] [emillybrowm] segment #0 written
[2024-01-24 00:11:57] [INFO] [emillybrowm] segment #1 written
[2024-01-24 00:11:57] [INFO] [emillybrowm] segment #2 written
```

&nbsp;

## Help

```bash
$ chaturbate-dvr -h

NAME:
   chaturbate-dvr - Records your favorite Chaturbate stream 😎🫵

USAGE:
   chaturbate-dvr [global options] command [command options]

VERSION:
   1.0.0

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --username value, -u value               channel username to record
   --gui-username value, --gui-u value      username for auth web (optional)
   --gui-password value, --gui-p value      password for auth web (optional)
   --framerate value, -f value              preferred framerate (default: 30)
   --interval value, -i value               minutes to check if the channel is online (default: 1)
   --resolution value, -r value             preferred resolution (default: 1080)
   --resolution-fallback value, --rf value  fallback to 'up' (larger) or 'down' (smaller) resolution if preferred resolution is not available (default: "down")
   --filename-pattern value, --fp value     filename pattern for videos (default: "videos/{{.Username}}_{{.Year}}-{{.Month}}-{{.Day}}_{{.Hour}}-{{.Minute}}-{{.Second}}{{if .Sequence}}_{{.Sequence}}{{end}}")
   --split-duration value, --sd value       minutes to split each video into segments ('0' to disable) (default: 0)
   --split-filesize value, --sf value       size in MB to split each video into segments ('0' to disable) (default: 0)
   --log-level value                        log level, availables: 'DEBUG', 'INFO', 'WARN', 'ERROR' (default: "INFO")
   --port value                             port to expose the web interface and API (default: "8080")
   --help, -h                               show help
   --version, -v                            print the version
```

**Examples**:

```yaml
# Records in 720p/60fps
$ chaturbate-dvr -u yamiodymel -r 720 -f 60

# Split the video every 30 minutes
$ chaturbate-dvr -u yamiodymel -sd 30

# Split the video every 1024 MB
$ chaturbate-dvr -u yamiodymel -sf 1024

# Change output filename pattern
$ chaturbate-dvr -u yamiodymel -fp video/{{.Username}}/{{.Year}}-{{.Month}}-{{.Day}}_{{.Hour}}-{{.Minute}}-{{.Second}}_{{.Sequence}}
```

※ When runs in Web UI mode, the settings will be default settings for Web UI to create channels.

&nbsp;

## 📺 Framerate & Resolution / Fallback

Fallback indicates what to do when there's no expected target resolution, situation:

```
Availables: 1080p, 720p, 240p

Resolution: 480p (fallback setted to: up)
    Result: 720p will be used

Resolution: 480p (fallback setted to: down)
    Result: 240p will be used
```

&nbsp;

## 📄 Filename Pattern

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

※ The file will be saved as `.ts` format and it's not configurable.

&nbsp;

## 🤔 Frequently Asked Questions

**Q: The program closes itself when I just open it on Windows**

A: Try to open the program in **Command Prompt**, the error message should appear, and create a new [Issue](https://github.com/teacat/chaturbate-dvr/issues) for it.

&nbsp;

**Q: Channel is online but the program says it's Offline**

A: The program might be blocked by Chaturbate or Cloudflare. If the Channel is in a private/ticket show, the program doesn't support it yet.

&nbsp;

**Q: `listen tcp :8080: bind: An attempt was made to access a socket in a way forbidden by its access permissions.`**

A: The port `8080` is already in use, change the port with `-port` option (e.g. `-port 8123`) and visit `http://localhost:8123`.

If the error still occur, run **Command Prompt** as Administrator, and type `net stop winnat` then `net start winnat`, and re-run the Chaturbate DVR again.

&nbsp;

**Q: `A connection attempt failed because the connected party did not properly respond after a period of time, or established connection failed because connected host has failed to respond.`**

A: Your network is unstable or being blocked by Chaturbate, the program can't help with the network issue. Usually happened when you are using VPN or Proxy.

&nbsp;

## 💬 Verbose Log

Change `-log-level` to `DEBUG` to see more details in terminal, like Duration and Size.

```yaml
# Availables: DEBUG, INFO, WARN, ERROR
$ chaturbate-dvr -u hepbugbear -log-level DEBUG
[2024-01-24 01:18:11] [INFO] [hepbugbear] segment #0 written
[2024-01-24 01:18:11] [DEBUG] [hepbugbear] duration: 00:00:06, size: 0.00 MiB
[2024-01-24 01:18:11] [INFO] [hepbugbear] segment #1 written
[2024-01-24 01:18:11] [DEBUG] [hepbugbear] duration: 00:00:06, size: 1.36 MiB
[2024-01-24 01:18:11] [INFO] [hepbugbear] segment #2 written
[2024-01-24 01:18:11] [DEBUG] [hepbugbear] duration: 00:00:06, size: 2.72 MiB
[2024-01-24 01:18:12] [DEBUG] [hepbugbear] segment #3 fetched
[2024-01-24 01:18:13] [INFO] [hepbugbear] segment #3 written
[2024-01-24 01:18:13] [DEBUG] [hepbugbear] duration: 00:00:10, size: 4.08 MiB
```
