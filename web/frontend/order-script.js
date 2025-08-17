document.addEventListener("DOMContentLoaded", function () {
  const input = document.getElementById("uuidInput");
  const output = document.getElementById("orderData");
  const button = document.getElementById("searchBtn");

  function isValidUUID(uuid) {
    const regex =
      /^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i;
    return regex.test(uuid);
  }

  function fetchOrder() {
    const uuid = input.value.trim();

    if (!isValidUUID(uuid)) {
      alert("Некорректный UUID!");
      output.style.display = "none"; // скрываем контейнер при неверном UUID
      return;
    }

    fetch(`/order/${uuid}?json=1`)
      .then((res) => {
        if (!res.ok) throw new Error("Заказ не найден");
        return res.json();
      })
      .then((data) => {
        const formatted = JSON.stringify(data, null, 2);
        output.style.display = "block";
        output.textContent = formatted;
      })
      .catch((err) => {
        alert(err.message);
        output.style.display = "none"; // скрываем контейнер при ошибке
      });
  }

  button.addEventListener("click", fetchOrder);

  input.addEventListener("keydown", function (event) {
    if (event.key === "Enter") {
      fetchOrder();
    }
  });
});
