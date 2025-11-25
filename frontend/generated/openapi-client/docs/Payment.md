
# Payment


## Properties

Name | Type
------------ | -------------
`id` | string
`merchant` | string
`status` | string
`amount` | string
`createdAt` | Date

## Example

```typescript
import type { Payment } from ''

// TODO: Update the object below with actual values
const example = {
  "id": 1,
  "merchant": Merchant A,
  "status": completed , processing , or failed,
  "amount": alice@example.com,
  "createdAt": null,
} satisfies Payment

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as Payment
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


