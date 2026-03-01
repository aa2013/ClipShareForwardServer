import {LocalThemeTypes} from "@/types";

export const localTheme = {
  get theme(): LocalThemeTypes {
    return getTheme()
  },
  setTheme(theme: LocalThemeTypes): void {
    if (theme === 'auto') {
      localStorage.removeItem('theme')
      return
    }
    localStorage.setItem('theme', theme)
  },
  clear(): void {
    localStorage.removeItem('theme')
  }
}
const getTheme = (): LocalThemeTypes => {
  const theme = localStorage.getItem("theme");
  if (theme === 'dark' || theme === 'light') return theme
  return 'auto'
}
