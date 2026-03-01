import {defineStore} from 'pinia'
import {computed, ref} from "vue";
import {GlobalDialogProps} from "@/types";

export const useGlobalDialog = defineStore('global-dialog', () => {

  const dialogShow = ref(false)
  const okLoading = ref(false)
  const defaultProps: GlobalDialogProps = {
    title: '',
    icon: 'mdi-information-outline',
    iconColor: 'primary',
    msg: '',
    cancelBtnText: '取消',
    cancelBtnColor: undefined,
    showCancelBtn: false,
    okBtnText: '确定',
    okBtnColor: 'primary',
    onCancel: () => {
    },
    onOk: () => {
    },
    persistent: false,
  }
  const dialogProps = ref<GlobalDialogProps>(defaultProps)
  const onDialogClose = () => {
    dialogShow.value = false
    const onCancel = dialogProps.value.onCancel
    onCancel && onCancel()
  }
  const onDialogOk = () => {
    const onOk = dialogProps.value.onOk
    if (!onOk) {
      return
    }
    const promise = onOk()
    if (promise instanceof Promise) {
      okLoading.value = true
      promise.finally(() => {
        okLoading.value = false
        dialogShow.value = false
      })
    } else {
      okLoading.value = false
      dialogShow.value = false
    }
  }
  const isOkLoading = computed(() => okLoading)
  const globalDialogProps = computed(() => dialogProps.value)
  const showGlobalDialog = (props: GlobalDialogProps) => {
    okLoading.value = false
    dialogProps.value = Object.assign({...defaultProps}, props)
    dialogShow.value = true
  }
  return {
    dialogShow, onDialogClose, onDialogOk, showGlobalDialog, globalDialogProps, isOkLoading
  }
})
