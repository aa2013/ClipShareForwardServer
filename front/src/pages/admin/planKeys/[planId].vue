<template>
  <v-card class="pa-2 h-[100%] flex flex-col" flat>
    <v-card-title class="d-flex align-center">
      <v-icon icon="mdi-cloud-key-outline" class="mr-2" color="primary"/>
      密钥管理
      <v-spacer/>
      <div class="w-[150px]">
        <v-autocomplete
          clearable
          color="primary"
          v-model="planId"
          hide-no-data
          hide-details
          label="套餐选择"
          density="compact"
          placeholder="请选择套餐"
          :items="plans"
          variant="outlined"
          :loading="!plans"
        />
      </div>

    </v-card-title>
    <v-divider/>
    <div class="flex-1-0 overflow-auto h-0" ref="tableContainer">
      <v-data-table-server
        v-model:page="pageData.page"
        :items-per-page="pageData.itemsPerPage"
        width="100%"
        items-per-page-text="分页数量 "
        :headers="tableHeaders"
        :items="keys"
        :height="tableHeight||tableContainer?.offsetHeight"
        no-data-text="暂无数据"
        :fixed-header="true"
        :loading="!keys"
        multi-sort
        :items-length="pageData.total"
        @update:sortBy="(sortBy)=>sortOptions=sortBy"
        sticky
      >
        <template v-slot:item.key="{ item }">
          <v-btn prepend-icon="mdi-eye" class="cursor-pointer"
                 color="primary" variant="text"
                 @click="()=>showKeyDialog(item.key)"
          >
            查看
          </v-btn>
        </template>
        <template v-slot:item.enable="{ item }">
          <v-switch
            density="comfortable"
            color="primary"
            v-model="item.enable"
            :loading="statusLoading[item.id]"
            hide-details
            inset
            @update:modelValue="(value:boolean|null)=>onKeyStatusChanged(item,!!value)"
          />
        </template>
        <template v-slot:bottom>

        </template>
        <template v-slot:loading>
          <v-skeleton-loader type="table-row@20"></v-skeleton-loader>
        </template>
      </v-data-table-server>
    </div>
    <div class="flex flex-wrap justify-between items-center ga-4 pa-2">
      <div>
        共 {{ pageData.total }} 条数据
      </div>
      <div class="flex overflow-auto items-center">
        <div class="min-w-[120px]">
          <v-select
            hide-details
            color="primary"
            density="compact"
            label="分页数量"
            v-model="pageData.itemsPerPage"
            :items="pageData.pageItems"
            variant="outlined"
          />
        </div>
        <v-pagination
          active-color="primary"
          show-first-last-page
          v-model="pageData.page"
          :length="pageCnt"
          total-visible="5"
        />

      </div>
    </div>
  </v-card>
</template>
<script setup lang="ts">
import {computed, onMounted, onUnmounted, ref, watch} from "vue";
import {AutoCompletedItem, PlanKey, TableSortOptions} from "@/types";
import {useRoute} from 'vue-router';
import * as planKeysReq from '@/network/details/planKeys'
import * as planReq from '@/network/details/plan'
import {useGlobalDialog} from '@/stores/dialog'
import {useGlobalSnackbar} from '@/stores/snackbar'
import useClipboard from 'vue-clipboard3'

const {toClipboard} = useClipboard()

const {showGlobalDialog} = useGlobalDialog()
const {showSnackbar} = useGlobalSnackbar()
const tableHeaders = [
  {
    title: 'id',
    align: 'start',
    key: 'id',
    minWidth: '80px',
  },
  {
    title: '套餐名称',
    align: 'start',
    key: 'planName',
    minWidth: '115px',
    sortable: false
  },
  {
    title: '密钥',
    align: 'center',
    key: 'key',
    minWidth: '100px',
    sortable: false
  },
  {
    title: '状态',
    align: 'start',
    key: 'enable',
    minWidth: '100px'
  },
  {
    title: '首次使用时间',
    align: 'start',
    key: 'useAt',
    minWidth: '206px'
  },
  {
    title: '创建时间',
    align: 'start',
    key: 'createdAt',
    minWidth: '206px'
  },
  {
    title: '备注',
    align: 'start',
    key: 'content',
    minWidth: '100px',
    sortable: false
  },
] as any[]
const keys = ref<PlanKey[]>()
const tableContainer = ref<HTMLElement>()
const tableHeight = ref<number>()
const sortOptions = ref<TableSortOptions[]>([])
const onWindowResize = () => {
  tableHeight.value = tableContainer.value?.offsetHeight
}
const pageData = ref({
  page: 1,
  pageItems: [10, 20, 50, 100,],
  itemsPerPage: 20,
  total: 0
})
const route = useRoute();
const pageCnt = computed(() => Math.ceil(pageData.value.total / pageData.value.itemsPerPage))
const planId = ref<string>()
const plans = ref<AutoCompletedItem[]>()
const statusLoading = ref<Record<string, any>>({})
watch(() => route.params.planId, (newId) => {
  planId.value = newId as string
  plans.value = undefined
  fetchPlans()
  fetchPlanKeys()
})
watch(() => [pageData.value.page, pageData.value.itemsPerPage, sortOptions.value], () => {
  fetchPlanKeys()
})
const onKeyStatusChanged = (item: PlanKey, value: boolean) => {
  statusLoading.value[item.id] = true
  planKeysReq.updateStatus(item.id, value).then(res => {
    showSnackbar({
      text: `${value ? '启用' : '停用'}${res ? '成功' : '失败'}`
    }, !res)
    if (!res) {
      item.enable = !value
    }
  }).catch(err => {
    console.log(err)
    item.enable = !value
  }).finally(() => {
    delete statusLoading.value[item.id]
  })
}
const fetchPlans = () => {
  planReq.getPlans().then(list => {
    plans.value = list.map(item => ({title: item.name, value: item.id!}))
  })
}
const fetchPlanKeys = () => {
  const sorts = JSON.parse(JSON.stringify(sortOptions.value)) as TableSortOptions[]
  sorts.forEach(sort => {
    //将所有大写转为下划线和小写字母，如 planId -> plan_id
    sort.key = sort.key.replace(/[A-Z]/g, match => `_${match.toLowerCase()}`)
  })
  planKeysReq.list({
    pageNum: pageData.value.page,
    pageSize: pageData.value.itemsPerPage,
    planId: !planId.value || planId.value as unknown == false ? undefined : planId.value,
    sorts: sorts,
  }).then(list => {
    pageData.value.total = list.total
    keys.value = list.rows
  })
}
onMounted(() => {
  fetchPlans()
  fetchPlanKeys()
  const id = route.params.planId as string
  planId.value = ((id as unknown as boolean) == false || !id) ? undefined : id
  window.addEventListener('resize', onWindowResize)
})
onUnmounted(() => {
  window.removeEventListener('resize', onWindowResize)
})
const showKeyDialog = (key: string) => {
  showGlobalDialog({
    title: "密钥",
    msg: key,
    cancelBtnText: '复制',
    cancelBtnColor: 'primary',
    showCancelBtn: true,
    onCancel() {
      setTimeout(() => {
        //由于模态框焦点问题，延迟复制
        toClipboard(key).then(() => {
          showSnackbar({
            text: "复制成功"
          })
        }).catch(err => {
          showSnackbar({
            text: "复制失败",
          }, true)
          console.log(err)
        })
      }, 500)
    },
    okBtnText: '关闭'
  })
}
</script>
<style scoped>

</style>
