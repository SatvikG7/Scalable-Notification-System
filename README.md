# Scalable Notification Service

## Overview

The Scalable Notification Service is designed to handle high volumes of notifications efficiently, with a focus on reliability, prioritization, and rate-limiting. It supports multiple notification channels like Email, SMS, and Push notifications, and processes them based on priority levels (High, Medium, Low).

## Project Structure

```
/config                  # Configuration files
/internal                # Private application and library code
.env                     # Environment configuration
.gitignore               # Git ignore rules
go.mod                   # Go module file
go.sum                   # Go dependencies
main.go                  # Main entry point of the application
README.md                # Project documentation
```

### Directory Details

- **/config**: Stores configuration files that are necessary for the application to run, including environment variables.
- **/internal**: Contains private application and library code that should not be exposed outside of this project.
- **main.go**: The main entry point of the application, where the server and worker pools are initialized.

## Features Implemented

### 1. Initialization and Setup
- **Basic HTTP Server Setup**: Using the Fiber framework, a simple HTTP server has been initialized with routes for the root path and a health check.
- **Environment Variables**: Configuration is managed using a `.env` file, ensuring that sensitive information is kept secure.

### 2. Message Queue Integration
- **RabbitMQ**: Implemented RabbitMQ for message queuing with different queues for different notification channels and priorities.
- **Producer Logic**: Notifications are sent to different queues based on their channel (e.g., Email, SMS, Push) and priority (e.g., Low, Medium, High).

### 3. Worker Pools and Rate Limiting
- **Worker Pools**: Defined worker pools for each queue, with the ability to reallocate workers based on queue activity.
- **Rate Limiting**: Implemented rate limiting for each worker to ensure external services are not overwhelmed.

### 4. Error Handling and Retry Mechanism
- **Retry Logic**: If a notification fails to send, it is retried based on specific rules. Unsuccessful notifications are logged to PostgreSQL for further analysis.
- **Logging**: Implemented comprehensive logging for tracking errors and system performance.

### 5. Worker Reallocation
- **Dynamic Worker Allocation**: Workers are dynamically allocated to different queues based on their load and priority. Idle queues have their workers reallocated to active queues.
- **Scheduler**: A scheduler balances the priority of queues, ensuring that high-priority notifications are processed promptly without overwhelming any single priority level.

## Future Scope

### 1. Real-time Monitoring with Prometheus and Grafana
- **Metrics Collection**: Integration with Prometheus for collecting metrics related to the performance of the notification service.
- **Visualization**: Setting up Grafana dashboards for real-time monitoring and visualization of key metrics, such as message throughput, worker utilization, and error rates.

## Installation and Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/scalable-notification-service.git
   cd scalable-notification-service
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Configure environment variables:
   - Create a `.env` file based on the `.env.example` provided in the repository.

4. Run the service:
   ```bash
   go run main.go
   ```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Acknowledgements

- Special thanks to the open-source community for the tools and frameworks used in this project.
