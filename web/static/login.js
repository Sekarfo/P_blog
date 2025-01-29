document.getElementById("loginForm").addEventListener("submit", async (event) => {
    event.preventDefault();
    const email = document.getElementById("email").value;
    const password = document.getElementById("password").value;

    try {
        const response = await fetch("/login", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ email, password }),
        });

        if (!response.ok) {
            throw new Error("Invalid email or password");
        }

        const { token, role } = await response.json();
        localStorage.setItem("token", token);
        localStorage.setItem("role", role); // Save user role
        window.location.href = "/posts.html";
    } catch (err) {
        document.getElementById("loginErrorMessage").textContent = err.message;
    }
});

document.getElementById("registerForm").addEventListener("submit", async (event) => {
    event.preventDefault();
    console.log("Register form submitted"); // log form submission

    const name = document.getElementById("registerName").value;
    const email = document.getElementById("registerEmail").value;
    const password = document.getElementById("registerPassword").value;

    try {
        const response = await fetch("/register", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ name, email, password }),
        });

        if (!response.ok) {
            throw new Error("Failed to register");
        }

        document.getElementById("registerErrorMessage").textContent = "Registration successful! Please check your email to verify your account.";
    } catch (err) {
        document.getElementById("registerErrorMessage").textContent = err.message;
    }
});