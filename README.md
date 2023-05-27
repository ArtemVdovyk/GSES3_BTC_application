# Description
The АРI service that will allow:

- find out the current exchange rate of bitcoin (BTC) in hryvnia (UAH)

- sign an email to receive information on changing the course

- a request that will send the current course to all subscribed users.

# Usage
## Run locally:
1. Create .env file inside **"GSES3_BTC_application"** folder
2. Specify necesary variables in **".env"** file
    ```
    EMAIL="<your-google-account>"
    EMAIL_PASSWORD="<your-google-account-app-password>"
    PORT="8080"
    ```
3. Ensure you are in the **"GSES3_BTC_application"** folder.

4. Run the following commands:
    ```
    go run .
    ```

## Run in Docker:
1. Create .env file inside **"GSES3_BTC_application"** folder
2. Specify necesary variables in **".env"** file
    ```
    EMAIL="<your-google-account>"
    EMAIL_PASSWORD="<your-google-account-app-password>"
    PORT="8080"
    ```
3. Ensure you are in the root folder with Dockerfile.
4. Run the following commands:
    ```
    docker build -t <image-name> .
    docker run -p 8080:8080 <image-name>

    ```

# Endpoints
Url example: **"http://localhost:8080/api/"**

Method **Get**:

1. *"/rate"* - Get the current exchange rate of bitcoin (BTC) in hryvnia (UAH).

Method **POST**:

1. *"/subscribe"* - sign an email to receive information on changing the course.
You need to sent request with "x-www-form-urlencoded" payload:
    ```
    email: <email-for-subscribing>
    ```
2. *"/sendEmails"* - send the current course to all subscribed users.
