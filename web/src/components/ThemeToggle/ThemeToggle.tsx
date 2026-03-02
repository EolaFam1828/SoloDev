import React from 'react'
import cx from 'classnames'
import { HalfMoon, SunLight } from 'iconoir-react'
import { useSoloDevTheme } from 'contexts/SoloDevThemeContext'
import css from './ThemeToggle.module.scss'

export const ThemeToggle: React.FC = () => {
  const { theme, setTheme } = useSoloDevTheme()

  return (
    <div className={css.main}>
      <span className={css.label}>Theme</span>
      <div className={css.toggle} role="tablist" aria-label="SoloDev theme">
        <button
          type="button"
          className={cx(css.option, { [css.active]: theme === 'light' })}
          onClick={() => setTheme('light')}
          role="tab"
          aria-selected={theme === 'light'}>
          <SunLight />
          <span>Light</span>
        </button>
        <button
          type="button"
          className={cx(css.option, { [css.active]: theme === 'dark' })}
          onClick={() => setTheme('dark')}
          role="tab"
          aria-selected={theme === 'dark'}>
          <HalfMoon />
          <span>Dark</span>
        </button>
      </div>
    </div>
  )
}
