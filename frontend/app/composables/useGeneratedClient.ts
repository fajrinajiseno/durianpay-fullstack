import { Configuration, DefaultApi } from '../../generated/openapi-client' // path may vary
import { useRuntimeConfig } from '#imports'
import { useAuthStore } from '../store/auth'

export function useGeneratedClient() {
  const config = useRuntimeConfig()
  const auth = useAuthStore()
  const apiBase = config.public.apiBase || 'http://localhost:8080'
  const cfg = new Configuration({
    basePath: apiBase,
    headers: {
      Authorization: `Bearer ${auth.user?.token}`
    }
  })
  const api = new DefaultApi(cfg)

  return { api }
}
