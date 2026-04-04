<template>
  <v-card class="pa-2 h-[100%] flex flex-col" flat>
    <v-card-title class="d-flex align-center">
      <v-icon icon="mdi-cloud-key-outline" class="mr-2" color="primary" />
      {{ t('planKeys.title') }}
      <v-spacer />
      <div class="w-[150px]">
        <v-autocomplete
          clearable
          color="primary"
          v-model="planId"
          hide-no-data
          hide-details
          :label="t('planKeys.planSelect')"
          density="compact"
          :placeholder="t('planKeys.planSelectPlaceholder')"
          :items="plans"
          variant="outlined"
          :loading="!plans"
        />
      </div>
    </v-card-title>
    <v-divider />
    <div class="flex-1-0 overflow-auto h-0" ref="tableContainer">
      <v-data-table-server
        v-model:page="pageData.page"
        :items-per-page="pageData.itemsPerPage"
        width="100%"
        :items-per-page-text="t('pagination.itemsPerPage')"
        :headers="tableHeaders"
        :items="keys"
        :height="tableHeight || tableContainer?.offsetHeight"
        :no-data-text="t('common.noData')"
        :fixed-header="true"
        :loading="!keys"
        multi-sort
        :items-length="pageData.total"
        @update:sortBy="(sortBy) => sortOptions = sortBy"
        sticky
      >
        <template #item.key="{ item }">
          <v-btn prepend-icon="mdi-eye" class="cursor-pointer" color="primary" variant="text" @click="() => showKeyDialog(item.key)">
            {{ t('common.view') }}
          </v-btn>
        </template>
        <template #item.enable="{ item }">
          <v-switch
            density="comfortable"
            color="primary"
            v-model="item.enable"
            :loading="statusLoading[item.id]"
            hide-details
            inset
            @update:modelValue="(value: boolean | null) => onKeyStatusChanged(item, !!value)"
          />
        </template>
        <template #bottom></template>
        <template #loading>
          <v-skeleton-loader type="table-row@20"></v-skeleton-loader>
        </template>
      </v-data-table-server>
    </div>
    <div class="flex flex-wrap justify-between items-center ga-4 pa-2">
      <div>
        {{ t('common.totalItems', { count: pageData.total }) }}
      </div>
      <div class="flex overflow-auto items-center">
        <div class="min-w-[120px]">
          <v-select
            hide-details
            color="primary"
            density="compact"
            :label="t('pagination.itemsPerPage')"
            v-model="pageData.itemsPerPage"
            :items="pageData.pageItems"
            variant="outlined"
          />
        </div>
        <v-pagination active-color="primary" show-first-last-page v-model="pageData.page" :length="pageCnt" total-visible="5" />
      </div>
    </div>
  </v-card>
</template>
<script setup lang="ts">
import { t } from '@/i18n'
import * as planKeysReq from '@/network/details/planKeys'
import * as planReq from '@/network/details/plan'
import { useGlobalDialog } from '@/stores/dialog'
import { useGlobalSnackbar } from '@/stores/snackbar'
import { AutoCompletedItem, PlanKey, TableSortOptions } from '@/types'
import useClipboard from 'vue-clipboard3'
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'

const { toClipboard } = useClipboard()
const { showGlobalDialog } = useGlobalDialog()
const { showSnackbar } = useGlobalSnackbar()

const tableHeaders = computed<any[]>(() => [
  {
    title: t('common.id'),
    align: 'start',
    key: 'id',
    minWidth: '80px',
  },
  {
    title: t('planKeys.planName'),
    align: 'start',
    key: 'planName',
    minWidth: '115px',
    sortable: false,
  },
  {
    title: t('common.key'),
    align: 'center',
    key: 'key',
    minWidth: '100px',
    sortable: false,
  },
  {
    title: t('common.status'),
    align: 'start',
    key: 'enable',
    minWidth: '100px',
  },
  {
    title: t('planKeys.firstUseTime'),
    align: 'start',
    key: 'useAt',
    minWidth: '206px',
  },
  {
    title: t('common.createdAt'),
    align: 'start',
    key: 'createdAt',
    minWidth: '206px',
  },
  {
    title: t('common.remark'),
    align: 'start',
    key: 'content',
    minWidth: '100px',
    sortable: false,
  },
])

const keys = ref<PlanKey[]>()
const tableContainer = ref<HTMLElement>()
const tableHeight = ref<number>()
const sortOptions = ref<TableSortOptions[]>([])
const onWindowResize = () => {
  tableHeight.value = tableContainer.value?.offsetHeight
}
const pageData = ref({
  page: 1,
  pageItems: [10, 20, 50, 100],
  itemsPerPage: 20,
  total: 0,
})
const route = useRoute()
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
      text: value
        ? t('common.enableResult', { result: res ? t('common.success') : t('common.operationFailed') })
        : t('common.disableResult', { result: res ? t('common.success') : t('common.operationFailed') }),
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
    plans.value = list.map(item => ({ title: item.name, value: item.id! }))
  })
}

const fetchPlanKeys = () => {
  const sorts = JSON.parse(JSON.stringify(sortOptions.value)) as TableSortOptions[]
  sorts.forEach(sort => {
    sort.key = sort.key.replace(/[A-Z]/g, match => `_${match.toLowerCase()}`)
  })
  planKeysReq.list({
    pageNum: pageData.value.page,
    pageSize: pageData.value.itemsPerPage,
    planId: !planId.value || planId.value as unknown === false ? undefined : planId.value,
    sorts,
  }).then(list => {
    pageData.value.total = list.total
    keys.value = list.rows
  })
}

onMounted(() => {
  fetchPlans()
  fetchPlanKeys()
  const id = route.params.planId as string
  planId.value = ((id as unknown as boolean) === false || !id) ? undefined : id
  window.addEventListener('resize', onWindowResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', onWindowResize)
})

const showKeyDialog = (key: string) => {
  showGlobalDialog({
    title: t('planKeys.showKeyTitle'),
    msg: key,
    cancelBtnText: t('common.copy'),
    cancelBtnColor: 'primary',
    showCancelBtn: true,
    onCancel() {
      setTimeout(() => {
        toClipboard(key).then(() => {
          showSnackbar({
            text: t('common.copySucceeded'),
          })
        }).catch(err => {
          showSnackbar({
            text: t('common.copyFailed'),
          }, true)
          console.log(err)
        })
      }, 500)
    },
    okBtnText: t('common.close'),
  })
}
</script>
