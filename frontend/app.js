async function login() {

    const email =
        document.getElementById("email").value;

    const password =
        document.getElementById("password").value;

    const response = await fetch(
        "http://localhost:8080/login",
        {
            method: "POST",

            headers: {
                "Content-Type": "application/json"
            },

            body: JSON.stringify({
                email: email,
                password: password
            })
        }
    );

    const data = await response.json();

    if (!response.ok) {
    document.getElementById("result").innerText =
        data.message;
    return;
}
    localStorage.setItem("token", data.token);

document.getElementById("result").innerText =
    data.message;

    console.log(data);
}