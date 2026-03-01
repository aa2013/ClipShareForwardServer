<template>
  <BaseChart id="connection-cnt-chart" :options="options"/>
</template>
<script setup lang="ts">
import {computed, defineModel} from 'vue'

import BaseChart from "@/components/charts/BaseChart.vue";
import * as echarts from "echarts";
import {ConnectionChartData} from "@/types";

const data = defineModel<ConnectionChartData[]>({
  get(value) {
    if (!value) {
      return []
    }
    return value
  }
})
const options = computed<echarts.EChartsOption>(() => {
  const totalCntData = data.value!.map(item => item.baseCnt+item.dataSyncCnt+item.fileSyncCnt)
  const baseCntData = data.value!.map(item => item.baseCnt)
  const dataSyncCntData = data.value!.map(item => item.dataSyncCnt)
  const fileSyncCntData = data.value!.map(item => item.fileSyncCnt)
  return {
    tooltip: {
      show: true,
      trigger: 'axis',
    },
    legend: {
      data: ['总连接数','设备连接数', '数据同步连接数', '文件同步连接数']
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: data.value!.map(item => item.time)
    },
    yAxis: {
      type: 'value',
      min: 0,
      max: Math.ceil(Math.max(...totalCntData) * 1.2),
    },
    series: [
      {
        name: '总连接数',
        type: 'line',
        data: totalCntData
      },
      {
        name: '设备连接数',
        type: 'line',
        data: baseCntData
      },
      {
        name: '数据同步连接数',
        type: 'line',
        data: dataSyncCntData
      },
      {
        name: '文件同步连接数',
        type: 'line',
        data: fileSyncCntData
      },
    ]
  }
})
</script>

<style scoped>

</style>
