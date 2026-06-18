function setResult(message, type) {
  const el = document.getElementById('result');
  if (!el) return;
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

    const token = data.token;

    function parseRoleFromToken(t) {
      if (!t) return 'customer';
      try {
        const parts = String(t).split('.');
        if (parts.length < 2) return 'customer';
        const base64Url = parts[1];
        const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
        const payload = JSON.parse(atob(base64));

        // Support possible claim key variations
        return payload.role ?? payload.Role ?? payload.user_role ?? payload.userRole ?? 'customer';
      } catch (_) {
        return 'customer';
      }
    }

    const userRole = parseRoleFromToken(token);

    if (userRole === 'admin') {
      window.location.href = 'admin_products.html';
    } else if (userRole === 'seller') {
      window.location.href = 'seller_products.html';
    } else {
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

    // Success state: Update UI and redirect
    setResult('Account created! Redirecting to login...', 'success');
    
    // Redirect to login page after a 1.5s delay
    setTimeout(() => {
      window.location.href = 'login.html';
    }, 1000);

  } catch (err) {
    setResult('Network error. Is the backend running on http://localhost:8080?', 'error');
  } finally {
    setSubmitting(false);
  }
}