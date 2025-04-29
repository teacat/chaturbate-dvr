# Chaturbate DVR

The program can records **multiple** Chaturbate streams, supports macOS, Windows, Linux, can be run on Docker.

For Chaturbate-**only**.

※ **[DMCA WARNING](https://www.dmca.com/)**: Contents on Chaturbate are copyrighted, you should not copy, share, distribute the content.

&nbsp;

## Getting Started

Download executable from **[Release](https://github.com/teacat/chaturbate-dvr/releases)** page (e.g., `x64_windows_chatubrate-dvr.exe`)

&nbsp;

**1. 🌐 Start the program with the Web UI**

```yaml
# Windows (or double-click `x64_windows_chatubrate-dvr.exe` to open)
$ x64_windows_chatubrate-dvr.exe

# macOS or Linux
$ ./x64_linux_chatubrate-dvr
```

Visit [`http://localhost:8080`](http://localhost:8080) to use the Web UI.

&nbsp;

**2. 💻 Run as a command-line tool**

```yaml
# Windows
$ x64_windows_chatubrate-dvr.exe -u CHANNEL_USERNAME

# macOS or Linux
$ ./x64_linux_chatubrate-dvr -u CHANNEL_USERNAME
```

This records the `CHANNEL_USERNAME` channel immediately, and the Web UI won't be available.

&nbsp;

**3. 🐳 Run on Docker**

```yaml
# Windows
$ x64_windows_chatubrate-dvr.exe -u CHANNEL_USERNAME

# macOS or Linux
$ ./x64_linux_chatubrate-dvr -u CHANNEL_USERNAME
```

&nbsp;

## Command-line

```bash
$ chaturbate-dvr -h

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
   --cf-cookie value                        Cloudflare cookie to bypass anti-bot page
   --user-agent value                       Custom user agent for when using cf-cookie
   --help, -h                               show help
   --version, -v                            print the version
```

**Examples**:

```yaml
# Records in 720p/60fps
$ ./x64_linux_chatubrate-dvr -u yamiodymel -r 720 -f 60

# Split the video every 30 minutes
$ ./x64_linux_chatubrate-dvr -u yamiodymel -sd 30

# Split the video every 1024 MB
$ ./x64_linux_chatubrate-dvr -u yamiodymel -sf 1024

# Change output filename pattern
$ ./x64_linux_chatubrate-dvr -u yamiodymel -fp video/{{.Username}}/{{.Year}}-{{.Month}}-{{.Day}}_{{.Hour}}-{{.Minute}}-{{.Second}}_{{.Sequence}}
```

※ In Web UI mode, the settings are used as the default values for creating channels.

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

**Q: The program closes itself when I just open it on Windows.**

> Try to open the program in **Command Prompt**, the error message should appear. Create a new [Issue](https://github.com/teacat/chaturbate-dvr/issues) for it.

&nbsp;

**Q: Error message `listen tcp :8080: bind: An attempt was made to access a socket in a way forbidden by its access permissions`**

> The port `8080` is already in use. Change the port using the `-p` option (e.g., `-p 8123`), then visit `http://localhost:8123`.
>
> If the error still occurs, run **Command Prompt** as Administrator, and enter the following commands:
>
> ```
> net stop winnat
> net start winnat
> ```
>
> After that, re-open Chaturbate DVR.

&nbsp;

**Q: Error message `A connection attempt failed because the connected party did not properly respond after a period of time, or established connection failed because connected host has failed to respond`**

> Your network is unstable or may be blocked by Chaturbate. This program can't fix network-related issues, which often occur when using a VPN or proxy.

&nbsp;

**Q: Error message `channel was blocked by Cloudflare`**

> Chaturbate has temporarily blocked your access due to scraping activity. Please refer to the [Cookies & User-Agent](#!) section above for more details.

&nbsp;
