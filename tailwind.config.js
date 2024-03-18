/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./**/*.{html,templ}",
    "./templates/**/*.templ",
  ],
  theme: {
    extend: {},
  },
  plugins: [
    require("@catppuccin/tailwindcss")({
      prefix: "ctp",
      defaultFlavour: "mocha",
    })
  ],
}

