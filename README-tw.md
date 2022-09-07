# Chaturbate 重新錄播 (Alpha)

這個程式能夠監聽指定的 Chaturbate 頻道，並且在該頻道開始直播時自動儲存影片至本機。這樣你就不會錯過任何精彩的事情

**警告**：在 Chaturbate 上的直播內容都有版權，你不應該複製、分享、散播這些內容。（想閱讀更多，請參閱 [DMCA](https://www.dmca.com/)）

**免責聲明**：因為這還在早期開發階段，錄播內容可能會有幀數遺失（3 小時直播遺失 20 秒內容），這仍然需要測試。

## 使用方式

這個程式能夠在 64 位元的 macOS、Linux、Windows（懶得編譯 32 位元的版本）上正常運作。你只需要進入 `/bin` 資料夾找到對應你的系統，然後在終端機執行該檔案即可。

```bash
$ chaturbate-dvr -u 好棒棒頻道名稱

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
2020/02/13 18:05:22 好棒棒頻道名稱 is online! fetching...
2020/02/13 18:05:24 the video will be saved as "2020-02-13_22-16-27.ts".
2020/02/13 18:05:28 fetching media_w402018999_b5128000_t64RlBTOjI5Ljk3_9134.ts (size: 936428)
2020/02/13 19:07:06 failed to fetch the video segments, will try again. (1/2)
2020/02/13 19:07:06 failed to fetch the video segments, will try again. (2/2)
2020/02/13 19:07:11 failed to fetch the video segments after retried, 好棒棒頻道名稱 might went offline.
2020/02/13 19:07:11 好棒棒頻道名稱 is not online, check again after 3 minute(s)...
```

## 說明

影片畫質永遠是以最高為優先，目前沒辦法更改（懶得寫成一個選項）。

```bash
NAME:
   chaturbate-dvr - watching a specified chaturbate channel and auto saves the stream as local file

USAGE:
   main [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --username value, -u value  channel username to watching
   --interval value, -i value  minutes to check if a channel goes online or not (default: 1)
   --strip value, -s value     MB sizes to split the video into chunks (default: 0)
   --help, -h                  show help (default: false)
```

## 中文對應

```
XXX is online! fetching...
XXX 正在線上！開始撈取實況內容…

the video will be saved as "XXX".
影片將會被保存為「XXX」。

fetching XXX.ts (size: XXX)
正在擷取 XXX.ts 片段（大小：XXX）

failed to fetch the video segments, will try again. (1/2)
無法取得影片段落，稍後會重新嘗試。(1/2)

failed to fetch the video segments after retried, XXX might went offline.
無法取得影片段落，XXX 可能已經結束直播了。

cannot find segment XXX, will try again. (1/5)
無法找到影片段落，燒後會重新嘗試。（1/5）

inserting XXX segment to the master file. (total: XXX)
正在插入片段 XXX 至主要影片檔案。（總共：XXX）

skipped XXX due to the empty body!
跳過 XXX 片段因為其為空白內容！

exceeded the specified stripping limit, creating new video file. (file: XXX)
達到影片分割上限，建立新的影片檔案（檔名：XXX）
```
