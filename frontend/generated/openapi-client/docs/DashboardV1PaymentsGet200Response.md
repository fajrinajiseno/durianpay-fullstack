
# DashboardV1PaymentsGet200Response


## Properties

Name | Type
------------ | -------------
`meta` | [PaginationMeta](PaginationMeta.md)
`summary` | [PaymentSummary](PaymentSummary.md)
`payments` | [Array&lt;Payment&gt;](Payment.md)

## Example

```typescript
import type { DashboardV1PaymentsGet200Response } from ''

// TODO: Update the object below with actual values
const example = {
  "meta": null,
  "summary": null,
  "payments": null,
} satisfies DashboardV1PaymentsGet200Response

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as DashboardV1PaymentsGet200Response
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


