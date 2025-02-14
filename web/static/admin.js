async function fetchSubscriptions() {
    try {
        const response = await fetch("http://localhost:8081/api/subscriptions/admin/subscriptions", {
            headers: {
                Authorization: `Bearer ${localStorage.getItem("token")}`,
                "Cache-Control": "no-cache, no-store, must-revalidate"
            }
        });

        if (!response.ok) {
            console.error("Failed to fetch subscriptions:", response.status, await response.text());
            return;
        }

        const subscriptions = await response.json();
        console.log("Subscriptions received:", subscriptions);

        const table = document.getElementById("subscriptionTable");
        if (!table) {
            console.error("Subscription table not found.");
            return;
        }

        table.innerHTML = `<tr>
            <th>User</th>
            <th>Status</th>
            <th>Requested At</th>
            <th>Start Date</th>
            <th>End Date</th>
            <th>Actions</th>
        </tr>`;

        if (subscriptions.length === 0) {
            table.innerHTML += "<tr><td colspan='6'>No subscriptions found.</td></tr>";
            return;
        }

        subscriptions.forEach(sub => {
            if (!sub.id) {
                console.error("Subscription ID is undefined:", sub);
                return;
            }

            const row = document.createElement("tr");
            row.innerHTML = `
                <td>${sub.user_id} - ${sub.user_name}</td>
                <td>${sub.status}</td>
                <td>${new Date(sub.requested_at).toLocaleString()}</td>
                <td>${sub.subscription_start ? new Date(sub.subscription_start).toLocaleDateString() : "N/A"}</td>
                <td>${sub.subscription_end ? new Date(sub.subscription_end).toLocaleDateString() : "N/A"}</td>
                <td>
                    ${sub.status === "pending_approval" ? `
                        <button class="approve-btn" data-id="${sub.id}">Approve</button>
                        <button class="reject-btn" data-id="${sub.id}">Reject</button>
                    ` : "N/A"}
                </td>
            `;
            table.appendChild(row);
        });

        document.querySelectorAll(".approve-btn").forEach(button => {
            button.addEventListener("click", function () {
                const id = this.getAttribute("data-id");
                approveSubscription(id);
            });
        });

        document.querySelectorAll(".reject-btn").forEach(button => {
            button.addEventListener("click", function () {
                const id = this.getAttribute("data-id");
                rejectSubscription(id);
            });
        });

        console.log("Subscription table updated successfully.");
    } catch (error) {
        console.error("Error fetching subscriptions:", error);
    }
}



async function approveSubscription(id) {
    try {
        console.log(`Approving subscription ID: ${id}`); // ✅ Debugging log

        const response = await fetch(`http://localhost:8081/api/subscriptions/admin/approve/${id}`, {
            method: "PATCH",
            headers: { 
                "Content-Type": "application/json",
                Authorization: `Bearer ${localStorage.getItem("token")}`
            }
        });

        const result = await response.text(); // Get full response text
        console.log("Approve response:", result); // ✅ Debugging log

        if (!response.ok) {
            throw new Error("Failed to approve subscription.");
        }

        alert("✅ Subscription approved! Email with receipt will be sent.");
        fetchSubscriptions(); // Refresh subscription list
    } catch (error) {
        console.error("❌ Error approving subscription:", error);
        alert("❌ Error approving subscription.");
    }
}


async function rejectSubscription(id) {
    await fetch(`http://localhost:8081/api/subscriptions/admin/reject/${id}`, {
        method: "PATCH",
        headers: { Authorization: `Bearer ${localStorage.getItem("token")}` }
    });
    alert("Subscription rejected!");
    fetchSubscriptions();
}

async function requireAdminAuth() {
    const token = localStorage.getItem("token");
    if (!token) {
        alert("You must be logged in as an admin.");
        window.location.href = "/login.html";
        return;
    }

    try {
        const response = await fetch("http://localhost:8080/api/profile", {
            headers: { Authorization: `Bearer ${token}` }
        });

        if (!response.ok) {
            const errorText = await response.text();
            console.error("Failed to verify authentication:", errorText);
            alert("Failed to verify authentication.");
            window.location.href = "/login.html";
            return;
        }

        const user = await response.json();
        console.log("User:", user);

        if (user.role !== "Admin") {
            alert("Access denied. Only admins can view this page.");
            window.location.href = "/posts.html";
        }
    } catch (err) {
        console.error("Error verifying authentication:", err);
        alert("Failed to verify authentication.");
        window.location.href = "/login.html";
    }
}

async function fetchUsers(page = 1, search = "") {
    console.log("Fetching users...");

    const response = await fetch(`/admin/users?page=${page}&search=${search}`, {
        headers: { Authorization: `Bearer ${localStorage.getItem("token")}` }
    });

    if (!response.ok) {
        console.error("Failed to fetch users:", response.status);
        return;
    }

    const users = await response.json();
    console.log("Users received:", users);

    const userTable = document.querySelector("#userTable tbody");

    if (!userTable) {
        console.error("User table tbody not found.");
        return;
    }

    userTable.innerHTML = "";

    users.data.forEach((user) => {
        const row = document.createElement("tr");
        row.innerHTML = `
            <td>${user.id}</td>
            <td>${user.name}</td>
            <td>${user.email}</td>
            <td>
                <select onchange="changeRole(${user.id}, this.value)">
                    <option value="Admin" ${user.role === "Admin" ? "selected" : ""}>Admin</option>
                    <option value="Writer" ${user.role === "Writer" ? "selected" : ""}>Writer</option>
                    <option value="Reader" ${user.role === "Reader" ? "selected" : ""}>Reader</option>
                </select>
            </td>
            <td>
                <button onclick="deleteUser(${user.id})">Delete</button>
            </td>
        `;
        userTable.appendChild(row);
    });
}

// Add event listener to the search button
document.getElementById("searchUserBtn").addEventListener("click", () => {
    const searchQuery = document.getElementById("searchUser").value;
    fetchUsers(1, searchQuery);
});

// Ensure both sections load data when the page loads
document.addEventListener("DOMContentLoaded", () => {
    requireAdminAuth();
    fetchUsers();
    fetchSubscriptions();
});