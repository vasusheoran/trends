/** @type {import('tailwindcss').Config} */
module.exports = {
    darkMode: 'media',
    content: [
        "./templates/**/*.{html,js,templ}",
        "./node_modules/flowbite/**/*.js",
        "/static/history.js"
    ],
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
        require('flowbite/plugin')({
            datatables: true,
        }),
    ],
  }
  