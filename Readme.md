# How to run

1. Setup environment

- Register a test account here: https://sandbox.vnpayment.vn/devreg

- Wait for VNPay send the configure information to registered email

- Copy and paste `vnp_TmnCode`, `vnp_HashSecret` to `.env` file.

1. Install dependencies

```
go mod tidy
go mod vendor
```

2. Run the application

Again from the server directory run:

```sh
go run ./main/main.go
```

3. Go to `localhost:10000/pay` to see the demo.

# One time payment flow

```mermaid
sequenceDiagram
    title Create payment request

    actor User
    participant Frontend
    participant Backend
    participant VNPay
    autonumber

    User ->> Frontend: Create order

    activate User
    activate Frontend
    Frontend ->> Backend: Request a payment url

    activate Backend
    Backend ->> Frontend: Return the payment url
    deactivate Backend

    Frontend ->> VNPay: Open payment url <br/> and browser move to the VNPay payment method page
    deactivate Frontend

    activate VNPay
    User ->> VNPay: Select payment method (Apple pay, Card pay, QR pay)

    deactivate User
    VNPay -->> Frontend: Forward user to the result page
    deactivate VNPay
```

````mermaid
sequenceDiagram
    title Update payment result

    actor Dev
    participant VNPay
    participant Backend
    autonumber

    Dev ->> VNPay: Setup IPN url
    activate VNPay
    VNPay ->> Backend: Send payment result on each new payment
    deactivate VNPay```
````
