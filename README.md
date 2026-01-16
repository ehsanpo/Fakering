# FakeRing

![FakeRing Logo](build/appicon.png)

**FakeRing** is a cross-platform utility built with Wails v3 that turns your monitors into ring lights. Perfect for video calls, streaming, or content creation when you don't have physical lighting gear.

## Screenshots

![FakeRing Mockup](demo/cover.png)

## Why this exists

Professional lighting is expensive and bulky. **FakeRing** leverages the most powerful light source already on your desk—your monitors. By rendering precisely controlled, high-intensity borders on your displays, it provides soft, adjustable lighting that eliminates shadows and makes you look great on camera without any extra hardware.

## Installation

### Prerequisites

- [Go 1.25+](https://golang.org/dl/)
- [Node.js 20+](https://nodejs.org/)
- [Wails v3](https://v3.wails.io/docs/getting-started/installation)

### Build from source

```bash
# Clone the repository
git clone https://github.com/ehsanpo/Fakering.git
cd Fakering

# Run in development mode
wails3 dev

# Build production binary
wails3 build
```

## Links

- [Documentation](#)
- [Website](#)
- [Blog](#)

## Tech Used

- **Backend**: Go with [Wails v3](https://v3.wails.io/)
- **Frontend**: React + Tailwind CSS
- **Core Logic**: Windows User32/GDI32 Native APIs for high-performance transparent overlays

## Inspiration

Inspired by the need for better lighting in remote work environments and the power of modern web technologies combined with native system performance.

## Contributing

Issues and Pull Requests are welcome! Please feel free to check the issues page or submit a new feature request.

## License

MIT © [Ehsan Pourhadi](https://github.com/ehsanpo)
