document.getElementById("signupForm").addEventListener("submit", function (e) {
  e.preventDefault();

  const username = document.getElementById("username").value.trim();
  const password = document.getElementById("password").value;

  fetch("/v1/signup", {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify({ username, password })
  })
  .then(res => {
    if (res.status != 201) {
        alert("username taken!");
        return;
    } else {
        alert("Signup successful!");
        window.location.href = "index.html";
    }
    
  })
  .catch(err => alert("Error: " + err.message));
});
