<template>
  <div :id="id" class="size-full"/>
</template>
<script setup lang="ts">
import {markRaw} from "vue";
import * as echarts from 'echarts';
import {onBeforeUnmount, onMounted, ref, watch} from "vue";
import {storeToRefs} from "pinia";
import {useLocalTheme} from "@/stores/theme";

const {currentTheme} = storeToRefs(useLocalTheme())
const {getAutoTheme} = useLocalTheme()
const props = defineProps<{
  id: string;
  options: echarts.EChartsOption;
}>();
const chart = ref<echarts.ECharts>();
const onResize = () => {
  chart.value?.resize();
}
const initEcharts = () => {
  chart.value?.dispose()
  const theme = currentTheme.value === 'auto' ? getAutoTheme() : currentTheme.value
  chart.value = markRaw(echarts.init(document.getElementById(props.id), theme));
  updateOptions()
}
const updateOptions = () => {
  chart.value?.setOption({...props.options, backgroundColor: 'transparent'});
}
onMounted(() => {
  initEcharts()
  window.addEventListener('resize', onResize);
})
onBeforeUnmount(() => {
  window.removeEventListener('resize', onResize);
});
watch([props], () => {
  updateOptions()
})
watch([currentTheme], () => {
  initEcharts()
})
</script>
<style scoped>

</style>
