import { createContext, useContext } from 'react'
import type { Dispatch, SetStateAction } from 'react'

export type SoloDevTheme = 'light' | 'dark'

export interface SoloDevThemeContextValue {
  theme: SoloDevTheme
  setTheme: Dispatch<SetStateAction<SoloDevTheme>>
  toggleTheme: () => void
}

export const SOLODEV_THEME_STORAGE_KEY = 'solodev.theme'

export const getDefaultSoloDevTheme = (): SoloDevTheme => {
  if (typeof window === 'undefined') {
    return 'light'
  }

  try {
    const storedTheme = window.localStorage.getItem(SOLODEV_THEME_STORAGE_KEY)
    if (storedTheme) {
      const parsedTheme = JSON.parse(storedTheme) as SoloDevTheme
      if (parsedTheme === 'light' || parsedTheme === 'dark') {
        return parsedTheme
      }
    }
  } catch {
    // Ignore storage parsing issues and fall back to the preferred SoloDev default.
  }

  return 'light'
}

const noop = () => undefined

export const SoloDevThemeContext = createContext<SoloDevThemeContextValue>({
  theme: 'light',
  setTheme: noop as Dispatch<SetStateAction<SoloDevTheme>>,
  toggleTheme: noop
})

export const useSoloDevTheme = (): SoloDevThemeContextValue => useContext(SoloDevThemeContext)
