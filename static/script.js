document.addEventListener("DOMContentLoaded", () => {
  // State
  let selectedCoffee = {
    id: "espresso",
    name: "Espresso",
    icon: "☕",
  }
  let weight = 0
  let autoOrder = false

  // DOM Elements
  const coffeeButtons = document.querySelectorAll(".coffee-button")
  const weightDisplay = document.getElementById("weight-display")
  const autoOrderSwitch = document.getElementById("auto-order-switch")
  const activeCoffee = document.querySelector('.coffee-button.active')
  const orderNow = document.querySelector('.order-now')

  if (activeCoffee === null) {
    orderNow.disabled = true
  }

  // Coffee Selection
  coffeeButtons.forEach((button) => {
    button.addEventListener("click", function () {
      coffeeButtons.forEach((btn) => {
        btn.classList.remove("active")
        btn.style.backgroundColor = ""
      })

      // Remove active class from all buttons
      coffeeButtons.forEach((btn) => btn.classList.remove("active"))
      const color = this.dataset.color
      // Add active class to clicked button
      this.classList.add("active")
      orderNow.disabled = false
      this.style.backgroundColor = color

      // Update selected coffee
      selectedCoffee = {
        id: this.dataset.id,
        name: this.dataset.name,
        icon: this.dataset.icon,
      }
    })
  })

  // Auto Order Toggle
  autoOrderSwitch.addEventListener("change", function (event) {
    const hasActiveCoffee = document.querySelector('.coffee-button.active');

    if (!hasActiveCoffee) {
      event.preventDefault();
      this.checked = !this.checked;
      return;
    }

    autoOrder = this.checked
    document.querySelector('.order-now').disabled = this.checked;
    console.log(this.checked)
  })

  // Simulate weight changes
  weight += 1
  weightDisplay.textContent =weight + "oz"

  const socket = new WebSocket("ws://localhost/ws");

  socket.onOpen = () => {
    console.log('Connected to server');
  }

  socket.onMessage = (event) => {
    console.log(event.data);
  }
})
