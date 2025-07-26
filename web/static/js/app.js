// ======= JWT парсинг из localStorage =======
function parseJwt(token) {
  try {
    const payload = token.split(".")[1];
    const json = atob(payload.replace(/-/g, "+").replace(/_/g, "/"));
    return JSON.parse(
      decodeURIComponent(
        json
          .split("")
          .map((c) => "%" + ("00" + c.charCodeAt(0).toString(16)).slice(-2))
          .join(""),
      ),
    );
  } catch (e) {
    console.error("[parseJwt] Invalid token", e);
    return null;
  }
}

// ======= Отображение кнопок в шапке =======
function setupAuthButtons() {
  const authContainer = document.querySelector(".header-auth");
  if (!authContainer) return;

  const token = localStorage.getItem("jwt");
  if (token) {
    const claims = parseJwt(token);
    if (claims?.role === "admin") {
      authContainer.innerHTML = `
        <a href="/stats" class="header-auth_link--l">Статистика</a>
        <a href="#" id="logout-btn" class="header-auth_link--r">Выйти</a>
      `;
    } else {
      authContainer.innerHTML = `
        <a href="#" id="logout-btn" class="header-auth_link--r">Выйти</a>
      `;
    }
    document.getElementById("logout-btn").onclick = (e) => {
      e.preventDefault();
      localStorage.removeItem("jwt");
      window.location.href = "/signin";
    };
  } else {
    authContainer.innerHTML = `
      <a href="/signin" class="header-auth_link--l">Войти</a>
      <a href="/signup" class="header-auth_link--r">Регистрация</a>
    `;
  }
}

// ======= Обработка формы входа =======
function initSigninForm() {
  const form = document.getElementById("signin-form");
  if (!form) return;

  form.addEventListener("submit", async (e) => {
    e.preventDefault();
    const data = {
      email: form.email.value,
      password: form.password.value,
      role: form.role?.value || "user",
    };
    try {
      const res = await fetch("/auth/signin", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      });
      const json = await res.json();
      if (!res.ok) {
        alert(json.error || "Ошибка авторизации");
        return;
      }
      localStorage.setItem("jwt", json.token);
      window.location.href = "/";
    } catch (err) {
      console.error(err);
      alert("Сетевая ошибка");
    }
  });
}

// ======= Обработка формы регистрации =======
function initSignupForm() {
  const form = document.getElementById("signup-form");
  if (!form) return;

  form.addEventListener("submit", async (e) => {
    e.preventDefault();
    const data = {
      name: form.name.value,
      email: form.email.value,
      password: form.password.value,
      role: form.role.value,
      is_blocked: false,
    };
    try {
      const res = await fetch("/auth/signup", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      });
      const json = await res.json();
      if (!res.ok) {
        alert(json.error || "Ошибка регистрации");
        return;
      }
      localStorage.setItem("jwt", json.token);
      window.location.href = "/";
    } catch (err) {
      console.error(err);
      alert("Сетевая ошибка");
    }
  });
}

// ======= Обработка формы сокращения ссылки =======
function initShortenForm() {
  const form = document.getElementById("shorten-form");
  const result = document.getElementById("short-url-result");
  if (!form || !result) return;

  form.addEventListener("submit", async (e) => {
    e.preventDefault();
    const urlValue = form.url.value.trim();
    if (!urlValue) return;

    result.textContent = "Сокращаем…";

    try {
      const token = localStorage.getItem("jwt");
      const res = await fetch("/users/links", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ url: urlValue }),
      });
      const json = await res.json();
      if (!res.ok) {
        result.textContent = json.error || "Ошибка";
        return;
      }

      const hash = json.hash;
      const shortLink = `/${hash}`;

      // Очищаем контейнер
      result.textContent = "Короткая ссылка: ";

      // Создаём <a> через DOM API
      const a = document.createElement("a");
      a.href = shortLink;
      a.target = "_blank";
      a.classList.add("shorten-result_link");
      a.textContent = shortLink;

      result.appendChild(a);

      form.url.value = "";
    } catch (err) {
      console.error(err);
      result.textContent = "Сетевая ошибка";
    }
  });
}

// ======= Запуск после загрузки =======
document.addEventListener("DOMContentLoaded", () => {
  setupAuthButtons();
  initSigninForm();
  initSignupForm();
  initShortenForm();
});
