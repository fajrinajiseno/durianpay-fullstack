import { describe, vi, it, expect, beforeEach, afterEach } from 'vitest'
import { mountSuspended, mockNuxtImport } from '@nuxt/test-utils/runtime'
import Dashboard from '~/pages/dashboard.vue'
import Table from '@nuxt/ui/components/Table.vue'
import Select from '@nuxt/ui/components/Select.vue'
import Input from '@nuxt/ui/components/Input.vue'

const userMock = vi.fn()
const dashboardV1PaymentsGetMock = vi.fn()
const dashboardV1PaymentIdReviewPutMock = vi.fn()
const toastAddMock = vi.fn()

describe('Dashboard Page', () => {
  beforeEach(() => {
    const { useGeneratedClientMock, useToastMock, useAuthStoreMock } =
      vi.hoisted(() => {
        return {
          useGeneratedClientMock: vi.fn(() => {
            return {
              api: {
                dashboardV1PaymentsGet: dashboardV1PaymentsGetMock,
                dashboardV1PaymentIdReviewPut: dashboardV1PaymentIdReviewPutMock
              }
            }
          }),
          useToastMock: vi.fn(() => {
            return {
              add: toastAddMock
            }
          }),
          useAuthStoreMock: vi.fn(() => {
            return {
              getUser: userMock
            }
          })
        }
      })
    mockNuxtImport('useGeneratedClient', () => {
      return useGeneratedClientMock
    })
    mockNuxtImport('useToast', () => {
      return useToastMock
    })
    mockNuxtImport('useAuthStore', () => {
      return useAuthStoreMock
    })
  })

  afterEach(() => {
    vi.clearAllMocks()
    vi.resetAllMocks()
    userMock.mockClear()
    dashboardV1PaymentsGetMock.mockClear()
    dashboardV1PaymentIdReviewPutMock.mockClear()
    toastAddMock.mockClear()
  })

  it('Success Render', async () => {
    dashboardV1PaymentsGetMock.mockResolvedValueOnce({
      meta: { limit: 10, offset: 0, total: 12 },
      payments: [
        {
          amount: '100',
          createdAt: '2025-11-24T01:10:25+07:00',
          id: '1',
          merchant: 'merchant 1',
          status: 'pending'
        },
        {
          amount: '200',
          createdAt: '2025-11-24T01:10:25+07:00',
          id: '2',
          merchant: 'merchant 2',
          status: 'completed'
        }
      ],
      summary: { completed: 9, failed: 2, pending: 1, total: 12 }
    })
    const page = await mountSuspended(Dashboard, { route: '/dashboard' })
    await page.vm.$nextTick()
    expect(dashboardV1PaymentsGetMock).toBeCalledWith({
      limit: 10,
      offset: 0,
      sort: '-created_at'
    })
    const UTable = page.findComponent(Table)
    expect(UTable.text()).toBe(
      'IDMerchantAmountDateStatusAction1merchant 1IDR 100.0011/24/2025, 1:10:25 AMpending2merchant 2IDR 200.0011/24/2025, 1:10:25 AMcompleted'
    )
    expect(page.find('[data-testid="payment-summary-total"]').text()).toContain(
      '12'
    )
    expect(
      page.find('[data-testid="payment-summary-completed"]').text()
    ).toContain('9')
    expect(
      page.find('[data-testid="payment-summary-failed"]').text()
    ).toContain('2')
    expect(
      page.find('[data-testid="payment-summary-pending"]').text()
    ).toContain('1')
    expect(page.find('[data-testid="payment-review-1"]').exists()).toBeFalsy()
    page.unmount()
  })

  it('Success Click Action Review', async () => {
    userMock.mockImplementation(() => {
      return {
        email: 'operation@test.com',
        role: 'operation',
        token: 'token'
      }
    })
    dashboardV1PaymentsGetMock.mockResolvedValueOnce({
      meta: { limit: 10, offset: 0, total: 12 },
      payments: [
        {
          amount: '100',
          createdAt: '2025-11-24T01:10:25+07:00',
          id: '1',
          merchant: 'merchant 1',
          status: 'pending'
        },
        {
          amount: '200',
          createdAt: '2025-11-24T01:10:25+07:00',
          id: '2',
          merchant: 'merchant 2',
          status: 'completed'
        }
      ],
      summary: { completed: 9, failed: 2, pending: 1, total: 12 }
    })
    dashboardV1PaymentIdReviewPutMock.mockResolvedValueOnce({
      message: 'success review'
    })
    const page = await mountSuspended(Dashboard, { route: '/dashboard' })
    await page.vm.$nextTick()
    const UTable = page.findComponent(Table)
    expect(UTable.text()).toBe(
      'IDMerchantAmountDateStatusAction1merchant 1IDR 100.0011/24/2025, 1:10:25 AMpendingReview2merchant 2IDR 200.0011/24/2025, 1:10:25 AMcompletedReview'
    )
    expect(page.find('[data-testid="payment-review-1"]').exists()).toBeTruthy()
    await page.find('[data-testid="payment-review-1"]').trigger('click')
    await page.vm.$nextTick()
    expect(toastAddMock).toBeCalledWith({
      title: 'success',
      description: 'success review'
    })
    page.unmount()
  })

  it('Success select status', async () => {
    dashboardV1PaymentsGetMock
      .mockResolvedValueOnce({
        meta: { limit: 10, offset: 0, total: 12 },
        payments: [
          {
            amount: '100',
            createdAt: '2025-11-24T01:10:25+07:00',
            id: '1',
            merchant: 'merchant 1',
            status: 'pending'
          },
          {
            amount: '200',
            createdAt: '2025-11-24T01:10:25+07:00',
            id: '2',
            merchant: 'merchant 2',
            status: 'completed'
          }
        ],
        summary: { completed: 9, failed: 2, pending: 1, total: 12 }
      })
      .mockResolvedValueOnce({
        meta: { limit: 10, offset: 0, total: 1 },
        payments: [
          {
            amount: '200',
            createdAt: '2025-11-24T01:10:25+07:00',
            id: '2',
            merchant: 'merchant 2',
            status: 'completed'
          }
        ],
        summary: { completed: 9, failed: 2, pending: 1, total: 12 }
      })
    const page = await mountSuspended(Dashboard, { route: '/dashboard' })
    await page.vm.$nextTick()
    expect(dashboardV1PaymentsGetMock).toHaveBeenLastCalledWith({
      limit: 10,
      offset: 0,
      sort: '-created_at'
    })
    const UTable = page.findComponent(Table)
    expect(UTable.text()).toBe(
      'IDMerchantAmountDateStatusAction1merchant 1IDR 100.0011/24/2025, 1:10:25 AMpending2merchant 2IDR 200.0011/24/2025, 1:10:25 AMcompleted'
    )
    const USelect = page.findComponent(Select)
    page.vm.selectedStatus = 'success'
    USelect.vm.$emit('change')
    await page.vm.$nextTick()
    expect(dashboardV1PaymentsGetMock).toHaveBeenLastCalledWith({
      status: 'success',
      limit: 10,
      offset: 0,
      sort: '-created_at'
    })
    page.unmount()
  })

  it('Success Search ID', async () => {
    vi.useFakeTimers()
    dashboardV1PaymentsGetMock
      .mockResolvedValueOnce({
        meta: { limit: 10, offset: 0, total: 12 },
        payments: [
          {
            amount: '100',
            createdAt: '2025-11-24T01:10:25+07:00',
            id: '1',
            merchant: 'merchant 1',
            status: 'pending'
          },
          {
            amount: '200',
            createdAt: '2025-11-24T01:10:25+07:00',
            id: '2',
            merchant: 'merchant 2',
            status: 'completed'
          }
        ],
        summary: { completed: 9, failed: 2, pending: 1, total: 12 }
      })
      .mockResolvedValueOnce({
        meta: { limit: 10, offset: 0, total: 1 },
        payments: [
          {
            amount: '200',
            createdAt: '2025-11-24T01:10:25+07:00',
            id: '1',
            merchant: 'merchant 2',
            status: 'completed'
          }
        ],
        summary: { completed: 9, failed: 2, pending: 1, total: 12 }
      })
    const page = await mountSuspended(Dashboard, { route: '/dashboard' })
    await page.vm.$nextTick()
    expect(dashboardV1PaymentsGetMock).toHaveBeenLastCalledWith({
      limit: 10,
      offset: 0,
      sort: '-created_at'
    })
    const UTable = page.findComponent(Table)
    expect(UTable.text()).toBe(
      'IDMerchantAmountDateStatusAction1merchant 1IDR 100.0011/24/2025, 1:10:25 AMpending2merchant 2IDR 200.0011/24/2025, 1:10:25 AMcompleted'
    )
    const UInput = page.findComponent(Input)
    UInput.vm.$emit('update:model-value', '1')
    vi.runOnlyPendingTimers()
    await page.vm.$nextTick()
    expect(dashboardV1PaymentsGetMock).toHaveBeenLastCalledWith({
      id: '1',
      limit: 10,
      offset: 0,
      sort: '-created_at'
    })
    page.unmount()
    vi.useRealTimers()
  })

  it('Success Sort Amount', async () => {
    dashboardV1PaymentsGetMock
      .mockResolvedValueOnce({
        meta: { limit: 10, offset: 0, total: 12 },
        payments: [
          {
            amount: '100',
            createdAt: '2025-11-24T01:10:25+07:00',
            id: '1',
            merchant: 'merchant 1',
            status: 'pending'
          },
          {
            amount: '200',
            createdAt: '2025-11-24T01:10:25+07:00',
            id: '2',
            merchant: 'merchant 2',
            status: 'completed'
          }
        ],
        summary: { completed: 9, failed: 2, pending: 1, total: 12 }
      })
      .mockResolvedValueOnce({
        meta: { limit: 10, offset: 0, total: 12 },
        payments: [
          {
            amount: '200',
            createdAt: '2025-11-24T01:10:25+07:00',
            id: '2',
            merchant: 'merchant 2',
            status: 'completed'
          },
          {
            amount: '100',
            createdAt: '2025-11-24T01:10:25+07:00',
            id: '1',
            merchant: 'merchant 1',
            status: 'pending'
          }
        ],
        summary: { completed: 9, failed: 2, pending: 1, total: 12 }
      })
    const page = await mountSuspended(Dashboard, { route: '/dashboard' })
    await page.vm.$nextTick()
    expect(dashboardV1PaymentsGetMock).toHaveBeenLastCalledWith({
      limit: 10,
      offset: 0,
      sort: '-created_at'
    })
    const UTable = page.findComponent(Table)
    expect(UTable.text()).toBe(
      'IDMerchantAmountDateStatusAction1merchant 1IDR 100.0011/24/2025, 1:10:25 AMpending2merchant 2IDR 200.0011/24/2025, 1:10:25 AMcompleted'
    )
    expect(
      page.find('[data-testid="payment-sort-amount"]').exists()
    ).toBeTruthy()
    await page.find('[data-testid="payment-sort-amount"]').trigger('click')
    await page.vm.$nextTick()
    expect(dashboardV1PaymentsGetMock).toHaveBeenLastCalledWith({
      limit: 10,
      offset: 0,
      sort: 'amount'
    })
    page.unmount()
  })

  it('Error render', async () => {
    dashboardV1PaymentsGetMock.mockRejectedValueOnce({
      code: '500',
      message: 'error'
    })
    const page = await mountSuspended(Dashboard, { route: '/dashboard' })
    await page.vm.$nextTick()
    expect(dashboardV1PaymentsGetMock).toHaveBeenLastCalledWith({
      limit: 10,
      offset: 0,
      sort: '-created_at'
    })
    expect(page.find('[data-testid="payment-error"]').exists()).toBeTruthy()
    page.unmount()
  })

  it('No Data', async () => {
    dashboardV1PaymentsGetMock.mockResolvedValueOnce({
      meta: { limit: 10, offset: 0, total: 0 },
      payments: [],
      summary: { completed: 0, failed: 0, pending: 0, total: 0 }
    })
    const page = await mountSuspended(Dashboard, { route: '/dashboard' })
    await page.vm.$nextTick()
    expect(dashboardV1PaymentsGetMock).toHaveBeenLastCalledWith({
      limit: 10,
      offset: 0,
      sort: '-created_at'
    })
    expect(page.find('[data-testid="payment-nodata"]').exists()).toBeTruthy()
    page.unmount()
  })
})
