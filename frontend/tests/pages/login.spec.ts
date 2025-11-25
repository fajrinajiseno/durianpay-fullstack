import { vi, it, expect, beforeEach, afterEach } from 'vitest'
import { mountSuspended, mockNuxtImport } from '@nuxt/test-utils/runtime'
import Login from '~/pages/login.vue'
import AuthForm from '@nuxt/ui/components/AuthForm.vue'

const dashboardV1AuthLoginPostMockMock = vi.fn()
const toastAddMock = vi.fn()

beforeEach(() => {
  const { useGeneratedClientMock, usehandleErrorMock, useToastMock } =
    vi.hoisted(() => {
      return {
        useGeneratedClientMock: vi.fn(() => {
          return {
            api: { dashboardV1AuthLoginPost: dashboardV1AuthLoginPostMockMock }
          }
        }),
        usehandleErrorMock: vi.fn(() => {
          return {
            code: '500',
            message: 'error'
          }
        }),
        useToastMock: vi.fn(() => {
          return {
            add: toastAddMock
          }
        })
      }
    })
  mockNuxtImport('useGeneratedClient', () => {
    return useGeneratedClientMock
  })
  mockNuxtImport('usehandleError', () => {
    return usehandleErrorMock
  })
  mockNuxtImport('useToast', () => {
    return useToastMock
  })
})

afterEach(() => {
  vi.clearAllMocks()
})

it('Success Login', async () => {
  dashboardV1AuthLoginPostMockMock.mockResolvedValueOnce({
    email: 'operation@test.com',
    role: 'cs',
    token: 'token'
  })
  const page = await mountSuspended(Login, { route: '/login' })
  const UAuthForm = page.findComponent(AuthForm)
  UAuthForm.vm.$emit('submit', {
    data: { email: 'operation@test.com', password: 'password' }
  })
  await page.vm.$nextTick()
  expect(dashboardV1AuthLoginPostMockMock).toBeCalledWith({
    dashboardV1AuthLoginPostRequest: {
      email: 'operation@test.com',
      password: 'password'
    }
  })
  expect(toastAddMock).toBeCalledWith({
    title: 'success login'
  })
})

it('Error Login', async () => {
  dashboardV1AuthLoginPostMockMock.mockRejectedValueOnce({
    response: {
      text: vi.fn().mockResolvedValueOnce({
        code: '500',
        message: 'error'
      })
    }
  })
  const page = await mountSuspended(Login, { route: '/login' })
  const UAuthForm = page.findComponent(AuthForm)
  UAuthForm.vm.$emit('submit', {
    data: { email: 'operation@test.com', password: 'password' }
  })
  await page.vm.$nextTick()
  expect(dashboardV1AuthLoginPostMockMock).toBeCalledWith({
    dashboardV1AuthLoginPostRequest: {
      email: 'operation@test.com',
      password: 'password'
    }
  })
  await page.vm.$nextTick()
  expect(toastAddMock).toBeCalledWith({
    title: '500',
    description: 'error',
    color: 'error'
  })
})
