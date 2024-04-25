const loginForm = document.getElementById("login-form");
loginForm.addEventListener("submit", (e) => {
    e.preventDefault();
    const email = e.target.elements.email.value;
    const password = e.target.elements.password.value;
    
    const data = {
        email: email,
        password: password
    };

    fetch('http://localhost:8081/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        return response.json();
    })
    .then(data => {
        console.log('Success:', data);
        console.log(data)
        document.cookie = `token=${data.token}; path=/; SameSite=None; Secure`;
        alert('Successful login');
        window.location.href = "culc.html";
    })
    .catch(error => {
        console.error('Error:', error);
        alert('Login error');
    });
});

