#

## 📺 Framerate & Resolution / Fallback

A Fallback indicates what to do when there's no expected target resolution.

```
Availables: 1080p, 720p, 240p

Resolution: 480p (Fallback: UP)
    Result: 720p will be used
```

## 📄 Filename Pattern

The format is based on [Go Template Syntax](https://pkg.go.dev/text/template), available variables are:

`{{.Username}}`, `{{.Year}}`, `{{.Month}}`, `{{.Day}}`, `{{.Hour}}`, `{{.Minute}}`, `{{.Second}}`, `{{.Sequence}}`

Default:

```
Pattern: video/{{.Username}}_{{.Year}}-{{.Month}}-{{.Day}}_{{.Hour}}-{{.Minute}}-{{.Second}}_{{.Sequence}}
 Output: video/yamiodymel_2024-01-02_13-45-00_0.ts
```

👀 Hide sequence if it's zero, for better looking.

```
Pattern: {{.Username}}_{{.Year}}-{{.Month}}-{{.Day}}_{{.Hour}}-{{.Minute}}-{{.Second}}{{if .Sequence}}_{{.Sequence}}{{end}}
 Output: yamiodymel_2024-01-02_13-45-00.ts    # Sequence won't be shown if it's zero.
 Output: yamiodymel_2024-01-02_13-45-00_1.ts
```

📁 Folders per each channel, non-exists folder will be created automatically.

```
Pattern: video/{{.Username}}/{{.Year}}-{{.Month}}-{{.Day}}_{{.Hour}}-{{.Minute}}-{{.Second}}_{{.Sequence}}
 Output: video/yamiodymel/2024-01-02_13-45-00_0.ts
```

※ The file will be saved as `.ts` and it's not configurable.
