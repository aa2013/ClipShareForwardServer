import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { t } from '@/i18n'
import { GlobalSnackbarProps } from '@/types'

export const useGlobalSnackbar = defineStore('global-snackbar', () => {
  const snackbarShow = ref(false)

  const getDefaultProps = (): GlobalSnackbarProps => ({
    text: '',
    icon: 'mdi-information-outline',
    color: 'primary',
    location: 'top',
    showAction: false,
    actionText: t('common.close'),
    actionColor: 'white',
    onActionClick: () => {},
    timeout: 1000,
  })

  const snackbarProps = ref<GlobalSnackbarProps>(getDefaultProps())
  const globalSnackbarProps = computed(() => snackbarProps.value)

  const onActionClick = () => {
    snackbarShow.value = false
    const cb = snackbarProps.value.onActionClick
    cb && cb()
  }

  const showSnackbar = (props: GlobalSnackbarProps, error?: boolean) => {
    if (error === true) {
      props.color = 'red'
    } else if (error === false) {
      props.color = 'primary'
    }
    snackbarProps.value = Object.assign({ ...getDefaultProps() }, props)
    snackbarShow.value = true
  }

  return {
    snackbarShow,
    showSnackbar,
    onActionClick,
    globalSnackbarProps,
  }
})
