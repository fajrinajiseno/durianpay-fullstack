<template>
  <div
    v-show="!auth.user"
    class="flex flex-col items-center justify-center gap-4 p-4"
  >
    <UPageCard class="w-full max-w-md">
      <UAuthForm
        :schema="schema"
        title="Login"
        description="Enter your credentials to access dashboard."
        icon="i-lucide-user"
        :fields="fields"
        :submit="{ label: 'Submit', block: true }"
        @submit="onSubmit"
      />
    </UPageCard>
  </div>
</template>

<script setup lang="ts">
import * as z from 'zod'
import type { FormSubmitEvent, AuthFormField } from '@nuxt/ui'
import { useAuthStore } from '~/store/auth'

definePageMeta({
  middleware: 'auth'
})

const auth = useAuthStore()
const toast = useToast()
const router = useRouter()

const fields: AuthFormField[] = [
  {
    name: 'email',
    type: 'email',
    label: 'Email',
    placeholder: 'Enter your email',
    required: true
  },
  {
    name: 'password',
    label: 'Password',
    type: 'password',
    placeholder: 'Enter your password',
    required: true
  }
]

const schema = z.object({
  email: z.email('Invalid email'),
  password: z.string('Password is required')
})

type Schema = z.output<typeof schema>

async function onSubmit(payload: FormSubmitEvent<Schema>) {
  const { api } = useGeneratedClient()
  try {
    const data = await api.dashboardV1AuthLoginPost({
      dashboardV1AuthLoginPostRequest: {
        email: payload.data.email,
        password: payload.data.password
      }
    })
    auth.setUser({
      email: data.email!,
      role: data.role!,
      token: data.token!
    })
    toast.add({ title: 'success login' })
    router.push({
      path: '/dashboard'
    })
  } catch (error) {
    const errorParsed = await usehandleError(error)
    toast.add({
      title: `${errorParsed.code}`,
      description: errorParsed.message,
      color: 'error'
    })
    console.log(errorParsed)
  }
}
</script>
