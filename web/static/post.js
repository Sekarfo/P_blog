const postID = new URLSearchParams(window.location.search).get("id");

async function fetchPostDetails() {
    const token = localStorage.getItem("token");

    if (!token) {
        alert("You need to log in first!");
        window.location.href = "/login.html";
        return;
    }

    const postID = new URLSearchParams(window.location.search).get("id");

    try {
        const response = await fetch(`/posts/${postID}`, {
            method: "GET",
            headers: {
                Authorization: `Bearer ${token}`,
            },
        });

        if (!response.ok) {
            throw new Error("Failed to fetch post details.");
        }

        const post = await response.json();
        console.log(post); // Debug: log post details

        // Render post details
        document.getElementById("postTitle").textContent = post.title;
        document.getElementById("postContent").textContent = post.content;
        document.getElementById("postAuthor").textContent = post.author;
        document.getElementById("postLikes").textContent = post.likes;

        // Check if user has liked the post
        const likeResponse = await fetch(`/posts/${postID}/likes`, {
            method: "GET",
            headers: {
                Authorization: `Bearer ${token}`,
            },
        });

        if (likeResponse.ok) {
            const { liked } = await likeResponse.json();
            if (liked) {
                document.getElementById("likeButton").style.display = "none";
                document.getElementById("unlikeButton").style.display = "inline-block";
            } else {
                document.getElementById("likeButton").style.display = "inline-block";
                document.getElementById("unlikeButton").style.display = "none";
            }
        }
    } catch (err) {
        console.error(err);
        alert("Error: " + err.message);
    }
}

async function fetchComments() {
    const token = localStorage.getItem("token");

    if (!token) {
        alert("You need to log in first!");
        window.location.href = "/login.html";
        return;
    }

    try {
        const response = await fetch(`/posts/${postID}/comments`, {
            method: "GET",
            headers: {
                Authorization: `Bearer ${token}`,
            },
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(`Failed to fetch comments: ${errorText}`);
        }

        const comments = await response.json();
        console.log(comments); // Debug: log comments

        const commentsList = document.getElementById("commentsList");
        commentsList.innerHTML = ""; // Clear the current list

        // Check if comments is not null or undefined
        if (comments && comments.length > 0) {
            comments.forEach((comment) => {
                const li = document.createElement("li");
                li.innerHTML = `
                    <p>${comment.content}</p>
                    <small>By: ${comment.author}, on ${new Date(comment.created).toLocaleString()}</small>
                `;
                commentsList.appendChild(li);
            });
        } else {
            commentsList.innerHTML = "<p>No comments found for this post.</p>";
        }
    } catch (err) {
        console.error(err);
        alert("Error: " + err.message);
    }
}

async function addComment(event) {
    event.preventDefault();

    const token = localStorage.getItem("token");

    if (!token) {
        alert("You need to log in first!");
        window.location.href = "/login.html";
        return;
    }

    const content = document.getElementById("commentContent").value;
    //const postID = new URLSearchParams(window.location.search).get("id");
    const postID = Number(new URLSearchParams(window.location.search).get("id"));


    // Log the request payload
    console.log("Adding comment with content:", content, "for post ID:", postID);

    try {
        const response = await fetch(`/comments`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: `Bearer ${token}`,
            },
            body: JSON.stringify({ content, post_id: Number(postID) }),
        });
        console.log("Server error response:", errorText);

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(`Failed to add comment: ${errorText}`);
        }

        document.getElementById("commentContent").value = ""; // Clear input
        fetchComments(); // Refresh comments
    } catch (err) {
        console.error(err);
        alert("Error: " + err.message);
    }
}

async function likePost() {
    try {
        const response = await fetch(`/posts/${postID}/like`, {
            method: "POST",
            headers: {
                Authorization: `Bearer ${localStorage.getItem("token")}`,
            },
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(`Failed to like post: ${errorText}`);
        }

        const data = await response.json();
        document.getElementById("postLikes").textContent = data.likes;
    } catch (err) {
        console.error(err);
        alert("Error: " + err.message);
    }
}

async function unlikePost() {
    try {
        const response = await fetch(`/posts/${postID}/unlike`, {
            method: "POST",
            headers: {
                Authorization: `Bearer ${localStorage.getItem("token")}`,
            },
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(`Failed to unlike post: ${errorText}`);
        }

        const data = await response.json();
        document.getElementById("postLikes").textContent = data.likes;
    } catch (err) {
        console.error(err);
        alert("Error: " + err.message);
    }
}

document.getElementById("likeButton").addEventListener("click", likePost);
document.getElementById("unlikeButton").addEventListener("click", unlikePost);
document.getElementById("commentForm").addEventListener("submit", addComment);

document.addEventListener("DOMContentLoaded", () => {
    fetchPostDetails();
    fetchComments();
});