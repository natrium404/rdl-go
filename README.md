# RDL (Reel Downloader)

RDL (short for Reel Downloader) is a cross-platform lightweight application (currently supporting Windows and Linux) designed to download reels. Built with modern technologies, RDL offers a clean and user-friendly interface.

## Features

- **Cross-platform support**: Runs on both Windows and Linux.
- **Original reel download**: Fetch and download the original video.
- **Separate audio/video download**: Download the video and audio separately.
- **Chromium-based scraping**: Uses chromedp to scrape reel pages.
- Automatic browser detection.
- Detects existing Chromium-based browsers on the system.
- If none is found, download a `chromium-headless-shell` to a local cache.
  - Windows: `%USERPROFILE%\AppData\Local\chromedp\`
  - Linux: `~/.cache/chromedp/`
- **Dark-themed UI**: Built with Svelte, the UI offers a sleek design and enhanced user experience.

## Tech Stack

- **Frontend**: [Svelte](https://svelte.dev/)
- **Backend**: [Go](https://go.dev/)
- **Application Framework**: [Wails](https://wails.io/)
- **Web Scraping**: [chromedp](https://github.com/chromedp/chromedp)

## TODO

- [ ] Auto-update mechanism
- [ ] Support for downloading reels from private accounts

## Installation

### Option 1: Binaries

Goto the [Releases](https://github.com/natrium404/rdl-go/releases) page.

Download the binary for your OS:

- Windows: `rdl.exe`
- Linux: `rdl`
  (Optional) Add the binary to your system’s PATH for easier access.

#### How to Run:

- **Windows**:
  Double-click rdl.exe or run via Command Prompt:

```sh
rdl.exe
```

- **Linux**:

```sh
chmod +x rdl
./rdl
```

**NOTE:** RDL will automatically check for a Chromium-based browser. If none is found, it will download a chromium-headless-shell to a local cache folder.

### Option 2: Development Setup

#### Prerequisites:

- Go (v1.20 or later)
- Node.js (for Svelte)
- Wails CLI

#### Steps:

- Clone the repo:

```sh
git clone https://github.com/yourusername/rdl.git
cd rdl
```

- Install frontend dependencies:

```sh
cd frontend
npm install
cd ..
```

- Run in development mode:

```sh
wails dev
```

- Build production binary (optional):

```sh
wails build
```

## Contributing

For contribution [CONTRIBUTING.md](CONTRIBUTING.md)

## License

This project is open source and available under the [MIT License](LICENSE).

## FIN

![SLEEPY](./frontend/public/sleepy-sleepy-cat.gif)
