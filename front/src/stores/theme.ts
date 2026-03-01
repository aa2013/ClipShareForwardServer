import {defineStore} from "pinia";
import {LocalThemeTypes} from "@/types";
import {computed, ref, watch} from "vue";
import {useTheme} from "vuetify";

export const useLocalTheme = defineStore('local-theme', () => {

  const theme = useTheme()

  const getCurrentTheme = (): LocalThemeTypes => {
    const theme = localStorage.getItem("theme");
    if (theme === 'dark' || theme === 'light') return theme
    return 'auto'
  }
  const getAutoTheme = () => {
    const theme = getCurrentTheme()
    if (theme != 'auto') return theme
    // 获取当前小时
    const currentHour = new Date().getHours();
    // 判断是白天还是晚上
    const isDaytime = currentHour >= 6 && currentHour < 18;
    return isDaytime ? 'light' : "dark"
  }
  const localTheme = ref<LocalThemeTypes>(getCurrentTheme())
  const currentTheme = computed(() => localTheme.value)
  const currentThemeIcon = computed(() => {
    if (localTheme.value === 'light') {
      return 'mdi-brightness-5'
    } else if (localTheme.value === 'dark') {
      return 'mdi-brightness-2'
    }
    return 'mdi-brightness-auto'
  })

  const setTheme = (newTheme: LocalThemeTypes) => {
    localStorage.setItem('theme', newTheme)
    localTheme.value = newTheme
  }
  const clearTheme = () => {
    setTheme('auto')
  }
  watch([localTheme], (newV, oldV) => {
    applyTheme(newV[0])
  })
  const applyTheme = (apply?: LocalThemeTypes) => {
    apply = apply ?? getCurrentTheme()
    theme.global.name.value = apply === 'auto' ? getAutoTheme() : getCurrentTheme()
  }
  return {
    currentTheme, currentThemeIcon, setTheme, clearTheme, applyTheme, getAutoTheme
  }
})
