# FakeRing Elite

![FakeRing Logo](fakering_logo_1768489730286.png)

**FakeRing Elite** is a high-performance cross-platform utility built with Wails v3 that turns your monitors into professional-grade ring lights. Perfect for video calls, streaming, or content creation when you don't have physical lighting gear.

## Screenshots

![FakeRing Elite Mockup](fakering_screenshot_mockup_1768489748878.png)

## Why this exists

Professional lighting is expensive and bulky. **FakeRing Elite** leverages the most powerful light source already on your desk—your monitors. By rendering precisely controlled, high-intensity borders on your displays, it provides soft, adjustable lighting that eliminates shadows and makes you look great on camera without any extra hardware.

## Installation

### Prerequisites
- [Go 1.25+](https://golang.org/dl/)
- [Node.js 20+](https://nodejs.org/)
- [Wails v3](https://v3.wails.io/docs/getting-started/installation)

### Build from source
```bash
# Clone the repository
git clone https://github.com/yourusername/Fakering.git
cd Fakering

# Run in development mode
wails3 dev

# Build production binary
wails3 build
```

## Example Usage

### App Configuration (Persistence)
The application uses a custom `AppConfig` helper to remember your settings across sessions.

```typescript
import { getSetting, setSetting } from './config-helper';

// Save monitor configuration
await setSetting('monitorSettings', {
  'Display 1': { enabled: true, color: '#FFEE08', brightness: 255, width: 50 },
  'Display 2': { enabled: false, color: '#FFFFFF', brightness: 128, width: 30 }
});

// Load global state
const isEnabled = await getSetting('enabled');
```

## Links

- [Documentation](#)
- [Website](#)
- [Blog](#)

## Tech Used

- **Backend**: Go with [Wails v3](https://v3.wails.io/)
- **Frontend**: React + Tailwind CSS
- **Design**: Glassmorphism / Cyberpunk Aesthetics
- **Core Logic**: Windows User32/GDI32 Native APIs for high-performance transparent overlays

## Inspiration

Inspired by the need for better lighting in remote work environments and the power of modern web technologies combined with native system performance.

## Contributing

Issues and Pull Requests are welcome! Please feel free to check the issues page or submit a new feature request.

## License

MIT © [Your Name/Company]
