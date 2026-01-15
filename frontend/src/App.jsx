import { useState, useEffect, useRef } from 'react'
import { RingLightService } from "../bindings/fakering";
import { getSetting, setSetting } from './config-helper';

function App() {
  const [monitors, setMonitors] = useState([]);
  const [monitorSettings, setMonitorSettings] = useState({});
  const saveTimeoutRef = useRef(null);

  useEffect(() => {
    const initialize = async () => {
      try {
        const list = await RingLightService.GetMonitors();
        setMonitors(list);
        
        // Load initial settings
        const savedSettings = await getSetting('monitorSettings');
        
        // Ensure backend global state is true (handled in backend now, but good to be sure)
        RingLightService.SetEnabled(true);

        const settings = {};
        list.forEach(m => {
          if (savedSettings && savedSettings[m]) {
            settings[m] = savedSettings[m];
            // Apply saved settings to backend immediately
            const s = savedSettings[m];
            const rgb = hexToRgb(s.color);
            if (rgb) RingLightService.SetColor(m, rgb.r, rgb.g, rgb.b);
            RingLightService.SetBrightness(m, s.brightness);
            RingLightService.SetWidth(m, s.width);
            RingLightService.ToggleMonitor(m, s.enabled);
          } else {
            settings[m] = {
              enabled: false,
              color: '#FFEE08',
              brightness: 200,
              width: 30
            };
            // Ensure backend is also off by default
            RingLightService.ToggleMonitor(m, false);
          }
        });
        setMonitorSettings(settings);
      } catch (e) {
        console.error("Failed to initialize:", e);
      }
    };
    
    // Small delay to ensure backend is ready
    const timer = setTimeout(initialize, 1000);
    return () => {
      clearTimeout(timer);
      if (saveTimeoutRef.current) clearTimeout(saveTimeoutRef.current);
    };
  }, []);

  const hexToRgb = (hex) => {
    const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);
    return result ? {
      r: parseInt(result[1], 16),
      g: parseInt(result[2], 16),
      b: parseInt(result[3], 16)
    } : null;
  };

  const hslToHex = (h, s, l) => {
    l /= 100;
    const a = s * Math.min(l, 1 - l) / 100;
    const f = n => {
      const k = (n + h / 30) % 12;
      const color = l - a * Math.max(Math.min(k - 3, 9 - k, 1), -1);
      return Math.round(255 * color).toString(16).padStart(2, '0');
    };
    return `#${f(0)}${f(8)}${f(4)}`;
  };

  const handleRingInteraction = (e, monitorName) => {
    const rect = e.currentTarget.getBoundingClientRect();
    const centerX = rect.left + rect.width / 2;
    const centerY = rect.top + rect.height / 2;
    
    // Get angle in radians, then convert to degrees (0-360)
    // Math.atan2 returns -PI to PI
    let angle = Math.atan2(e.clientY - centerY, e.clientX - centerX) * (180 / Math.PI);
    
    // Adjust so 0 is red (at the top or right depending on gradient start)
    // Conic gradient by default starts at top (0deg). 
    // Atan2 0 is right. -90 is top. 
    angle = (angle + 90 + 360) % 360;
    
    const hexColor = hslToHex(angle, 100, 50);
    updateSetting(monitorName, 'color', hexColor);
  };

  const debouncedSave = (settings) => {
    if (saveTimeoutRef.current) clearTimeout(saveTimeoutRef.current);
    saveTimeoutRef.current = setTimeout(() => {
      setSetting('monitorSettings', settings);
    }, 500);
  };

  const updateSetting = (name, key, value) => {
    const newSettings = {
      ...monitorSettings,
      [name]: { ...monitorSettings[name], [key]: value }
    };
    setMonitorSettings(newSettings);
    debouncedSave(newSettings);

    switch(key) {
      case 'color':
        const rgb = hexToRgb(value);
        if (rgb) RingLightService.SetColor(name, rgb.r, rgb.g, rgb.b);
        break;
      case 'brightness':
        RingLightService.SetBrightness(name, value);
        break;
      case 'width':
        RingLightService.SetWidth(name, value);
        break;
      case 'enabled':
        RingLightService.ToggleMonitor(name, value);
        break;
    }
  };



  return (
    <div className="flex flex-col h-screen w-screen overflow-hidden bg-[#0a0c14] text-slate-200 select-none items-center">
      {/* Main Content */}
      <main className="w-full max-w-7xl flex flex-col p-10 overflow-y-auto">
        <header className="flex items-center justify-between mb-10 w-full">
          <h1 className="text-3xl font-bold tracking-tight text-white">FakeRing</h1>
        </header>

        <section className="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-8">
          {monitors.length > 0 ? monitors.map((m, i) => {
            const s = monitorSettings[m] || { enabled: false, color: '#FFEE08', brightness: 200, width: 30 };
            return (
              <div key={m} className={`glass rounded-[32px] transition-all duration-500 border-2 flex flex-col p-6 ${s.enabled ? 'gap-8 border-yellow-500/20 bg-slate-800/5' : 'gap-0 border-transparent bg-slate-900/10 opacity-60'}`}>
                <div className="flex items-center justify-between w-full">
                  <div className="flex flex-col">
                    <span className="text-[10px] uppercase font-bold text-slate-500 tracking-widest mb-1">Display Interface</span>
                    <span className="text-sm font-semibold text-slate-300">{m}</span>
                  </div>
                  <label className="switch">
                    <input 
                      type="checkbox" 
                      checked={s.enabled} 
                      onChange={() => updateSetting(m, 'enabled', !s.enabled)} 
                    />
                    <span className="slider"></span>
                  </label>
                </div>

                {s.enabled && (
                  <div className="flex flex-col gap-8 w-full animate-in fade-in slide-in-from-top-4 duration-500">
                    <div className="flex justify-center py-4">
                      <div 
                        className="relative group color-ring-container w-44 h-44" 
                        onMouseDown={(e) => handleRingInteraction(e, m)}
                        onMouseMove={(e) => e.buttons === 1 && handleRingInteraction(e, m)}
                      >
                        <div className="color-inner-circle border-[8px]">
                           <div 
                             className="w-24 h-24 rounded-full shadow-2xl transition-all duration-300" 
                             style={{ 
                               backgroundColor: s.color,
                               boxShadow: `0 0 30px ${s.color}66`
                             }}
                           />
                        </div>
                      </div>
                    </div>

                    <div className="flex flex-col gap-6 mt-auto">
                      <div className="flex flex-col gap-3">
                        <div className="flex justify-between items-end">
                          <span className="text-xs text-slate-500 font-bold uppercase tracking-widest">Intensity</span>
                          <span className="text-slate-400 text-sm font-mono">{Math.round((s.brightness/255)*100)}</span>
                        </div>
                        <input 
                          type="range" 
                          min="0" 
                          max="255" 
                          value={s.brightness} 
                          onChange={(e) => updateSetting(m, 'brightness', parseInt(e.target.value))} 
                          style={{
                            background: `linear-gradient(to right, #FFEE08 0%, #FFEE08 ${(s.brightness/255)*100}%, rgba(255, 255, 255, 0.1) ${(s.brightness/255)*100}%, rgba(255, 255, 255, 0.1) 100%)`
                          }}
                        />
                      </div>

                      <div className="flex flex-col gap-3">
                        <div className="flex justify-between items-end">
                          <span className="text-xs text-slate-500 font-bold uppercase tracking-widest">Core Weight</span>
                          <span className="text-slate-400 text-sm font-mono">{s.width}</span>
                        </div>
                        <input 
                          type="range" 
                          min="1" 
                          max="150" 
                          value={s.width} 
                          onChange={(e) => updateSetting(m, 'width', parseInt(e.target.value))} 
                          style={{
                            background: `linear-gradient(to right, #FFEE08 0%, #FFEE08 ${((s.width-1)/149)*100}%, rgba(255, 255, 255, 0.1) ${((s.width-1)/149)*100}%, rgba(255, 255, 255, 0.1) 100%)`
                          }}
                        />
                      </div>
                    </div>
                  </div>
                )}
              </div>
            )
          }) : (
            <div className="col-span-full py-20 text-center flex flex-col items-center gap-4">
              <div className="w-10 h-10 border-2 border-yellow-500/20 border-t-yellow-500 rounded-full animate-spin"></div>
              <p className="text-slate-500 font-medium">Initializing display protocols...</p>
            </div>
          )}
        </section>
      </main>
    </div>
  )
}




export default App
