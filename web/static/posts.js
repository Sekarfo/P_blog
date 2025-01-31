async function fetchPosts() {
    const token = localStorage.getItem("token");

    if (!token) {
        alert("You need to log in first!");
        window.location.href = "/login.html";
        return;
    }

    console.log("Token:", token); // Debug: log the token value

    try {
        const response = await fetch("/posts", {
            method: "GET",
            headers: {
                Authorization: `Bearer ${token}`, // Include the Authorization header
            },
        });

        if (!response.ok) {
            throw new Error("Failed to fetch posts. Please check your authorization.");
        }

        const posts = await response.json(); // Parse the JSON response
        console.log(posts); // Debug: log the fetched posts

        const postList = document.getElementById("postList");
        postList.innerHTML = ""; // Clear the existing list

        posts.forEach((post) => {
            const li = document.createElement("li");
            li.innerHTML = `
                <h3><a href="/post.html?id=${post.id}">${post.title}</a></h3>
                <p>${post.content.substring(0, 100)}...</p>
                <small>Author: ${post.author}</small>
            `;
            postList.appendChild(li); // Add the post to the list
        });
    } catch (err) {
        console.error(err); // Log
        alert("Error: " + err.message);
    }
}

async function createPost(event) {
    event.preventDefault();

    const title = document.getElementById("postTitle").value;
    const content = document.getElementById("postContent").value;

    const response = await fetch("/posts", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
        body: JSON.stringify({ title, content }),
    });

    if (!response.ok) {
        const errorText = await response.text();
        alert(`Error: ${errorText}`);
        return;
    }

    fetchPosts();
}

async function deletePost(postID) {
    const response = await fetch(`/posts/${postID}`, {
        method: "DELETE",
        headers: {
            Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
    });

    if (!response.ok) {
        const errorText = await response.text();
        alert(`Error: ${errorText}`);
        return;
    }

    fetchPosts();
}

async function fetchNews() {
    try {
        const response = await fetch("https://newsapi.org/v2/top-headlines?country=us&apiKey=ec79b26e36ca4f3ba4f320ceec94a351");

        if (!response.ok) {
            throw new Error("Failed to fetch news articles.");
        }

        const newsData = await response.json();
        console.log(newsData); // Debug: log the fetched news articles

        const newsList = document.getElementById("newsList");
        newsList.innerHTML = ""; // Clear the existing list

        newsData.articles.forEach((article) => {
            const li = document.createElement("li");
            li.innerHTML = `
                <h3><a href="${article.url}" target="_blank">${article.title}</a></h3>
                <p>${article.description}</p>
                <small>Source: ${article.source.name} | Published at: ${new Date(article.publishedAt).toLocaleString()}</small>
            `;
            newsList.appendChild(li); // Add the news article to the list
        });
    } catch (err) {
        console.error(err); // Log the error
        alert("Error: " + err.message); // Display an alert if something goes wrong
    }
}

function requireAuth() {
    const token = localStorage.getItem("token");
    console.log("Token:", token); // Debug token value
    if (!token) {
        alert("You must be logged in to access this page.");
        window.location.href = "/login.html";
    }
}

document.addEventListener("DOMContentLoaded", () => {
    requireAuth();

    const role = localStorage.getItem("role");
    if (role !== "Admin") {
        const adminLink = document.querySelector('a[href="/admin.html"]');
        if (adminLink) adminLink.style.display = "none";
    }

    fetchPosts();
    fetchNews();
});

document.getElementById("createPostForm").addEventListener("submit", createPost);
document.addEventListener("DOMContentLoaded", fetchPosts);