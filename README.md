# Chaturbate DVR

The program watches a specified Chaturbate channel and save the stream in real-time when the channel goes online, so you won't miss anything.

**Warning**: The streaming content on Chaturbate is copyrighted, you should not copy, share, distribute the content. (for more information, check [DMCA](https://www.dmca.com/))

## Usage

The program works for 64-bit macOS, Linux, Windows (too lazy to compile for 32-bit). Just get in the `/bin` folder and find your operating system then execute the program in terminal.

```bash
$ chaturbate-dvr -u my_lovely_channel_name

 .o88b. db   db  .d8b.  d888888b db    db d8888b. d8888b.  .d8b.  d888888b d88888b
d8P  Y8 88   88 d8' `8b `~~88~~' 88    88 88  `8D 88  `8D d8' `8b `~~88~~' 88'
8P      88ooo88 88ooo88    88    88    88 88oobY' 88oooY' 88ooo88    88    88ooooo
8b      88~~~88 88~~~88    88    88    88 88`8b   88~~~b. 88~~~88    88    88~~~~~
Y8b  d8 88   88 88   88    88    88b  d88 88 `88. 88   8D 88   88    88    88.
 `Y88P' YP   YP YP   YP    YP    ~Y8888P' 88   YD Y8888P' YP   YP    YP    Y88888P
d8888b. db    db d8888b.
88  `8D 88    88 88  `8D
88   88 Y8    8P 88oobY'
88   88 `8b  d8' 88`8b
88  .8D  `8bd8'  88 `88.
Y8888D'    YP    88   YD
---
2020/02/13 18:05:22 my_lovely_channel_name is online! fetching...
2020/02/13 18:05:24 the video will be saved as "2020-02-13 18:05:22.1344318 +0800 CST m=+0.885404701.ts".
2020/02/13 18:05:28 fetching media_w402018999_b5128000_t64RlBTOjI5Ljk3_9134.ts (size: 936428)
```

## Help

```bash
NAME:
   chaturbate-dvr - watching a specified chaturbate channel and auto saved to local file

USAGE:
   main [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --username value, -u value  channel username to watching
   --interval value, -i value  minutes to check if a channel goes online or not (default: 1)
   --help, -h                  show help (default: false)
```