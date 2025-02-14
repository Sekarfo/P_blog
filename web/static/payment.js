document.addEventListener("DOMContentLoaded", function () {
    const paymentForm = document.getElementById("paymentForm");
    const message = document.getElementById("message");

    // Get subscription ID from query params
    const urlParams = new URLSearchParams(window.location.search);
    const subscriptionId = urlParams.get("subscriptionId");
    if (!subscriptionId) {
        alert("Invalid subscription request. Please try again.");
        window.location.href = "/subscribe.html"; // Redirect back if no ID
    }
    document.getElementById("subscriptionId").value = subscriptionId;

    paymentForm.addEventListener("submit", async function (event) {
        event.preventDefault();

        const paymentData = {
            cardNumber: document.getElementById("cardNumber").value,
            expirationDate: document.getElementById("expirationDate").value,
            cvv: document.getElementById("cvv").value,
            name: document.getElementById("name").value,
            address: document.getElementById("address").value,
            subscription_id: Number(subscriptionId) // ✅ Convert to number
        };
              

        try {
            const response = await fetch("http://localhost:8081/api/subscriptions/payment", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: `Bearer ${localStorage.getItem("token")}`
                },
                body: JSON.stringify(paymentData)
            });

            const result = await response.json();
            if (response.ok) {
                message.textContent = result.payment_success ? "Payment Successful!" : "Payment Declined.";
                message.style.color = result.payment_success ? "green" : "red";

                if (result.payment_success) {
                    setTimeout(() => {
                        window.location.href = "/posts.html"; // Redirect to profile after success
                    }, 2000);
                }
            } else {
                message.textContent = result.error || "Payment failed.";
                message.style.color = "red";
            }
        } catch (error) {
            console.error("Error processing payment:", error);
            message.textContent = "An error occurred.";
            message.style.color = "red";
        }
    });
});
