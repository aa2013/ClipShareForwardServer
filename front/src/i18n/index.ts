import { computed, ref } from 'vue'

export type AppLocale = 'en-US' | 'zh-CN'

interface MessageObject {
  [key: string]: string | MessageObject
}

type MessageValue = string | MessageObject

type Messages = Record<AppLocale, MessageObject>

const STORAGE_KEY = 'locale'

export const localeOptions: Array<{ value: AppLocale, label: string }> = [
  { value: 'en-US', label: 'English (US)' },
  { value: 'zh-CN', label: '简体中文' },
]

const messages: Messages = {
  'en-US': {
    common: {
      add: 'Add',
      action: 'Action',
      cancel: 'Cancel',
      close: 'Close',
      confirm: 'Confirm',
      copy: 'Copy',
      createdAt: 'Created At',
      delete: 'Delete',
      description: 'Description',
      edit: 'Edit',
      enabled: 'Enabled',
      englishUs: 'English (US)',
      generate: 'Generate',
      id: 'ID',
      key: 'Key',
      loading: 'Loading',
      login: 'Login',
      logout: 'Logout',
      mode: 'Mode',
      name: 'Name',
      no: 'No',
      noData: 'No data',
      ok: 'OK',
      operationFailed: 'Operation failed',
      plan: 'Plan',
      platform: 'Platform',
      publicMode: 'Public Mode',
      remark: 'Remark',
      required: 'Required.',
      save: 'Save',
      search: 'Search',
      seconds: 'seconds',
      status: 'Status',
      submit: 'Submit',
      success: 'Success',
      time: 'Time',
      totalItems: 'Total {count} items',
      unlimited: 'Unlimited',
      view: 'View',
      yes: 'Yes',
      chineseSimplified: 'Simplified Chinese',
      username: 'Username',
      password: 'Password',
      deviceId: 'Device ID',
      deviceName: 'Device Name',
      deviceDescription: 'Device Description',
      appVersion: 'App Version',
      targetDevice: 'Target Device',
      connectionTime: 'Connection Time',
      currentSpeed: 'Current Speed',
      trafficUsed: 'Traffic Used',
      whitelist: 'Whitelist',
      refresh: 'Refresh',
      saveResult: 'Update {result}',
      enableResult: 'Enable {result}',
      disableResult: 'Disable {result}',
      generateResult: 'Generate {result}',
      addResult: 'Add {result}',
      updateResult: 'Update {result}',
      copySucceeded: 'Copied successfully',
      copyFailed: 'Copy failed',
    },
    language: {
      label: 'Language',
      select: 'Select Language',
    },
    theme: {
      auto: 'Auto',
      light: 'Light',
      dark: 'Dark',
    },
    validation: {
      required: 'Required.',
      cannotBeNegative: 'Cannot be negative.',
      cannotBeDecimal: 'Cannot be decimal.',
      maximum: 'Maximum {value}',
      deviceIdExists: 'The device ID already exists',
    },
    loginPage: {
      title: 'Forward Server',
      usernameRequired: 'Please enter a username',
      passwordRequired: 'Please enter a password',
      loginButton: 'Login',
    },
    nav: {
      connectionInfo: 'Connections',
      planManagement: 'Plan Management',
      plans: 'Plans',
      keyManagement: 'Key Management',
      keyVerification: 'Key Verification',
      settings: 'Settings',
      logs: 'Logs',
    },
    drawer: {
      logoutConfirm: 'Are you sure you want to log out?',
    },
    overview: {
      title: 'Overview',
      trafficTrends: 'Traffic Trends',
      deviceConnectionStatus: 'Device Connection Status',
      totalConnections: 'Total Connections',
      deviceCount: 'Devices',
      syncConnections: 'Sync Connections',
      fileConnections: 'File Connections',
      deviceConnectionsTab: 'Device Connections',
      dataSyncTab: 'Data Sync',
      fileSyncTab: 'File Sync',
      clickToViewDeviceId: 'Click to view device ID',
      disconnectConnection: 'Disconnect connection',
      deviceIdTitle: '{name} Device ID',
      disconnectPrompt: 'Disconnect this connection: {id}?',
      disconnectRelatedConnections: 'Related connections will also be disconnected!',
      disconnectedConnection: 'Disconnected this connection',
      disconnectedConnectionAndRelated: 'Disconnected this connection and related connections',
    },
    trends: {
      network: 'Network',
      usage: 'Usage',
      connection: 'Connections',
    },
    charts: {
      totalConnections: 'Total Connections',
      deviceConnections: 'Device Connections',
      dataSyncConnections: 'Data Sync Connections',
      fileSyncConnections: 'File Sync Connections',
      totalSpeed: 'Total Speed',
      dataSync: 'Data Sync',
      fileSync: 'File Sync',
      currentConnectionTraffic: 'Current Connection Traffic',
      unitPerSecond: 'Unit: {unit}/s',
      unit: 'Unit: {unit}',
    },
    logs: {
      title: 'Runtime Logs',
      level: 'Log Level',
      content: 'Log Content',
    },
    settings: {
      loginSettings: 'Login',
      loginExpiredSeconds: 'Login timeout (seconds)',
      modeSelection: 'Mode Selection',
      fileTransferLimit: 'File Transfer Limit (effective in public mode)',
      allowFileTransfer: 'Allow file sync',
      fileTransferRateLimit: 'Rate limit (KB/s)',
      whitelistDevices: 'Whitelist Devices',
      deviceCard: 'Device {index}',
      addDevice: 'Add Device',
      logs: 'Logs',
      logMemoryBufferSize: 'In-memory log buffer size',
      deleteWhitelistDeviceTitle: 'Prompt',
      deleteWhitelistDeviceMessage: 'Delete this whitelist device: {name} ({id})?',
    },
    verifyKey: {
      title: 'Key Verification',
      inputLabel: 'Key',
      inputPlaceholder: 'Please enter a key',
      submit: 'Verify',
      notFound: 'This key does not exist',
    },
    planKeys: {
      title: 'Key Management',
      planSelect: 'Select Plan',
      planSelectPlaceholder: 'Please select a plan',
      itemsPerPage: 'Items per page',
      planName: 'Plan Name',
      firstUseTime: 'First Use Time',
      showKeyTitle: 'Key',
    },
    plans: {
      deviceLimit: 'Simultaneous devices:',
      rateLimit: 'Rate limit:',
      lifespan: 'Validity:',
      days: 'days',
      devicesUnit: 'devices',
      generateKeys: 'Generate Keys',
      clickToDisable: 'Click to Disable',
      clickToEnable: 'Click to Enable',
      expand: 'Expand',
      collapse: 'Collapse',
      addPlan: 'Add Plan',
      editPlan: 'Edit Plan',
      addOrEditPlan: '{action} Plan',
      planName: 'Plan Name',
      rateLimitLabel: 'Rate Limit',
      validityDays: 'Validity (days)',
      deviceLimitLabel: 'Concurrent Device Limit',
      remark: 'Remark',
      submitAction: 'Submit {action}',
      generateKeysTitle: 'Generate Keys',
      generateCount: 'Generate Count',
    },
    pagination: {
      itemsPerPage: 'Items per page',
    },
    notFound: {
      redirectToLogin: 'Redirecting to login in {seconds} seconds...',
    },
    footer: {
      copyright: 'Copyright © {year} ClipShare. All Rights Reserved.',
    },
    errors: {
      title: 'Error',
    },
  },
  'zh-CN': {
    common: {
      add: '添加',
      action: '操作',
      cancel: '取消',
      close: '关闭',
      confirm: '确认',
      copy: '复制',
      createdAt: '创建时间',
      delete: '删除',
      description: '描述',
      edit: '修改',
      enabled: '启用',
      englishUs: 'English (US)',
      generate: '生成',
      id: 'ID',
      key: '密钥',
      loading: '加载中',
      login: '登录',
      logout: '退出登录',
      mode: '模式',
      name: '名称',
      no: '否',
      noData: '暂无数据',
      ok: '确定',
      operationFailed: '操作失败',
      plan: '套餐',
      platform: '平台',
      publicMode: '公开模式',
      remark: '备注',
      required: '必填项。',
      save: '保存',
      search: '搜索',
      seconds: '秒',
      status: '状态',
      submit: '提交',
      success: '成功',
      time: '时间',
      totalItems: '共 {count} 条数据',
      unlimited: '无限制',
      view: '查看',
      yes: '是',
      chineseSimplified: '简体中文',
      username: '用户名',
      password: '密码',
      deviceId: '设备 ID',
      deviceName: '设备名称',
      deviceDescription: '设备描述',
      appVersion: 'App 版本',
      targetDevice: '目标设备',
      connectionTime: '连接时间',
      currentSpeed: '当前速度',
      trafficUsed: '已用流量',
      whitelist: '白名单',
      refresh: '刷新',
      saveResult: '更新{result}',
      enableResult: '启用{result}',
      disableResult: '停用{result}',
      generateResult: '生成{result}',
      addResult: '添加{result}',
      updateResult: '更新{result}',
      copySucceeded: '复制成功',
      copyFailed: '复制失败',
    },
    language: {
      label: '语言',
      select: '选择语言',
    },
    theme: {
      auto: '自动',
      light: '亮色模式',
      dark: '暗色模式',
    },
    validation: {
      required: '必填项。',
      cannotBeNegative: '不能为负数。',
      cannotBeDecimal: '不能为小数。',
      maximum: '最大值为 {value}',
      deviceIdExists: '设备 ID 已存在',
    },
    loginPage: {
      title: 'Forward Server',
      usernameRequired: '请输入用户名',
      passwordRequired: '请输入密码',
      loginButton: '登录',
    },
    nav: {
      connectionInfo: '连接信息',
      planManagement: '套餐管理',
      plans: '套餐管理',
      keyManagement: '密钥管理',
      keyVerification: '密钥验证',
      settings: '配置修改',
      logs: '运行日志',
    },
    drawer: {
      logoutConfirm: '是否确认退出登录？',
    },
    overview: {
      title: '概览',
      trafficTrends: '流量趋势',
      deviceConnectionStatus: '设备连接状态',
      totalConnections: '总连接数',
      deviceCount: '设备数',
      syncConnections: '同步连接数',
      fileConnections: '文件连接数',
      deviceConnectionsTab: '设备连接',
      dataSyncTab: '数据同步',
      fileSyncTab: '文件同步',
      clickToViewDeviceId: '点击查看设备 ID',
      disconnectConnection: '断开连接',
      deviceIdTitle: '{name} 设备 ID',
      disconnectPrompt: '是否断开该连接: {id}？',
      disconnectRelatedConnections: '相关连接都会被断开！',
      disconnectedConnection: '已断开该连接',
      disconnectedConnectionAndRelated: '已断开该连接及相关连接',
    },
    trends: {
      network: '网络',
      usage: '使用量',
      connection: '连接数',
    },
    charts: {
      totalConnections: '总连接数',
      deviceConnections: '设备连接数',
      dataSyncConnections: '数据同步连接数',
      fileSyncConnections: '文件同步连接数',
      totalSpeed: '总速度',
      dataSync: '数据同步',
      fileSync: '文件同步',
      currentConnectionTraffic: '当前连接使用量',
      unitPerSecond: '单位：{unit}/s',
      unit: '单位：{unit}',
    },
    logs: {
      title: '运行日志',
      level: '日志级别',
      content: '日志内容',
    },
    settings: {
      loginSettings: '登录',
      loginExpiredSeconds: '登录超时时间（秒）',
      modeSelection: '模式选择',
      fileTransferLimit: '文件传输限制（公开模式下生效）',
      allowFileTransfer: '允许文件同步',
      fileTransferRateLimit: '速率限制（KB/s）',
      whitelistDevices: '白名单设备',
      deviceCard: '设备 {index}',
      addDevice: '添加设备',
      logs: '日志',
      logMemoryBufferSize: '日志内存缓冲区大小',
      deleteWhitelistDeviceTitle: '提示',
      deleteWhitelistDeviceMessage: '是否删除该白名单设备: {name}（{id}）？',
    },
    verifyKey: {
      title: '密钥验证',
      inputLabel: '密钥',
      inputPlaceholder: '请输入密钥',
      submit: '提交验证',
      notFound: '该密钥不存在',
    },
    planKeys: {
      title: '密钥管理',
      planSelect: '套餐选择',
      planSelectPlaceholder: '请选择套餐',
      itemsPerPage: '分页数量',
      planName: '套餐名称',
      firstUseTime: '首次使用时间',
      showKeyTitle: '密钥',
    },
    plans: {
      deviceLimit: '设备同时连接数:',
      rateLimit: '限速:',
      lifespan: '有效期:',
      days: '天',
      devicesUnit: '台',
      generateKeys: '生成密钥',
      clickToDisable: '点击停用',
      clickToEnable: '点击启用',
      expand: '展开',
      collapse: '收起',
      addPlan: '新增套餐',
      editPlan: '修改套餐',
      addOrEditPlan: '{action}套餐',
      planName: '套餐名称',
      rateLimitLabel: '限速速率',
      validityDays: '有效期（天）',
      deviceLimitLabel: '同时使用设备数量',
      remark: '备注',
      submitAction: '提交{action}',
      generateKeysTitle: '生成密钥',
      generateCount: '生成数量',
    },
    pagination: {
      itemsPerPage: '分页数量',
    },
    notFound: {
      redirectToLogin: '{seconds} 秒后跳转登录...',
    },
    footer: {
      copyright: '版权所有 © {year} ClipShare。保留所有权利。',
    },
    errors: {
      title: '错误',
    },
  },
}

const getByPath = (obj: MessageObject, key: string): string | undefined => {
  const value = key.split('.').reduce<MessageValue | undefined>((acc, part) => {
    if (!acc || typeof acc === 'string') return acc
    return acc[part]
  }, obj)
  return typeof value === 'string' ? value : undefined
}

const resolveLocale = (): AppLocale => {
  if (typeof window === 'undefined') {
    return 'en-US'
  }
  const stored = localStorage.getItem(STORAGE_KEY)
  if (stored === 'en-US' || stored === 'zh-CN') {
    return stored
  }
  const browserLocale = navigator.language.toLowerCase()
  if (browserLocale.startsWith('zh')) {
    return 'zh-CN'
  }
  return 'en-US'
}

export const currentLocale = ref<AppLocale>(resolveLocale())

const syncDocumentLang = (locale: AppLocale) => {
  if (typeof document !== 'undefined') {
    document.documentElement.lang = locale
  }
}

syncDocumentLang(currentLocale.value)

export const setLocale = (locale: AppLocale) => {
  currentLocale.value = locale
  if (typeof window !== 'undefined') {
    localStorage.setItem(STORAGE_KEY, locale)
  }
  syncDocumentLang(locale)
}

export const t = (key: string, params: Record<string, string | number> = {}) => {
  const message = getByPath(messages[currentLocale.value], key)
    ?? getByPath(messages['en-US'], key)
    ?? key

  return message.replace(/\{(\w+)\}/g, (_, name: string) => String(params[name] ?? `{${name}}`))
}

export const useLocale = () => ({
  locale: computed(() => currentLocale.value),
  localeOptions: computed(() => localeOptions),
})
