<template>
  <div class="size-full">
    <div class="my-2">
      <div class="flex items-center">
        <v-icon icon="mdi-chart-arc" color="primary" class="mr-2"/>
        <span class="font-weight-bold text-xl">概览</span>
      </div>
      <div class="h-[1px] bg-[#efefef] mt-2"/>
    </div>
    <v-row class="py-4 px-2">
      <v-col v-for="(item,i) in statusCards" :key="i" cols="6" sm="3" md="3" lg="3" xl="3"
             class="flex justify-between">
        <div class="h-[60px]">
          <div class="font-bold text-lg p-2 flex items-center">
            <v-icon :icon="item.icon" size="small" color="primary" class="mr-2"/>
            {{ item.name }}
          </div>
          <div class="font-bold text-xl pl-9">
            {{ item.value }}
          </div>
        </div>
      </v-col>
    </v-row>
    <div class="my-5">
      <div class="flex justify-between">
        <div class="flex items-center">
          <v-icon icon="mdi-chart-bell-curve-cumulative" color="primary" class="mr-2"/>
          <span class="font-weight-bold text-xl">流量趋势</span>
        </div>
        <div class="flex items-center">
          <traffic-trends-btn-group v-model="trafficTrends" v-if="!isMobile"/>
        </div>
      </div>
      <div class="h-[1px] bg-[#efefef] mt-2"/>
    </div>
    <traffic-trends-btn-group v-model="trafficTrends" v-if="isMobile" class="justify-center"/>
    <div class="py-4 px-2 position-relative h-[380px]">
      <network-chart
        v-model="netSpeedChartData"
        class="h-[350px] position-absolute"
        :style="{
        opacity:trafficTrends=='network'?1:0,
        zIndex:trafficTrends=='network'?1:-1
      }"/>
      <traffic-chart
        v-model="trafficChartData"
        class="h-[350px] position-absolute"
        :style="{
        opacity:trafficTrends=='usage'?1:0,
        zIndex:trafficTrends=='usage'?1:-1
      }"/>
      <connections-cnt-chart
        v-model="connectionCntChartData"
        class="h-[350px] position-absolute"
        :style="{
        opacity:trafficTrends=='connection'?1:0,
        zIndex:trafficTrends=='connection'?1:-1
      }"/>
    </div>
    <div class="my-5">
      <div class="flex items-center">
        <v-icon icon="mdi-cellphone-link" color="primary" class="mr-2"/>
        <span class="font-weight-bold text-xl">设备连接状态</span>
      </div>
      <div class="h-[1px] bg-[#efefef] mt-2"/>
    </div>
    <div class="flex flex-col">
      <v-tabs
        v-model="tab"
        color="primary"
        direction="horizontal"
      >
        <v-tab prepend-icon="mdi-account" text="设备连接" value="base"/>
        <v-tab prepend-icon="mdi-lock" text="数据同步" value="dataSync"/>
        <v-tab prepend-icon="mdi-access-point" text="文件同步" value="fileSync"/>
      </v-tabs>

      <v-tabs-window v-model="tab" class="flex-1">
        <v-tabs-window-item v-for="connType in connTypes"
                            :value="connType" :key="connType">
          <v-data-table
            width="100%"
            items-per-page-text="分页数量 "
            :custom-filter="filterOnlyCapsText"
            :headers="tableHeaders"
            :items="tables[connType]"
            no-data-text="暂无数据"
            :fixed-header="true"
            sticky
          >
            <template v-slot:item.unlimited="{ item }">
              <v-chip
                :color="item.unlimited ? 'green' : 'grey'"
                :text="item.unlimited ? '是' : '否'"
                class="text-uppercase"
                size="small"
                label
              />
            </template>
            <template v-slot:item.selfName="{ item }">
              <div class="flex items-center">
                <v-tooltip
                  :text="item.selfName"
                  location="top"
                >
                  <template v-slot:activator="{ props }">
                    <div v-bind="props" class="whitespace-nowrap overflow-hidden w-[80px] truncate ">
                      {{ item.selfName }}
                    </div>
                  </template>

                </v-tooltip>
                <v-tooltip
                  text="点击查看设备id"
                  location="top"
                >
                  <template v-slot:activator="{ props }">
                    <v-btn v-bind="props" size="small"
                           color="blue"
                           variant="text"
                           elevation="0"
                           icon="mdi-identifier"
                           @click="()=>showDevIdDialog(item.selfId,item.selfName)"
                    />
                  </template>

                </v-tooltip>
              </div>
            </template>
            <template v-slot:item.targetName="{ item }">
              <div class="flex items-center">
                <v-tooltip
                  :text="item.targetName"
                  location="top"
                >
                  <template v-slot:activator="{ props }">
                    <div v-bind="props" class="whitespace-nowrap overflow-hidden w-[80px] truncate ">
                      {{ item.targetName }}
                    </div>
                  </template>

                </v-tooltip>
                <v-tooltip
                  text="点击查看设备id"
                  location="top"
                >
                  <template v-slot:activator="{ props }">
                    <v-btn v-bind="props" size="small"
                           color="blue" v-if="item.targetId"
                           variant="text"
                           elevation="0"
                           icon="mdi-identifier"
                           @click="()=>showDevIdDialog(item.targetId!,item.targetName!)"
                    />
                  </template>

                </v-tooltip>
              </div>
            </template>
            <template v-slot:item.actions="{item}">
              <v-tooltip
                text="断开连接"
                location="bottom"
              >
                <template v-slot:activator="{ props }">
                  <v-btn v-bind="props" size="small"
                         color="red"
                         variant="text"
                         elevation="0"
                         icon="mdi-close-circle-outline"
                         @click="()=>disconnect(connType,item)"
                  />
                </template>

              </v-tooltip>
            </template>
          </v-data-table>
        </v-tabs-window-item>
      </v-tabs-window>
    </div>
  </div>
</template>
<script setup lang="ts">
import {
  ConnectionChartData,
  ConnectionStatusResp,
  ConnTableItem,
  connTypes,
  NetworkChartData,
  TrafficChartData,
  TrafficTrends
} from "@/types.d";

const {showGlobalDialog} = useGlobalDialog()
const {showSnackbar} = useGlobalSnackbar()
import NetworkChart from "@/components/charts/NetworkChart.vue";
import {computed, onMounted, onUnmounted, ref, watchEffect} from "vue";
import * as connReq from "@/network/details/connection";
import ConnectionsCntChart from "@/components/charts/ConnectionsCntChart.vue";
import TrafficTrendsBtnGroup from "@/components/TrafficTrendsBtnGroup.vue";
import {useDisplay} from 'vuetify'
import dayjs from "dayjs";
import TrafficChart from "@/components/charts/TrafficChart.vue";
import {useGlobalDialog} from "@/stores/dialog";
import {useGlobalSnackbar} from "@/stores/snackbar";

const {mobile} = useDisplay()

const trafficTrends = ref<TrafficTrends>('network')
const tab = ref<keyof ConnectionStatusResp>('base')
const isMobile = ref(mobile.value)
const connectionCntChartData = ref<ConnectionChartData[]>([])
const netSpeedChartData = ref<NetworkChartData[]>([])
const trafficChartData = ref<TrafficChartData[]>([])
watchEffect(() => {
  isMobile.value = mobile.value
})
const chartDataLen = 10
const tableHeaders = [
  {
    title: '连接设备',
    align: 'start',
    key: 'selfName',
    minWidth: '115px'
  },
  {
    title: '平台',
    align: 'start',
    key: 'platform',
    minWidth: '105px'
  },
  {
    title: 'App版本',
    align: 'start',
    key: 'appVersion',
    minWidth: '135px'
  },
  {
    title: '目标设备',
    align: 'start',
    key: 'targetName',
    minWidth: '135px'
  },
  {
    title: '连接时间',
    align: 'start',
    key: 'createTime',
    minWidth: '175px'
  },
  {
    title: '当前速度',
    align: 'start',
    key: 'speed',
    minWidth: '115px'
  },
  {
    title: '已用流量',
    align: 'start',
    key: 'transferredBytes',
    minWidth: '115px'
  },
  {
    title: '白名单',
    align: 'start',
    key: 'unlimited',
    minWidth: '95px'
  },
  {
    title: '操作',
    key: 'actions',
    value: 'actions',
    align: 'start',
    minWidth: '85px',
    fixed: true,
    sortable: false
  },
] as any[]
const connStatus = ref<ConnectionStatusResp>({
  base: [],
  dataSync: [],
  fileSync: []
})
let connStatusTimer: number
let fetchChartDataTimer: number
const statusCards = computed(() => {
  return [
    {
      icon: 'mdi-connection',
      name: '总连接数',
      value: (connStatus.value?.base.length ?? 0) + (connStatus.value?.dataSync.length ?? 0) + (connStatus.value?.fileSync.length ?? 0),
      color: 'primary'
    },
    {
      icon: 'mdi-cellphone-link',
      name: '设备数',
      value: connStatus.value?.base.length ?? 0,
      color: 'primary'
    },
    {
      icon: 'mdi-sync',
      name: '同步连接数',
      value: connStatus.value?.dataSync.length ?? 0,
      color: 'primary'
    },
    {
      icon: 'mdi-file-arrow-up-down-outline',
      name: '文件连接数',
      value: connStatus.value?.fileSync.length ?? 0,
      color: 'primary'
    },
  ]
})
const filterOnlyCapsText = (value: string, query: string, item: any) => {
  return value != null &&
    query != null &&
    true &&
    value.toString().toLocaleUpperCase().indexOf(query) !== -1
}
onMounted(() => {
  fetchConnectionStatus()
  fetchChartsData()
  fetchChartDataTimer = setInterval(fetchChartsData, 1000)
  connStatusTimer = setInterval(fetchConnectionStatus, 1000)
})
const fetchChartsData = () => {
  connReq.getChartsData().then(data => {
    const connCnt = data.connectionCnt;
    const traffic = data.traffic;
    const netSpeed = data.netSpeed;
    const time = dayjs().format('HH:mm:ss')
    connectionCntChartData.value = [...connectionCntChartData.value, {...connCnt, time}].slice(-1 * chartDataLen);
    trafficChartData.value = [...trafficChartData.value, {traffic, time}].slice(-1 * chartDataLen)
    netSpeedChartData.value = [...netSpeedChartData.value, {...netSpeed, time}].slice(-1 * chartDataLen)
  })
}
const tables = computed<Record<keyof ConnectionStatusResp, ConnTableItem[]>>(() => {
  const unlimitedDeviceIds: string[] = []
  const baseConns: ConnTableItem[] = connStatus.value.base.map(item => {
    if (item.unlimited) {
      unlimitedDeviceIds.push(item.self.devId)
    }
    return {
      appVersion: item.self.appVersion,
      createTime: item.createTime,
      platform: item.self.platform,
      selfId: item.self.devId,
      selfName: item.self.devName,
      speed: item.speed,
      targetId: item.target?.devId,
      targetName: item.target?.devName,
      transferredBytes: item.transferredBytes,
      unlimited: item.unlimited,
    } as ConnTableItem
  })
  const dataSyncConns: ConnTableItem[] = connStatus.value.dataSync.map(item => ({
    appVersion: item.self.appVersion,
    createTime: item.createTime,
    platform: item.self.platform,
    selfId: item.self.devId,
    selfName: item.self.devName,
    speed: item.speed,
    targetId: item.target?.devId,
    targetName: item.target?.devName,
    transferredBytes: item.transferredBytes,
    unlimited: unlimitedDeviceIds.includes(item.self.devId)
  } as ConnTableItem))
  const fileSyncConns: ConnTableItem[] = connStatus.value.fileSync.map(item => ({
    appVersion: item.self.appVersion,
    createTime: item.createTime,
    platform: item.self.platform,
    selfId: item.self.devId,
    selfName: item.self.devName,
    speed: item.speed,
    targetId: item.target?.devId,
    targetName: item.target?.devName,
    transferredBytes: item.transferredBytes,
    unlimited: unlimitedDeviceIds.includes(item.self.devId)
  } as ConnTableItem))
  return {
    base: baseConns,
    dataSync: dataSyncConns,
    fileSync: fileSyncConns,
  }
})

onUnmounted(() => {
  clearInterval(connStatusTimer)
  clearInterval(fetchChartDataTimer)
})
const fetchConnectionStatus = () => {
  connReq.getConnectionStatus().then(data => {
    connStatus.value = data
  })
}
const showDevIdDialog = (id: string, devName: string) => {
  showGlobalDialog({
    title: devName + " 设备 id",
    msg: id,
  })
}
const disconnect = (connType: keyof ConnectionStatusResp, item: ConnTableItem) => {
  const isBase = connType == 'base'
  showGlobalDialog({
    title: "提示",
    msg: `是否断开该连接: ${item.selfId}？${isBase ? '<br>相关连接都会被断开！' : ''}`,
    showCancelBtn: true,
    onOk() {
      const key = isBase ? item.selfId : `${item.selfId}->${item.targetId}`
      connReq.forcedDisconnection({connType, key}).then(res => {
        if (res) {
          fetchConnectionStatus()
          fetchChartsData()
          showSnackbar({
            text: `已断开该连接${isBase ? '及相关连接' : ''}`
          })
        } else {
          showSnackbar({
            text: `操作失败`
          }, true)
        }
      })
    },
  })
}
</script>
<style scoped>

</style>
