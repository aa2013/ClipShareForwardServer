export interface ReqResponse {
  code: number;
  msg: string;
  data: Record<string, any>;
}

export interface LoginResp {
  token: string;
}

export interface LocalUser {
  username: string;
  token: string;
}

export interface GlobalDialogProps {
  title: string;
  icon?: string;
  iconColor?: string;
  msg: string
  cancelBtnText?: string;
  cancelBtnColor?: string;
  showCancelBtn?: boolean;
  okBtnText?: string;
  okBtnColor?: string;
  onCancel?: () => void;
  onOk?: () => void | Promise<void>;
  persistent?: boolean;
}

export interface GlobalSnackbarProps {
  text: string;
  icon?: string;
  color?: string;
  location?: Anchor
  showAction?: boolean
  actionText?: string;
  actionColor?: string
  onActionClick?: () => void;
  timeout?: number;
}

export interface DevInfo {
  devId: string;
  devName: string;
  platform: string;
  appVersion: string;
}

export interface ConnectionStatus {
  self: DevInfo;
  target?: DevInfo;
  connType: string;
  createTime: string;
  speed: string;
  transferredBytes: string;
  unlimited: boolean
}

export interface ConnectionStatusResp {
  base: ConnectionStatus[],
  dataSync: ConnectionStatus[],
  fileSync: ConnectionStatus[]
}

export const connTypes: (keyof ConnectionStatusResp)[] = ['base', 'dataSync', 'fileSync']

export interface ConnTableItem {
  selfId: string,
  selfName: string,
  platform: string,
  appVersion: string,
  targetId?: string,
  targetName?: string,
  createTime: string,
  speed: string,
  transferredBytes: string,
  unlimited: boolean,
}

export interface LoginSettings {
  loginExpiredSeconds: number
}

export interface FileTransferLimit {
  fileTransferEnabled: boolean,
  fileTransferRateLimit: number
}

export interface UnlimitedDevice {
  id: string,
  name: string
  desc?: string
}

export interface LogConfig {
  memoryBufferSize: number,
}

export interface Log {
  level: string,
  time: string,
  content: string
}

export interface LogResp {
  log: string,
  time: string,
}

export type LocalThemeTypes = 'dark' | 'light' | 'auto'

export interface ThemeMode {
  mode: LocalThemeTypes
  name: string
  icon: string
}

export type LogLevel = 'info' | 'warn' | 'error'
export type SysConfig = LoginSettings & FileTransferLimit & {
  unlimitedDevices: UnlimitedDevice[]
  log: LogConfig,
  publicMode: boolean
}
export type TrafficTrends = 'network' | 'usage' | 'connection'

export interface NetworkChartData {
  dataSync: number;
  fileSync: number;
  time: string;
}

export interface ConnectionChartData {
  baseCnt: number;
  dataSyncCnt: number;
  fileSyncCnt: number;
  time: string;
}

export interface TrafficChartData {
  traffic: number;
  time: string;
}

export interface ChartData {
  netSpeed: NetworkChartData;
  traffic: number;
  connectionCnt: ConnectionChartData;
}

export interface ForcedDisconnectionDto {
  connType: keyof ConnectionStatusResp;
  key: string;
}

export interface PlanType {
  id?: string,
  name: string,
  rate?: number,
  lifespan?: number,
  deviceLimit?: number
  remark?: string
  enable: boolean
}

export interface PlanKey {
  id: number,
  key: string,
  planId: string,
  planName: string,
  useAt: string,
  createdAt: string,
  enable: boolean
  remark?: string
}

export interface GenPlanKeyInfo {
  id: string,
  size: number
}

export interface PageParams {
  pageNum: number,
  pageSize: number,
}

export type PageData<T> = PageParams & {
  total: number
  rows: T[]
}

export interface DrawerMenu {
  text: string
  value: string
  icon: string
  route?: string
  defaultParams?: any[]
  children?: DrawerMenu[]
}

export interface AutoCompletedItem {
  title: string,
  value: string
}

export interface TableSortOptions {
  key: string,
  order: 'asc' | 'desc'
}
