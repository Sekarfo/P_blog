document.addEventListener('DOMContentLoaded', () => {
    fetchAndDisplayArticles();
    /* document.getElementById('createUserForm').addEventListener('submit', createUser);
    document.getElementById('updateUserForm').addEventListener('submit', updateUser);
    document.getElementById('findUserForm').addEventListener('submit', findUserById);
    fetchUsers(); */
});

document.addEventListener('DOMContentLoaded', () => {
    fetchAndDisplayArticles();
});

function fetchAndDisplayArticles(page = 1) {
    const query = document.getElementById('query').value || 'top-headlines'; // Default query
    const sortBy = document.getElementById('sortBy').value || 'relevancy'; // Default sort by
    const pageSize = document.getElementById('pageSize').value || 5; // Default page size
    const language = document.getElementById('language').value || 'en'; // Default language
    const url = `/api/articles?q=${query}&sortBy=${sortBy}&pageSize=${pageSize}&page=${page}&language=${language}`;

    fetch(url)
        .then(response => response.json())
        .then(data => {
            const articlesDiv = document.getElementById('articles');
            if (!articlesDiv) {
                console.error('Element with ID "articles" not found.');
                return;
            }
            articlesDiv.innerHTML = ''; // Clear previous articles

            // Display total results
            const totalResultsElement = document.createElement('p');
            totalResultsElement.textContent = `Total Results: ${data.Total}`;
            articlesDiv.appendChild(totalResultsElement);

            data.articles.forEach(article => {
                const articleElement = document.createElement('div');
                articleElement.className = 'article';
                articleElement.innerHTML = `
                    <h3>${article.title}</h3>
                    <p>${article.description}</p>
                    <a href="${article.url}" target="_blank">Read more</a>
                    <p><strong>Author:</strong> ${article.author}</p>
                    <p><strong>Published At:</strong> ${new Date(article.publishedAt).toLocaleString()}</p>
                `;
                articlesDiv.appendChild(articleElement);
            });

            // Update pagination
            updatePagination(page, pageSize, data.Total);
        })
        .catch(error => console.error('Error fetching articles:', error));
}

function updatePagination(currentPage, pageSize, totalResults) {
    const paginationDiv = document.getElementById('pagination');
    paginationDiv.innerHTML = ''; // Clear previous pagination

    const totalPages = Math.ceil(totalResults / pageSize);
    const maxPages = 5; // Set the maximum number of pages to display
    const pagesToShow = Math.min(totalPages, maxPages);

    for (let i = 1; i <= pagesToShow; i++) {
        const button = document.createElement('button');
        button.textContent = i;
        button.className = (i === currentPage) ? 'active' : '';
        button.onclick = () => fetchAndDisplayArticles(i);
        paginationDiv.appendChild(button);
    }

    if (totalPages > maxPages) {
        const morePages = document.createElement('span');
        morePages.textContent = '...';
        paginationDiv.appendChild(morePages);
    }
}

/* let usersPerPage = 5;
let currentPage = 1;
let filteredUsers = [];

// Fetch all users and initialize pagination
function fetchUsers() {
    fetch('/api/users')
        .then(response => response.json())
        .then(users => {
            filteredUsers = users;
            displayUsers();
        })
        .catch(error => console.error('Error fetching users:', error));
}

// Display users with pagination
function displayUsers() {
    const startIndex = (currentPage - 1) * usersPerPage;
    const endIndex = startIndex + usersPerPage;
    const usersToShow = filteredUsers.slice(startIndex, endIndex);

    const tbody = document.querySelector('#usersTable tbody');
    tbody.innerHTML = '';

    usersToShow.forEach(user => {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${user.id}</td>
            <td>${user.name}</td>
            <td>${user.email}</td>
            <td><button onclick="deleteUser(${user.id})">Delete</button></td>
        `;
        tbody.appendChild(row);
    });

    setupPagination();
}

// Setup pagination buttons
function setupPagination() {
    const totalPages = Math.ceil(filteredUsers.length / usersPerPage);
    const paginationDiv = document.getElementById('pagination');
    paginationDiv.innerHTML = '';

    for (let i = 1; i <= totalPages; i++) {
        const button = document.createElement('button');
        button.textContent = i;
        button.className = (i === currentPage) ? 'active' : '';
        button.onclick = () => {
            currentPage = i;
            displayUsers();
        };
        paginationDiv.appendChild(button);
    }
}

// Filter users by name
function filterUsers() {
    const filterValue = document.getElementById('filterName').value.toLowerCase();
    fetch('/api/users')
        .then(response => response.json())
        .then(users => {
            filteredUsers = users.filter(user => user.name.toLowerCase().includes(filterValue));
            currentPage = 1;
            displayUsers();
        })
        .catch(error => console.error('Error filtering users:', error));
}

// Remaining functions: createUser, updateUser, findUserById, deleteUser
// (These functions stay as you provided)
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

 */