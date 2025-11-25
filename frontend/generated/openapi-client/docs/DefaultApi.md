# DefaultApi

All URIs are relative to *http://localhost:8080*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**dashboardV1AuthLoginPost**](DefaultApi.md#dashboardv1authloginpostoperation) | **POST** /dashboard/v1/auth/login | Login with email + password |
| [**dashboardV1PaymentIdReviewPut**](DefaultApi.md#dashboardv1paymentidreviewput) | **PUT** /dashboard/v1/payment/{id}/review | Allows marking a payment as reviewed only by operation role |
| [**dashboardV1PaymentsGet**](DefaultApi.md#dashboardv1paymentsget) | **GET** /dashboard/v1/payments | List of payments |



## dashboardV1AuthLoginPost

> User dashboardV1AuthLoginPost(dashboardV1AuthLoginPostRequest)

Login with email + password

### Example

```ts
import {
  Configuration,
  DefaultApi,
} from '';
import type { DashboardV1AuthLoginPostOperationRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new DefaultApi();

  const body = {
    // DashboardV1AuthLoginPostRequest
    dashboardV1AuthLoginPostRequest: ...,
  } satisfies DashboardV1AuthLoginPostOperationRequest;

  try {
    const data = await api.dashboardV1AuthLoginPost(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **dashboardV1AuthLoginPostRequest** | [DashboardV1AuthLoginPostRequest](DashboardV1AuthLoginPostRequest.md) |  | |

### Return type

[**User**](User.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | return token and user information |  -  |
| **401** | Authentication failed or missing credentials |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## dashboardV1PaymentIdReviewPut

> DashboardV1PaymentIdReviewPut200Response dashboardV1PaymentIdReviewPut(id)

Allows marking a payment as reviewed only by operation role

### Example

```ts
import {
  Configuration,
  DefaultApi,
} from '';
import type { DashboardV1PaymentIdReviewPutRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new DefaultApi(config);

  const body = {
    // string
    id: id_example,
  } satisfies DashboardV1PaymentIdReviewPutRequest;

  try {
    const data = await api.dashboardV1PaymentIdReviewPut(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **id** | `string` |  | [Defaults to `undefined`] |

### Return type

[**DashboardV1PaymentIdReviewPut200Response**](DashboardV1PaymentIdReviewPut200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Payment review message |  -  |
| **401** | Authentication failed or missing credentials |  -  |
| **403** | User is authenticated but not authorized |  -  |
| **404** | Resource not found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## dashboardV1PaymentsGet

> DashboardV1PaymentsGet200Response dashboardV1PaymentsGet(limit, offset, sort, status, id)

List of payments

### Example

```ts
import {
  Configuration,
  DefaultApi,
} from '';
import type { DashboardV1PaymentsGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new DefaultApi(config);

  const body = {
    // number | Limit number of items to return (max 100) (optional)
    limit: 56,
    // number | Offset from start (0-based) (optional)
    offset: 56,
    // string | Comma-separated sort fields. Common patterns: `-created_at` (prefix `-` = desc) `amount` (no prefix `-` = asc)  (optional)
    sort: -created_at,
    // string | status of payment (completed , processing , or failed) (optional)
    status: status_example,
    // string | payment id (optional)
    id: id_example,
  } satisfies DashboardV1PaymentsGetRequest;

  try {
    const data = await api.dashboardV1PaymentsGet(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **limit** | `number` | Limit number of items to return (max 100) | [Optional] [Defaults to `20`] |
| **offset** | `number` | Offset from start (0-based) | [Optional] [Defaults to `0`] |
| **sort** | `string` | Comma-separated sort fields. Common patterns: &#x60;-created_at&#x60; (prefix &#x60;-&#x60; &#x3D; desc) &#x60;amount&#x60; (no prefix &#x60;-&#x60; &#x3D; asc)  | [Optional] [Defaults to `undefined`] |
| **status** | `string` | status of payment (completed , processing , or failed) | [Optional] [Defaults to `undefined`] |
| **id** | `string` | payment id | [Optional] [Defaults to `undefined`] |

### Return type

[**DashboardV1PaymentsGet200Response**](DashboardV1PaymentsGet200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Payment List |  -  |
| **401** | Authentication failed or missing credentials |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

