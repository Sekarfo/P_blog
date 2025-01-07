// static/script.js

document.addEventListener('DOMContentLoaded', () => {
    // Initialize forms
    document.getElementById('createUserForm').addEventListener('submit', createUser);
    document.getElementById('updateUserForm').addEventListener('submit', updateUser);
    document.getElementById('findUserForm').addEventListener('submit', findUserById);
});

// Function to create a new user
function createUser(event) {
    event.preventDefault(); // Prevent form from submitting the traditional way

    const name = document.getElementById('createName').value.trim();
    const email = document.getElementById('createEmail').value.trim();

    fetch('/api/users', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ name, email })
    })
        .then(response => {
            if (!response.ok) {
                return response.text().then(text => { throw new Error(text) });
            }
            return response.json();
        })
        .then(data => {
            alert(`User created with ID: ${data.id}`);
            document.getElementById('createUserForm').reset();
            fetchUsers(); // Refresh the users table
        })
        .catch(error => {
            alert(`Error: ${error.message}`);
        });
}


function updateUser(event) {
    event.preventDefault();

    const id = document.getElementById('updateId').value.trim();
    const name = document.getElementById('updateName').value.trim();
    const email = document.getElementById('updateEmail').value.trim();

    if (!id) {
        alert('User ID is required for updating.');
        return;
    }

    const updateData = {};
    if (name) updateData.name = name;
    if (email) updateData.email = email;

    if (Object.keys(updateData).length === 0) {
        alert('Please provide at least one field to update.');
        return;
    }

    fetch(`/api/users/${id}`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(updateData)
    })
        .then(response => {
            if (!response.ok) {
                return response.text().then(text => { throw new Error(text) });
            }
            return response.json();
        })
        .then(data => {
            alert(`User with ID: ${data.id} updated successfully.`);
            document.getElementById('updateUserForm').reset();
            fetchUsers(); // Refresh the users table
        })
        .catch(error => {
            alert(`Error: ${error.message}`);
        });
}


function findUserById(event) {
    event.preventDefault();

    const id = document.getElementById('findId').value.trim();

    if (!id) {
        alert('Please enter a User ID.');
        return;
    }

    fetch(`/api/users/${id}`)
        .then(response => {
            if (!response.ok) {
                return response.text().then(text => { throw new Error(text) });
            }
            return response.json();
        })
        .then(user => {
            displayFindUserResult(user);
            document.getElementById('findUserForm').reset();
        })
        .catch(error => {
            displayFindUserResult(null, error.message);
        });
}

// Function to display the found user in the findUserResult div
function displayFindUserResult(user, error = null) {
    const resultDiv = document.getElementById('findUserResult');
    resultDiv.innerHTML = ''; // Clear previous results

    if (error) {
        resultDiv.innerHTML = `<p style="color: red;">Error: ${error}</p>`;
        return;
    }

    if (user) {
        const table = document.createElement('table');
        table.innerHTML = `
            <thead>
                <tr>
                    <th>ID</th>
                    <th>Name</th>
                    <th>Email</th>
                </tr>
            </thead>
            <tbody>
                <tr>
                    <td>${user.id}</td>
                    <td>${user.name}</td>
                    <td>${user.email}</td>
                </tr>
            </tbody>
        `;
        resultDiv.appendChild(table);
    } else {
        resultDiv.innerHTML = '<p>User not found.</p>';
    }
}

// Function to fetch and display all users
function fetchUsers() {
    fetch('/api/users')
        .then(response => response.json())
        .then(data => {
            const tbody = document.querySelector('#usersTable tbody');
            tbody.innerHTML = ''; // Clear existing rows

            data.forEach(user => {
                const row = document.createElement('tr');

                // ID
                const idCell = document.createElement('td');
                idCell.textContent = user.id;
                row.appendChild(idCell);

                // Name
                const nameCell = document.createElement('td');
                nameCell.textContent = user.name;
                row.appendChild(nameCell);

                // Email
                const emailCell = document.createElement('td');
                emailCell.textContent = user.email;
                row.appendChild(emailCell);

                // Actions (Delete Button)
                const actionsCell = document.createElement('td');
                const deleteButton = document.createElement('button');
                deleteButton.textContent = 'Delete';
                deleteButton.classList.add('actions-button');
                deleteButton.onclick = () => deleteUser(user.id);
                actionsCell.appendChild(deleteButton);
                row.appendChild(actionsCell);

                tbody.appendChild(row);
            });
        })
        .catch(error => {
            console.error('Error fetching users:', error);
            alert('Failed to fetch users.');
        });
}

// Function to delete a user
function deleteUser(id) {
    if (!confirm(`Are you sure you want to delete user with ID: ${id}?`)) {
        return;
    }

    fetch(`/api/users/${id}`, {
        method: 'DELETE'
    })
        .then(response => {
            if (response.status === 204) {
                alert(`User with ID: ${id} deleted successfully.`);
                fetchUsers(); // Refresh the users table
            } else {
                return response.text().then(text => { throw new Error(text) });
            }
        })
        .catch(error => {
            alert(`Error: ${error.message}`);
        });
}
