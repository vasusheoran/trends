/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./templates/**/*.{html,js,templ}"],
    theme: {
      extend: {
        gridTemplateColumns: {
          // Simple 16 column grid
          '15': 'repeat(15, minmax(0, 1fr))',

          // Complex site-specific column configuration
          'footer': '200px minmax(900px, 1fr) 100px',
        }
      },
    },
    plugins: [
      require('@tailwindcss/forms'),
      require('@tailwindcss/typography'),
    ],
  }
  