const displayButton = document.getElementById("click-me")

function displayLyrics() {
  alert("Too bad crodie")
  console.log("I watch the moon\n I watch you\n So long nice to know you, I'll be moving on")
}

displayButton.addEventListener("click", (_e) => {
  displayLyrics()
})
