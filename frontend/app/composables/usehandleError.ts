import type { ErrorContext, ModelError } from '../../generated/openapi-client'

export async function usehandleError(error: unknown): Promise<ModelError> {
  const text = await (error as ErrorContext).response?.text()
  if (text) {
    const errorParsed: ModelError = JSON.parse(text)
    return errorParsed
  }
  return { code: 0, message: '' }
}
