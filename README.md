Introduction
============

Sandman is a backend service that helps to decouple highly coupled REST API microservices. It does this by using a message queue (such as SQS) to allow the microservices to communicate with each other asynchronously.

How Sandman works
=================

Sandman works by reading tasks from a message queue and executing the actions specified in the task. A task is a JSON object that defines the actions to be performed and any required parameters or dependencies.

When a task is pulled from the queue, Sandman runs the task in a new goroutine. This allows Sandman to process multiple tasks concurrently and improve the overall performance of the system.

If the task's "request" action is successful, Sandman will delete the task from the queue. If the "request" action fails or times out, Sandman will retry the task a certain number of times (as specified in the "retries" field of the "request" object) before deleting it from the queue.

If a task has any callbacks defined in the "response" object, Sandman will execute those callbacks if the "request" action is successful. If a callback fails, Sandman will delete the main task from the queue and create a new task for the failed callback. This allows the failed callback to be retried, potentially allowing it to succeed on a subsequent attempt.

Tasks
-----

Tasks are the primary way that microservices communicate with each other through Sandman. Tasks are JSON objects that represent actions to be performed by the microservices. Each task has a set of fields that define the action to be performed and any necessary context or data.

Tasks are stored in a message queue (such as SQS) and are processed by Sandman. When a task is processed, Sandman sends an HTTP request to the specified "sub" field of the "request" object within the task. The response of this request is then processed according to the "response" object within the task.

If the "request" is successful and the "response" is as expected, Sandman will delete the task from the queue. If the "request" fails or times out, Sandman will retry the task a certain number of times (as specified in the "retries" field of the "request" object) before deleting it from the queue.

Tasks Definition
----------------

A task in Sandman is a JSON object that represents an action or set of actions to be performed by the service. Each task has a set of fields that define the task's properties and behavior.

#### Required Fields

-   `name`: A string that specifies the name of the task. This field is used for logging and identification purposes.

-   `traceId`: A string that specifies a unique identifier for the task. This field is used for tracing and debugging purposes.

-   `origin`: A string that specifies the source or origin of the task. This field is used for logging and identification purposes.

-   `request`: An object that specifies the HTTP request to be sent by the task. This field is required and must include the following sub-fields:

	-   `method`: A string that specifies the HTTP method of the request (e.g. "GET", "POST", etc.).

	-   `sub`: A string that specifies the URL or endpoint of the request.

	-   `headers`: An object that specifies the HTTP headers to be sent with the request. This field is optional.

	-   `body`: An object that specifies the body of the request. This field is optional.

#### Optional Fields

-   `groupReference`: A string that specifies a reference or identifier for the task's group or set. This field is used for logging and identification purposes and is optional.

-   `sandmanVersion`: A string that specifies the version of Sandman that the task is compatible with. This field is used for logging and identification purposes and is optional.

-   `intent`: A string that specifies the purpose or goal of the task. This field is used for logging and identification purposes and is optional.

-   `description`: A string that provides a more detailed description of the task. This field is used for logging and identification purposes and is optional.

-   `journey`: A string that specifies the journey or flow that the task is a part of. This field is used for logging and identification purposes and is optional.

-   `owner`: A string that specifies the owner or responsible party for the task. This field is used for logging and identification purposes and is optional.

-   `response`: An object that specifies the expected response of the task's "request". This field is optional and, if included, must include the following sub-fields:

	-   `successStatus`: An integer that specifies the expected HTTP status code of the response.

	-   `hasCallBack`: A boolean that specifies whether the task has any callbacks. If set to true, the "response" object must also include a "callbacks" array.

	-   `callbacks`: An array of objects that specify the callbacks to be performed by the task. Each callback object must include a "whenRequestStatus" field and a "request" object.

Architecture
------------

Sandman is designed using the Ports and Adapters pattern, which helps to decouple the internal logic of the service from its external dependencies. The internal logic of Sandman communicates with the external dependencies (such as the message queue and any external APIs or services it communicates with) through the defined ports (interfaces), rather than directly interacting with the dependencies' APIs.

This allows for more flexibility in terms of the dependencies used by Sandman, as the internal logic does not need to be modified when changing or adding new dependencies. It also makes it easier to test and maintain the service, as the internal logic is isolated from the dependencies.

Callbacks
---------

If a task has a "response" object with a "hasCallBack" field set to true, Sandman will process any callbacks specified in the "callbacks" field of the "response" object. Callbacks are used to provide feedback or perform additional actions based on the response of the initial "request".

Each callback has a "whenRequestStatus" field that specifies the HTTP status code of the initial "request" that will trigger the callback. If the initial "request" has a status code that matches the "whenRequestStatus" field of a callback, the callback will be processed.

Like tasks, callbacks are processed by sending an HTTP request to the specified "sub" field of the callback's "request" object. The response of this request is then processed according to the "response" object of the callback.

If the callback's "request" is successful and the "response" is as expected, Sandman will delete the callback from the queue (if applicable). If the callback's "request" fails or times out, Sandman will retry the callback a certain number of times (as specified in the "retries" field of the callback's "request" object) before deleting it from the queue.

Configuration
-------------

The visibility timeout is a setting in a message queue that determines how long a task remains invisible to other consumers after it is read from the queue. In the context of Sandman, the visibility timeout should be set to a value that is long enough to allow Sandman to process the task and its callbacks, but not so long that it causes delays or bottlenecks in the system.

There are a few factors to consider when setting the visibility timeout for Sandman:

1.  The complexity of the task: If the task is complex and requires a lot of processing or external API calls, the visibility timeout should be set to a longer value to allow Sandman to complete the task.

2.  The reliability of the microservices: If the microservices that Sandman communicates with are unreliable or prone to timeouts, the visibility timeout should be set to a longer value to allow for retries.

3.  The desired response time: If fast response times are important for the system, the visibility timeout should be set to a shorter value to allow Sandman to process tasks quickly.

In general, it is a good idea to start with a visibility timeout of around 10-20 seconds and adjust it based on the needs of the system. It is also a good idea to monitor the task processing times and the message queue depths to ensure that the visibility timeout is set appropriately.

Sandman can be configured through a set of environment variables or through a configuration file. The following configuration options are available:

-   `QUEUE_PROVIDER`: Specifies the message queue provider to use (e.g. "sqs").

-   `QUEUE_URL`: Specifies the URL of the message queue.

-   `AWS_REGION`: (Required for SQS queue provider) Specifies the AWS region of the message queue.

-   `AWS_ACCESS_KEY_ID`: (Required for SQS queue provider) Specifies the AWS access key ID for authentication.

-   `AWS_SECRET_ACCESS_KEY`: (Required for SQS queue provider) Specifies the AWS secret access key for authentication.

Deployment
----------

Sandman can be deployed as a containerized service, using a container orchestration platform such as Kubernetes. The following steps can be used to deploy Sandman:

1.  Build the Sandman Docker image.

2.  Push the Docker image to a container registry.

3.  Create a Container runtime deployment that references the Docker image and sets the necessary environment variables or configuration file.

4.  Use the Kubernetes deployment to create one or more replicas of the Sandman service.

Monitoring and Logging
----------------------

Sandman can be monitored and logged through the use of tools such as AWS CloudWatch, Prometheus and Grafana for metrics and ELK stack for logs. The following metrics and logs can be useful for monitoring and debugging Sandman:

-   Task processing time: The amount of time it takes for Sandman to process a task.

-   Task processing success rate: The percentage of tasks that are successfully processed by Sandman.

-   Callback processing success rate: The percentage of callbacks that are successfully processed by Sandman.

-   HTTP request success rate: The percentage of HTTP requests made by Sandman that are successful.

-   HTTP request latency: The amount of time it takes for Sandman to complete an HTTP request.

-   Task and callback JSON: The JSON for tasks and callbacks that are processed by Sandman.

By monitoring and logging these metrics and logs, you can get a better understanding of how Sandman is performing and identify any issues or bottlenecks that may need to be addressed.

Maintenance
-----------

Sandman requires regular maintenance to ensure that it is functioning correctly and efficiently. The following maintenance tasks should be performed on a regular basis:

-   Monitor and log the metrics and logs described in the "Monitoring and Logging" section to identify any issues or bottlenecks.

-   Monitor the message queue for any tasks or callbacks that may be stuck or unable to be processed.

-   Monitor the microservices that communicate with Sandman to ensure that they are functioning correctly and sending and receiving tasks as expected.

-   Perform regular updates to Sandman and the microservices to ensure that they are using the latest versions and patches.

By performing these maintenance tasks regularly, you can help to ensure that Sandman is running smoothly and efficiently, and that any issues are identified and addressed in a timely manner.

Additional Features
-------------------

The following additional features are available in Sandman:

-   "mapToBody" field: Allows you to map data from the initial "request" or "response" to the "request" or "response" of a callback. This can be useful if you want to pass data from the initial "request" or "response" to the callback. To use this feature, you can specify a "mapToBody" array within the "response" object of a task or callback. Each element in the array should have a "queryField" and "targetField" field, which specify the field in the initial "request" or "response" to be mapped and the field in the callback's "request" or "response" to be mapped to, respectively.

For example:

`1Copy code{ 2 "response": { 3 "hasCallBack": true, 4 "callbacks": [ 5 { 6 "whenRequestStatus": 200, 7 "mapToBody": [ 8 { 9 "queryField": "paymentId", 10 "targetField": "paymentReference" 11 } 12 ], 13 "request": { 14 "sub": "https://commission-backend.svc.dev.ind.vm.co.mz/api/v1/payments/3/feedback", 15 "body": { 16 "paymentReference": "" 17 } 18 } 19 } 20 ] 21 } 22} 23`

In this example, the "paymentId" field from the initial "request" or "response" would be mapped to the "paymentReference" field in the callback's "request" body.

Conclusion
----------

Sandman is a powerful backend service that helps to decouple highly coupled REST API microservices. By using the Ports and Adapters pattern and a message queue, Sandman allows microservices to communicate with each other asynchronously, improving the scalability and reliability of the system. With features such as callbacks and the ability to map data between requests and responses, Sandman provides a flexible and adaptable solution for managing tasks and actions within a microservice architecture.

By following best practices for deployment, monitoring, logging, and maintenance, you can ensure that Sandman is running smoothly and efficiently, and can quickly identify and resolve any issues that may arise.
