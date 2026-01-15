import { App } from '../bindings/fakering'

export async function loadConfig() {
  try {
    const config = await App.LoadConfig()
    return config
  } catch (error) {
    console.error('Failed to load config:', error)
    return null
  }
}

export async function saveConfig(config) {
  try {
    await App.SaveConfig(config)
    return true
  } catch (error) {
    console.error('Failed to save config:', error)
    return false
  }
}

export async function getSetting(key) {
  try {
    return await App.GetSetting(key)
  } catch (error) {
    console.error('Failed to get setting:', error)
    return null
  }
}

export async function setSetting(key, value) {
  try {
    await App.SetSetting(key, value)
    return true
  } catch (error) {
    console.error('Failed to set setting:', error)
    return false
  }
}
