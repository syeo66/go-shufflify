document.addEventListener('DOMContentLoaded', () => {
  let style = document.createElement('style')

  function doResize() {
    const elm = document.querySelector('.queue')
    if (elm) {
      const top = elm.getBoundingClientRect().top
      const viewPortHeight = window.innerHeight
      const maxHeight = viewPortHeight - top - 16;
      style.remove()
      style = document.createElement('style')
      style.innerHTML = `
        .queue {
          max-height: ${maxHeight}px;
        }
      `
      document.head.appendChild(style)
    }
  }

  window.addEventListener('resize', doResize)
  doResize()

  setInterval(doResize, 5000)
})
