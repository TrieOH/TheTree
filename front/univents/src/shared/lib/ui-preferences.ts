export type ThemeMode = 'auto' | 'light' | 'dark'

export const THEME_STORAGE_KEY = 'theme'
export const INPLACE_EDIT_STORAGE_KEY = 'inplace-edit'
export const UI_PREFERENCES_CHANGE_EVENT = 'ui-preferences-change'

export function readThemePreference(): ThemeMode {
  if (typeof window === 'undefined') return 'auto'

  const stored = window.localStorage.getItem(THEME_STORAGE_KEY)
  return stored === 'light' || stored === 'dark' || stored === 'auto' ? stored : 'auto'
}

export function writeThemePreference(theme: ThemeMode) {
  if (typeof window === 'undefined') return
  window.localStorage.setItem(THEME_STORAGE_KEY, theme)
  window.dispatchEvent(new Event(UI_PREFERENCES_CHANGE_EVENT))
}

export function readInplaceEditPreference(): boolean {
  if (typeof window === 'undefined') return true

  const stored = window.localStorage.getItem(INPLACE_EDIT_STORAGE_KEY)
  return stored === null ? true : stored === 'true'
}

export function writeInplaceEditPreference(enabled: boolean) {
  if (typeof window === 'undefined') return
  window.localStorage.setItem(INPLACE_EDIT_STORAGE_KEY, String(enabled))
  window.dispatchEvent(new Event(UI_PREFERENCES_CHANGE_EVENT))
}

export function applyThemePreference(theme: ThemeMode) {
  if (typeof document === 'undefined') return

  const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
  const resolved = theme === 'auto' ? (prefersDark ? 'dark' : 'light') : theme
  const root = document.documentElement

  root.classList.remove('light', 'dark')
  root.classList.add(resolved)

  if (theme === 'auto') root.removeAttribute('data-theme')
  else root.setAttribute('data-theme', theme)

  root.style.colorScheme = resolved
}
