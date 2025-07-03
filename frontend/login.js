document.getElementById("loginForm").addEventListener("submit", function (e) {
  e.preventDefault();

  const username = document.getElementById("username").value.trim();
  const password = document.getElementById("password").value;

  fetch("/v1/login", {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify({ username, password })
  })
  .then(res => {
    if (!res.ok) throw new Error("Invalid credentials");
    return res.json();
  })
  .then(data => {
    localStorage.setItem("token", data.token);
    localStorage.setItem("username", username);
    alert("Login successful!");
    window.location.href = "index.html";
  })
  .catch(err => {
    alert("Login failed: " + err.message);
  });
});
