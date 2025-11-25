
# PaymentSummary


## Properties

Name | Type
------------ | -------------
`total` | number
`completed` | number
`failed` | number
`pending` | number

## Example

```typescript
import type { PaymentSummary } from ''

// TODO: Update the object below with actual values
const example = {
  "total": 40,
  "completed": 20,
  "failed": 10,
  "pending": 10,
} satisfies PaymentSummary

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as PaymentSummary
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


