function setResult(message, type) {
  const el = document.getElementById('result');
  el.classList.remove('error', 'success');
  if (type) el.classList.add(type);
  el.innerText = message || '';
}

function setSubmitting(isSubmitting) {
  const btn = document.getElementById('submitBtn');
  if (!btn) return;
  btn.disabled = isSubmitting;
}

async function login() {
  try {
    setSubmitting(true);
    setResult('', null);

    const email = document.getElementById('email').value.trim();
    const password = document.getElementById('password').value;

    const response = await fetch('http://localhost:8080/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password }),
    });

    const data = await response.json().catch(() => ({}));

    if (!response.ok) {
      setResult(data.message || 'Login failed', 'error');
      return;
    }

    localStorage.setItem('token', data.token);
    setResult(data.message || 'Login successful', 'success');

    // Redirect to products page
    const postLoginLink = document.getElementById('postLoginLink');
    if (postLoginLink) {
      postLoginLink.style.display = 'inline';
    } else {
      // fallback navigation
      window.location.href = 'products.html';
    }

  } catch (err) {
    setResult('Network error. Is the backend running on http://localhost:8080?', 'error');
  } finally {
    setSubmitting(false);
  }
}

async function signup() {
  try {
    setSubmitting(true);
    setResult('', null);

    const name = document.getElementById('name').value.trim();
    const email = document.getElementById('email').value.trim();
    const password = document.getElementById('password').value;
    const role = document.getElementById('role').value;

    const response = await fetch('http://localhost:8080/signup', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name, email, password, role }),
    });

    const data = await response.json().catch(() => ({}));

    if (!response.ok) {
      setResult(data.message || 'Signup failed', 'error');
      return;
    }

    setResult(data.message || 'User created successfully', 'success');
  } catch (err) {
    setResult('Network error. Is the backend running on http://localhost:8080?', 'error');
  } finally {
    setSubmitting(false);
  }
}

