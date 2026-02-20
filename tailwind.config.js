/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './pages/**/*.{js,ts,jsx,tsx,mdx}',
    './components/**/*.{js,ts,jsx,tsx,mdx}',
    './app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        surface: '#f7f7f8',
        'surface-dark': '#111113',
        text: {
          primary: '#1d1d22',
          secondary: '#55556d',
          tertiary: '#8e8ea0',
        },
        border: '#e5e5ea',
        'border-dark': '#2c2c30',
        accent: '#2563eb',
        'accent-hover': '#1d4ed8',
      },
      fontFamily: {
        sans: ['var(--font-geist-sans)'],
        mono: ['var(--font-geist-mono)'],
      },
      maxWidth: {
        container: '800px',
      },
    },
  },
  plugins: [],
}
