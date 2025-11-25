<template>
  <div>
    <h1 class="text-2xl font-bold mb-4">Payments</h1>
    <div class="w-full space-y-4 pb-4">
      <div v-if="isShowFilter" class="flex justify-between items-center">
        <div class="flex items-center">
          <UInput
            v-model="paymentId"
            :loading="isLoading"
            :disabled="isLoading"
            placeholder="Filter Payment ID..."
            class="mr-2"
            data-testid="payment-search-id"
            @update:model-value="searchPaymentId"
          />
          <USelect
            v-model="selectedStatus"
            :loading="isLoading"
            :disabled="isLoading"
            placeholder="Select status"
            :items="statusOption"
            data-testid="payment-select-status"
            class="w-48 mr-2"
            @change="updateSelectedStatus"
          />
          <UButton
            v-if="isShowReset"
            color="neutral"
            variant="outline"
            data-testid="payment-reset"
            @click="resetFilter"
            >Reset Filter</UButton
          >
        </div>
        <div class="flex">
          <div class="mr-2"><b>Payment Summary:</b></div>
          <div class="mr-2">
            <UBadge color="info" variant="subtle" class="capitalize mr-1"
              >Total</UBadge
            >
            <b data-testid="payment-summary-total">: {{ summary?.total }}</b>
          </div>
          <div class="mr-2">
            <UBadge color="success" variant="subtle" class="capitalize mr-1"
              >completed</UBadge
            >
            <b data-testid="payment-summary-completed"
              >: {{ summary?.completed }}</b
            >
          </div>
          <div class="mr-2">
            <UBadge color="error" variant="subtle" class="capitalize mr-1"
              >failed</UBadge
            >
            <b data-testid="payment-summary-failed">: {{ summary?.failed }}</b>
          </div>
          <div>
            <UBadge color="neutral" variant="subtle" class="capitalize mr-1"
              >pending</UBadge
            >
            <b data-testid="payment-summary-pending"
              >: {{ summary?.pending }}</b
            >
          </div>
        </div>
      </div>
      <UTable
        ref="table"
        :loading="isLoading"
        :data="payments"
        :columns="columns"
        class="flex-1"
      >
        <template #empty>
          <div v-if="error">
            <UEmpty
              title="Error"
              :description="emptyText"
              data-testid="payment-error"
            >
              <template #actions>
                <UButton
                  icon="i-lucide-refresh-cw"
                  :loading="isLoading"
                  :disabled="isLoading"
                  size="xl"
                  class="mt-2"
                  @click="refresh"
                  >Refresh</UButton
                >
              </template>
            </UEmpty>
          </div>
          <div v-else-if="!isLoading && payments?.length === 0">
            <UEmpty title="No data found" description="please reset filter">
              <template #actions>
                <UButton
                  icon="i-lucide-refresh-cw"
                  :loading="isLoading"
                  :disabled="isLoading"
                  size="xl"
                  class="mt-2"
                  data-testid="payment-nodata"
                  @click="reset"
                  >Reset</UButton
                >
              </template>
            </UEmpty>
          </div>
          <div v-else>
            {{ emptyText }}
          </div>
        </template>
      </UTable>

      <div
        v-if="isShowingPagination"
        class="flex justify-center border-t border-default pt-4"
      >
        <UPagination
          :default-page="pagination.pageIndex + 1"
          :items-per-page="pagination.pageSize"
          :total="total"
          @update:page="(p) => updatePage(p)"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Payment } from '../../generated/openapi-client'
import type { TableColumn } from '@nuxt/ui'
import { debounce } from 'lodash-es'

enum Status {
  completed = 'completed',
  failed = 'failed',
  pending = 'pending'
}

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})

const route = useRoute()
const { api } = useGeneratedClient()
const auth = useAuthStore()
const toast = useToast()

const UBadge = resolveComponent('UBadge')
const UButton = resolveComponent('UButton')
const statusOption = ref([Status.pending, Status.completed, Status.failed])
const paymentId = ref('')
const selectedStatus = ref()
const offset = ref(0)
const limit = ref(10)
const sort = ref('-created_at')

const {
  data: paymentData,
  error,
  status,
  refresh
} = await useAsyncData(
  route.fullPath,
  async () => {
    const paymentResponse = await api.dashboardV1PaymentsGet({
      ...(paymentId.value ? { id: paymentId.value } : {}),
      ...(selectedStatus.value ? { status: selectedStatus.value } : {}),
      offset: offset.value,
      limit: limit.value,
      sort: sort.value
    })
    return { paymentResponse }
  },
  { server: false }
)

const pagination = ref({
  pageIndex: offset.value / limit.value,
  pageSize: limit.value
})

const isLoading = computed(() => {
  return status.value === 'idle' || status.value === 'pending'
})

const payments = computed(() => {
  return paymentData?.value?.paymentResponse.payments?.map((item) => ({
    ...item,
    action: auth.getUser()?.role === 'operation'
  }))
})

const summary = computed(() => {
  return paymentData?.value?.paymentResponse.summary
})

const total = computed(() => {
  return paymentData.value?.paymentResponse.meta?.total
})

const isShowFilter = computed(() => {
  return !isLoading.value && !error.value
})

const isShowReset = computed(() => {
  return paymentId.value || selectedStatus.value
})

const isShowingPagination = computed(() => {
  return !isLoading.value && !error.value && total.value
})

const emptyText = computed(() => {
  if (isLoading.value) {
    return 'Loading Data...'
  } else if (error.value) {
    return error.value?.message
  } else {
    return 'No Data'
  }
})

function updatePage(page: number) {
  const pageIndex = page - 1
  offset.value = pageIndex * limit.value
  pagination.value.pageIndex = pageIndex
  refresh()
}

const searchPaymentId = debounce(() => {
  offset.value = 0
  pagination.value.pageIndex = 0
  refresh()
}, 250)

function updateSelectedStatus() {
  offset.value = 0
  pagination.value.pageIndex = 0
  refresh()
}

function updateSort(newSort: string) {
  sort.value = newSort
  refresh()
}

function resetFilter() {
  paymentId.value = ''
  selectedStatus.value = undefined
  offset.value = 0
  pagination.value.pageIndex = 0
  refresh()
}

function reset() {
  offset.value = 0
  sort.value = '-created_at'
  paymentId.value = ''
  selectedStatus.value = undefined
  pagination.value.pageIndex = 0
  refresh()
}

async function handleClickReview(id: string) {
  try {
    const data = await api.dashboardV1PaymentIdReviewPut({
      id
    })
    toast.add({ title: 'success', description: data.message })
  } catch (error) {
    const errorParsed = await usehandleError(error)
    toast.add({
      title: `${errorParsed.code}`,
      description: errorParsed.message
    })
  }
}

const columns: TableColumn<Payment>[] = [
  {
    accessorKey: 'id',
    header: 'ID'
  },
  {
    accessorKey: 'merchant',
    header: 'Merchant'
  },
  {
    accessorKey: 'amount',
    header: ({ column }) => {
      const isSorted = column.getIsSorted()

      return h(UButton, {
        color: 'neutral',
        variant: 'ghost',
        label: 'Amount',
        icon: isSorted
          ? isSorted === 'asc'
            ? 'i-lucide-arrow-up-narrow-wide'
            : 'i-lucide-arrow-down-wide-narrow'
          : 'i-lucide-arrow-up-down',
        class: '-mx-2.5 text-right',
        'data-testid': 'payment-sort-amount',
        onClick: () => {
          column.toggleSorting(column.getIsSorted() === 'asc')
          updateSort(column.getIsSorted() === 'asc' ? 'amount' : '-amount')
        }
      })
    },
    cell: ({ row }) => {
      const amount = Number.parseFloat(row.getValue('amount'))
      const formatted = new Intl.NumberFormat('en-US', {
        style: 'currency',
        currency: 'IDR'
      }).format(amount)
      return h('div', { class: 'text-right font-medium' }, formatted)
    }
  },
  {
    accessorKey: 'createdAt',
    header: ({ column }) => {
      const isSorted = column.getIsSorted()

      return h(UButton, {
        color: 'neutral',
        variant: 'ghost',
        label: 'Date',
        icon: isSorted
          ? isSorted === 'asc'
            ? 'i-lucide-arrow-up-narrow-wide'
            : 'i-lucide-arrow-down-wide-narrow'
          : 'i-lucide-arrow-up-down',
        class: '-mx-2.5',
        'data-testid': 'payment-sort-created_at',
        onClick: () => {
          column.toggleSorting(column.getIsSorted() === 'asc')
          updateSort(
            column.getIsSorted() === 'asc' ? 'created_at' : '-created_at'
          )
        }
      })
    },
    cell: ({ row }) => {
      return new Date(row.getValue('createdAt')).toLocaleString()
    }
  },
  {
    accessorKey: 'status',
    header: 'Status',
    cell: ({ row }) => {
      const color = {
        [Status.completed]: 'success' as const,
        [Status.failed]: 'error' as const,
        [Status.pending]: 'neutral' as const
      }[row.getValue('status') as string]

      return h(UBadge, { class: 'capitalize', variant: 'subtle', color }, () =>
        row.getValue('status')
      )
    }
  },
  {
    accessorKey: 'action',
    header: 'Action',
    cell: ({ row }) => {
      return row.getValue('action')
        ? h(
            UButton,
            {
              'data-testid': `payment-review-${row.getValue('id')}`,
              onClick: () => handleClickReview(row.getValue('id'))
            },
            () => 'Review'
          )
        : ''
    }
  }
]
</script>
