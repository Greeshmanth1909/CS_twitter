function login() {
  alert("Login clicked");
}

function signup() {
  alert("Signup clicked");
}

function submitPost() {
  const content = document.getElementById("postContent").value;
  if (!content.trim()) {
    alert("Post cannot be empty.");
    return;
  }
  jwt = localStorage.getItem("token");
  if (jwt == null) {
    alert("Please Login before creating a post");
    return;
  }

  fetch("/v1/create-post", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "Authorization": `Bearer ${jwt}`
    },
    body: JSON.stringify({ post: content })
  })
  .then(res => {
    if (!res.ok) {
        console.log(res);
        alert("Please Login!");
        return;
    }
    alert("Post submitted!");
    document.getElementById("postContent").value = "";
    loadPosts();
  })
  .catch(err => alert("Error: " + err.message));
}

function submitComment(postId) {
  const input = document.getElementById(`comment-input-${postId}`);
  const comment = input.value;
  console.log("comment", comment)
  if (!comment) return alert("Comment cannot be empty.");
  jwt = localStorage.getItem("token");
  if (jwt == null) {
    alert("Please Login before creating a post");
    return;
  }

  fetch("/v1/create-comment", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "Authorization": `Bearer ${jwt}`

    },
    body: JSON.stringify({ comment: comment, post_id: postId })
  })
  .then(res => {
    console.log(res);
    if (!res.ok) throw new Error("Failed to comment");
    input.value = "";
    loadPosts();
  })
  .catch(err => alert("Error: " + err.message));
}

function loadPosts() {
  fetch("/v1/posts")
    .then(res => res.json())
    .then(data => {
      const feed = document.getElementById("feed");
      feed.innerHTML = "";

      data.forEach((post) => {
        const postDiv = document.createElement("div");
        postDiv.className = "post";

        const user = document.createElement("h3");
        user.textContent = post.Username;

        const content = document.createElement("p");
        content.textContent = post.Post;

        const commentsDiv = document.createElement("div");
        post.Comments.forEach(([text, commenter]) => {
            const c = document.createElement("div");
            c.className = "comment";

            var formattedText = text.replace(/\*\*(.*?)\*\*/g, "<b>$1</b>");
            formattedText = formattedText.replace(/\*(.*?)\*/g, "<em>$1</em>");

            // for hyperlinks [text](url)
            formattedText = formattedText.replace(/\[([^\]]+)\]\((https?:\/\/[^\s)]+)\)/g, "<a href=$2>$1</a>");

            c.innerHTML = `<strong>${commenter}:</strong> ${formattedText}`;
            commentsDiv.appendChild(c);
        });


        const commentInput = document.createElement("input");
        commentInput.type = "text";
        commentInput.placeholder = "Write a comment...";
        commentInput.id = `comment-input-${post.Post_id}`;

        const commentButton = document.createElement("button");
        commentButton.textContent = "Post Comment";
        commentButton.onclick = () => submitComment(post.Post_id);

        const commentForm = document.createElement("div");
        commentForm.className = "comment-form";
        commentForm.appendChild(commentInput);
        commentForm.appendChild(commentButton);

        postDiv.appendChild(user);
        postDiv.appendChild(content);
        postDiv.appendChild(commentsDiv);
        postDiv.appendChild(commentForm);

        feed.appendChild(postDiv);
      });
    })
    .catch(err => {
      document.getElementById("feed").textContent = "Failed to load posts.";
      console.error(err);
    });
}

window.onload = loadPosts;
