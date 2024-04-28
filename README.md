# Moving-Window

## Overview

The Request Counter Service is a simple HTTP server designed to count the number of requests received within the last 60 seconds. This document outlines the design and implementation details of the service.

## Assumption
- The current requests will be counted.
- If the service remains down for a long time before restarting successfully, the downtime duration will not be counted. In an extreme case, if the downtime exceeds 60 seconds, the count will be returned as 0.

## Architecture

The service is implemented in Go and uses an in-memory data structure to store the count of requests. It provides a RESTful API endpoint to retrieve the current count. The service also persists the count to disk to recover from restarts.

### Components

- **Server**: The core component that handles incoming HTTP requests, maintains request counts, and serves the API endpoint.
- **Buckets Array**: A fixed-size array that stores the count of requests for each second within a 60-second sliding window.
- **Mutex**: A synchronization primitive used to ensure thread-safe access to the buckets array.
- **Ticker**: A time.Ticker used to trigger periodic operations such as resetting expired buckets and persisting data to disk.

## API Specification

### Endpoint

- `GET /api/requests`: Retrieves the total number of requests received in the last 60 seconds.

### Response Format

The response is a JSON object with the following structure:

```json
{
  "status": "Success",
  "code": 200,
  "data": {
    "total_requests": 42
  }
}
```

- `status`: A string indicating the result of the API call. It can be "Success" or "Error".
- `code`: An HTTP status code representing the result of the API call.
- `data`: An object containing the actual response data. Here, it contains the `total_requests` field with the count of requests.

## Implementation Details

### Request Counting

The service divides time into 60 one-second buckets. Each incoming request increments the count in the current second's bucket. The total count is the sum of all buckets.

### Sliding Window

Every second, the service moves the window by resetting the count in the bucket that falls out of the 60-second range and subtracting its value from the total count.

### Persistence

The service periodically writes the current state (buckets array and total count) to a JSON file. This allows the service to recover the last known count after a restart.

### Concurrency

Access to shared data is protected by a mutex to prevent race conditions. Atomic operations are used where possible to maintain consistency of the total count.

## Running the Service

The service can be started by running the main function, which sets up the HTTP server and binds the API endpoint to the appropriate handler function.

## Error Handling

Errors during persistence are logged but do not interrupt the operation of the service. Client-facing errors are communicated through the `status` and `code` fields in the JSON response.

## Future Work

- Data statistics may have errors, and there can be data loss in persistent files. Consider using a finer time granularity bucket, such as milliseconds.
- The assumption mentioned above does not take into account clock drift or backward adjustments.

## Conclusion

The Request Counter Service provides a simple and efficient way to track the number of HTTP requests over a sliding window of time. Its RESTful API and JSON response format make it easy to integrate with other services or monitoring tools.
