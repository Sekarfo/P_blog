async function fetchUsers(page = 1, search = "") {
    const response = await fetch(`/admin/users?page=${page}&search=${search}`, {
        headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
    });

    if (!response.ok) {
        const errorText = await response.text();
        alert(`Error: ${errorText}`);
        return;
    }

    const users = await response.json();

    const userTable = document.querySelector("#userTable tbody");
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

    document.getElementById("currentPage").textContent = users.current_page;
    document.getElementById("prevPage").disabled = !users.prev_page;
    document.getElementById("nextPage").disabled = !users.next_page;

    currentPage = users.current_page;
}

async function changeRole(userID, roleID) {
    const roles = { Admin: 1, Writer: 2, Reader: 3 }; // Ensure role mapping matches your database
    try {
        const response = await fetch("/admin/users/role", {
            method: "PATCH",
            headers: {
                "Content-Type": "application/json",
                Authorization: `Bearer ${localStorage.getItem("token")}`,
            },
            body: JSON.stringify({ user_id: userID, role_id: roles[roleID] }),
        });

        if (!response.ok) {
            throw new Error("Failed to update role");
        }

        alert("Role updated successfully");
    } catch (err) {
        console.error(err);
        alert(err.message);
    }
}


async function deleteUser(userID) {
    const response = await fetch(`/admin/users/${userID}`, {
        method: "DELETE",
        headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
    });

    if (!response.ok) {
        const errorText = await response.text();
        alert(`Error: ${errorText}`);
        return;
    }

    fetchUsers(currentPage);
}

function filterUsers() {
    const searchValue = document.getElementById("searchUser").value;
    fetchUsers(1, searchValue);
}

function requireAuth() {
    const token = localStorage.getItem("token");
    if (!token) {
        alert("You must be logged in to access this page.");
        window.location.href = "/login.html";
    }
}

document.addEventListener("DOMContentLoaded", requireAuth);

document.getElementById("prevPage").addEventListener("click", () => {
    if (currentPage > 1) fetchUsers(currentPage - 1);
});
document.getElementById("nextPage").addEventListener("click", () => {
    fetchUsers(currentPage + 1);
});

document.addEventListener("DOMContentLoaded", () => fetchUsers());