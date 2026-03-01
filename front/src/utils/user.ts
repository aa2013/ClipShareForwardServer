import {LocalUser} from "@/types";

export const local = {
  get username(): string {
    return getUser().username
  },
  get token(): string {
    return getUser().token
  },
  get user(): LocalUser {
    return getUser()
  },
  clearUser(): void {
    localStorage.removeItem('user')
  }
}
const getUser = (): LocalUser => {
  const userStr = localStorage.getItem("user");
  if (!userStr) {
    return {token: '', username: ''}
  }
  try {
    return JSON.parse(userStr!) as LocalUser
  } catch (_) {
    return {token: '', username: ''}
  }
}
