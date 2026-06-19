# BeatportDL

Beatport & Beatsource downloader (FLAC, AAC)

*Requires an active [Beatport](https://stream.beatport.com/) or [Beatsource](https://stream.beatsource.com/) streaming plan.*

![Screenshot](/screenshots/main.png?raw=true "Screenshot")

What's new in this fork
---
This fork keeps BeatportDL's behaviour and configuration **identical**, while making it far easier to build and trimming the codebase:

- **No CGO, no C toolchain.** Audio tagging now uses [go-taglib](https://github.com/sentriz/go-taglib) — the real [TagLib](https://taglib.org/) library compiled to WebAssembly and run via [wazero](https://github.com/tetratelabs/wazero) — instead of CGO bindings. You no longer need TagLib, zlib or the Zig toolchain installed: `go build` and cross-compilation are pure Go. Tags are still produced by TagLib, so files are tagged exactly as before.
- **One-command cross-compilation.** The Makefile is now plain `GOOS`/`GOARCH` builds — a single `make` produces binaries for macOS, Linux and Windows with no per-platform C library paths.
- **Leaner download code.** The chart, playlist and artist download paths shared a lot of duplicated logic; they now run through a single shared helper (~165 fewer lines), which keeps behaviour consistent across link types and reduces the surface for bugs.

> **Requirements changed:** building now needs **Go 1.25+**. The advanced M4A `_raw` tag option now goes through TagLib's standard property mapping; the default config doesn't use it, so most setups are unaffected.

Building & testing
---
> **Prerequisites:** [Go 1.25+](https://go.dev/dl/), and an active [Beatport](https://stream.beatport.com/)/[Beatsource](https://stream.beatsource.com/) plan. [ffmpeg](https://www.ffmpeg.org/download.html) is only needed for the `medium-hls` quality option.

**1. Clone this branch and build**
```shell
git clone -b simplify/dedup-and-drop-cgo https://github.com/brian0h3c/beatportdl.git
cd beatportdl
go build ./cmd/beatportdl          # produces ./beatportdl
```
No TagLib/zlib/Zig install required. To prove there's no hidden C dependency, build with CGO disabled:
```shell
CGO_ENABLED=0 go build ./cmd/beatportdl
```
To build every release binary (output in `./bin`), run `make` — see [Building](#building).

**2. First run — create the config**

Run it once and answer the prompts (username, password, downloads directory, quality). This writes `beatportdl-config.yml`, and on a successful login `beatportdl-credentials.json`:
```shell
./beatportdl
```

**3. Download something**

Pass one or more URLs as arguments:
```shell
./beatportdl https://www.beatport.com/track/strobe/1696999
```
or run `./beatportdl` with no arguments and type a search query when prompted (add `@beatsource` to search Beatsource instead):
```shell
./beatportdl
Enter url or search query: deadmau5 strobe
```
Tracks, releases, playlists, charts, labels and artists are all supported — just paste the URL. Set `sort_by_context: true` in the config to group downloads into per-release/playlist/chart folders.

**4. Verify the tags and cover art**

Open a downloaded file in your player or DJ software and confirm the metadata and artwork look correct. You can also inspect it from the terminal:
```shell
ffprobe -hide_banner "your-downloaded-track.flac"   # shows tags + the attached cover stream
```
FLAC, AAC/M4A and `medium-hls` downloads should all be tagged correctly — artist, title, album, BPM, key, ISRC, label and embedded artwork included.

Setup
---
1. [Download](https://github.com/unspok3n/beatportdl/releases/) or [build](#building) BeatportDL.

     *Compiled binaries for Windows, macOS (amd64, arm64) and Linux (amd64, arm64) are available on the [Releases](https://github.com/unspok3n/beatportdl/releases) page.* \
     *Don't forget to set the execute permission on unix systems, e.g., chmod +x beatportdl-darwin-arm64*

2. Run beatportdl (e.g. `./beatportdl-darwin-arm64`), then specify the:
   - Beatport username
   - Beatport password
   - Downloads directory
   - Audio quality

3. OPTIONAL: Customize a config file. Create a new config file by running:
```shell
./beatportdl
```
This will create a new `beatportdl-config.yml` file. You can put the following options and values into the config file:

---
| Option                        | Default Value                             | Type       | Description                                                                                                                                                                               |
|-------------------------------|-------------------------------------------|------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `username`                    |                                           | String     | Beatport username                                                                                                                                                                         |
| `password`                    |                                           | String     | Beatport password                                                                                                                                                                         |
| `quality`                     | lossless                                  | String     | Download quality *(medium-hls, medium, high, lossless)*                                                                                                                                   |
| `show_progress`               | true                                      | Boolean    | Enable progress bars                                                                                                                                                                      |
| `write_error_log`             | false                                     | Boolean    | Write errors to `error.log`                                                                                                                                                               |
| `max_download_workers`        | 15                                        | Integer    | Concurrent download jobs limit                                                                                                                                                            |
| `max_global_workers`          | 15                                        | Integer    | Concurrent global jobs limit                                                                                                                                                              |
| `downloads_directory`         |                                           | String     | Location for the downloads directory                                                                                                                                                      |
| `sort_by_context`             | false                                     | Boolean    | Create a directory for each release, playlist, chart, label, or artist                                                                                                                    |
| `sort_by_label`               | false                                     | Boolean    | Use label names as parent directories for releases (requires `sort_by_context`)                                                                                                           |
| `force_release_directories`   | false                                     | Boolean    | Create release directories inside chart and playlist folders (requires `sort_by_context`)                                                                                                 |
| `track_exists`                | update                                    | String     | Behavior when track file already exists                                                                                                                                                   |
| `track_number_padding`        | 2                                         | Integer    | Track number padding for filenames and tag mappings (when using `track_number_with_padding` or `release_track_count_with_padding`)<br/> Set to 0 for dynamic padding based on track count |
| `cover_size`                  | 1400x1400                                 | String     | Cover art size for `keep_cover` and track metadata (if `fix_tags` is enabled)  *[max: 1400x1400]*                                                                                         |
| `keep_cover`                  | false                                     | Boolean    | Download cover art file (cover.jpg) to the context directory (requires `sort_by_context`)                                                                                                 |
| `fix_tags`                    | true                                      | Boolean    | Enable tag writing capabilities                                                                                                                                                           |
| `tag_mappings`                | *Listed below*                            | String Map | Custom tag mappings                                                                                                                                                                       |
| `track_file_template`         | {number}. {artists} - {name} ({mix_name}) | String     | Track filename template                                                                                                                                                                   |
| `release_directory_template`  | [{catalog_number}] {artists} - {name}     | String     | Release directory template                                                                                                                                                                |
| `playlist_directory_template` | {name} [{created_date}]                   | String     | Playlist directory template                                                                                                                                                               |
| `chart_directory_template`    | {name} [{published_date}]                 | String     | Chart directory template                                                                                                                                                                  |
| `label_directory_template`    | {name} [{updated_date}]                   | String     | Label directory template                                                                                                                                                                  |
| `artist_directory_template`   | {name}                                    | String     | Artist directory template                                                                                                                                                                 |
| `whitespace_character`        |                                           | String     | Whitespace character for track filenames and release directories                                                                                                                          |
| `artists_limit`               | 3                                         | Integer    | Maximum number of artists allowed before replacing with `artists_short_form` (affects directories, filenames, and search results)                                                         |
| `artists_short_form`          | VA                                        | String     | Custom string to represent "Various Artists"                                                                                                                                              |
| `key_system`                  | standard-short                            | String     | Music key system used in filenames and tags                                                                                                                                               |
| `proxy`                       |                                           | String     | Proxy URL                                                                                                                                                                                 |

If the Beatport credentials are correct, you should also see the file `beatportdl-credentials.json` appear in the BeatportDL directory.
*If you accidentally entered an incorrect password and got an error, you can always manually edit the config file*

Download quality options, per Beatport/Beatsource subscription type:

| Option       | Description                                                                                                  | Requires at least              | Notes                                                                   |
|--------------|--------------------------------------------------------------------------------------------------------------|--------------------------------|-------------------------------------------------------------------------|
| `medium-hls` | 128 kbps AAC through `/stream` endpoint (IMPORTANT: requires [ffmpeg](https://www.ffmpeg.org/download.html)) | Essential / Beatsource         | Same as `medium` on Advanced but uses a slightly slower download method |
| `medium`     | 128 kbps AAC                                                                                                 | Advanced / Beatsource Pro+     |                                                                         |
| `high`       | 256 kbps AAC                                                                                                 | Professional / Beatsource Pro+ |                                                                         |
| `lossless`   | 44.1 kHz FLAC                                                                                                | Professional / Beatsource Pro+ |                                                                         |
**Pioneer DJ / AlphaTheta player compatibility:**
* **`lossless` (FLAC)** — highest quality. Plays in rekordbox (Mac/PC/iOS) and on modern standalone players: **CDJ-3000, CDJ-2000NXS2, XDJ-XZ, XDJ-RX3/RX2, XDJ-RR, XDJ-AZ, OPUS-QUAD**. Older players (CDJ-2000NXS / CDJ-900NXS and earlier) **don't read FLAC** — use AAC for those.
* **`high` / `medium` / `medium-hls` (AAC / `.m4a`)** — play in rekordbox and **virtually all** Pioneer DJ / AlphaTheta gear, including older CDJs that don't support FLAC. Use `high` for the best AAC quality.

> If you're playing off USB on a standalone CDJ/XDJ, you still analyze and export the tracks with **rekordbox** first. rekordbox itself supports both FLAC and AAC.
Available `track_exists` options:
* `error` Log error and skip
* `skip` Skip silently
* `overwrite` Re-download
* `update` Update tags

Available template keywords for filenames and directories (`*_template`):
* Track: `id`,`name`,`mix_name`,`slug`,`artists`,`remixers`,`number`,`length`,`key`,`bpm`,`genre`,`subgenre`,`genre_with_subgenre`,`subgenre_or_genre`,`isrc`,`label`
* Release: `id`,`name`,`slug`,`artists`,`remixers`,`date`,`year`,`track_count`,`bpm_range`,`catalog_number`,`upc`,`label`
* Playlist: `id`,`name`,`first_genre`,`track_count`,`bpm_range`,`length`,`created_date`,`updated_date`
* Chart: `id`,`name`,`slug`,`first_genre`,`track_count`,`creator`,`created_date`,`published_date`,`updated_date`
* Artist: `id`, `name`, `slug`
* Label: `id`, `name`, `slug`, `created_date`, `updated_date`

Default `tag_mappings` config:
```yaml
tag_mappings:
   flac:
      track_name: "TITLE"
      track_artists: "ARTIST"
      track_number: "TRACKNUMBER"
      track_subgenre_or_genre: "GENRE"
      track_key: "KEY"
      track_bpm: "BPM"
      track_isrc: "ISRC"
   
      release_name: "ALBUM"
      release_artists: "ALBUMARTIST"
      release_date: "DATE"
      release_track_count: "TOTALTRACKS"
      release_catalog_number: "CATALOGNUMBER"
      release_label: "LABEL"
   m4a:
      track_name: "TITLE"
      track_artists: "ARTIST"
      track_number: "TRACKNUMBER"
      track_genre: "GENRE"
      track_key: "KEY"
      track_bpm: "BPM"
      track_isrc: "ISRC"
   
      release_name: "ALBUM"
      release_artists: "ALBUMARTIST"
      release_date: "DATE"
      release_track_count: "TOTALTRACKS"
      release_catalog_number: "CATALOGNUMBER"
      release_label: "LABEL"
```

As you can see, each key here represents a predefined value from either a release or a track that you can use to customize what is written to which tags. When you add an entry in the mappings for any format (for e.g., `flac`), only the tags that you specify will be written.

All tags by default are converted to uppercase, but since some M4A players might not recognize it, you can write the tag in lowercase and add the `_raw` suffix to bypass the conversion. *(This applies to M4A tags only)*

> **Fork note:** tags are now written by the bundled WebAssembly TagLib, which applies its own property mapping. The `_raw` suffix is still accepted and stripped, but the resulting atom casing follows TagLib's behaviour rather than being written verbatim. If you rely on a specific lower-cased atom (e.g. Traktor's `initialkey`), verify it in your player before bulk-downloading.

For e.g., Traktor doesn't recognize the track key tag in uppercase, so you have to add:
```yaml
tag_mappings:
   m4a:
      track_key: "initialkey_raw"
```

Available `tag_mappings` keys: `track_id`,`track_url`,`track_name`,`track_artists`,`track_artists_limited`,`track_remixers`,`track_remixers_limited`,`track_number`,`track_number_with_padding`,`track_number_with_total`,`track_genre`,`track_subgenre`,`track_genre_with_subgenre`,`track_subgenre_or_genre`,`track_key`,`track_bpm`,`track_isrc`,`release_id`,`release_url`,`release_name`,`release_artists`,`release_artists_limited`,`release_remixers`,`release_remixers_limited`,`release_date`,`release_year`,`release_track_count`,`release_track_count_with_padding`,`release_catalog_number`,`release_upc`,`release_label`,`release_label_url`

Available `key_system` options:

| System           | Example           |
|------------------|-------------------|
| `standard`       | Eb Minor, F Major |
| `standard-short` | Ebm, F            |
| `openkey`        | 7m, 12d           |
| `camelot`        | 2A, 7B            |

Proxy URL format example: `http://username:password@127.0.0.1:8080`

Usage
---

Run BeatportDL and enter Beatport or Beatsource URL or search query:
```shell
./beatportdl
Enter url or search query:
```
By default, search returns the results from beatport, if you want to search on beatsource instead, include `@beatsource` tag in the query

...or specify the URL using positional arguments:
```shell
./beatportdl https://www.beatport.com/track/strobe/1696999 https://www.beatport.com/track/move-for-me/591753
```
...or provide a text file with urls (separated by a newline)
```shell
./beatportdl file.txt file2.txt
```

URL types that are currently supported: **Tracks, Releases, Playlists, Charts, Labels, Artists**

Building
---
Required dependencies:
* [Go](https://go.dev/dl/) >= 1.25

BeatportDL handles audio metadata with [go-taglib](https://github.com/sentriz/go-taglib), which embeds [TagLib](https://taglib.org/) compiled to WebAssembly and runs it via [wazero](https://github.com/tetratelabs/wazero). There is **no CGO and no C toolchain required**, so building and cross-compilation are plain Go.

Build for your current platform:
```shell
go build ./cmd/beatportdl
```

The Makefile builds release binaries for every supported platform (output in `./bin`):
```shell
make                # all targets
make darwin-arm64
make darwin-amd64
make linux-amd64
make linux-arm64
make windows-amd64
```

> **Note:** the `medium-hls` quality option requires [ffmpeg](https://www.ffmpeg.org/download.html) at runtime to remux the downloaded stream segments.
