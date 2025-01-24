document.addEventListener('DOMContentLoaded', () => {
    const loginForm = document.getElementById('loginForm');
    if (loginForm) {
        loginForm.addEventListener('submit', async (event) => {
            event.preventDefault(); // Prevent default form submission
            const email = document.getElementById('email').value;
            const password = document.getElementById('password').value;

            try {
                const response = await fetch('/login', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/x-www-form-urlencoded',
                    },
                    body: new URLSearchParams({
                        email: email,
                        password: password,
                    }),
                    credentials: 'include', // Include cookies
                });

                if (response.ok) {
                    const data = await response.json();
                    if (data.status === 'success') {
                        window.location.href = '/profile'; // Redirect on success
                    } else {
                        alert(data.message || 'Invalid credentials');
                    }
                } else {
                    alert('Invalid email/password or Email not Verified');
                }
            } catch (error) {
                console.error('Error:', error);
                alert('An error occurred during login');
            }
        });
    }

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