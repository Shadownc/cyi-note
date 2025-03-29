/** @type {import('tailwindcss').Config} */
import typography from '@tailwindcss/typography'

export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        // 默认主题色
        primary: {
          50: 'var(--color-primary-50)',
          100: 'var(--color-primary-100)',
          200: 'var(--color-primary-200)',
          300: 'var(--color-primary-300)',
          400: 'var(--color-primary-400)',
          500: 'var(--color-primary-500)',
          600: 'var(--color-primary-600)',
          700: 'var(--color-primary-700)',
          800: 'var(--color-primary-800)',
          900: 'var(--color-primary-900)',
        },
        // 文本颜色
        'text-primary': 'var(--color-text)',
        'text-secondary': 'var(--color-text-secondary)',
        'text-tertiary': 'var(--color-text-tertiary, #6b7280)',
        'dark-text-primary': 'var(--color-text-dark, #f3f4f6)',
        'dark-text-secondary': 'var(--color-text-dark-secondary, #d1d5db)',
        'dark-text-tertiary': 'var(--color-text-tertiary-dark, #9ca3af)',
        // 背景颜色
        'background': 'var(--color-bg, #f9fafb)',
        'dark-background': 'var(--color-bg-dark, #1f2937)', 
        'card': 'var(--color-card-bg, #ffffff)',
        'dark-card': 'var(--color-card-bg-dark, #374151)',
        // 边框颜色
        'border': 'var(--color-border, #e5e7eb)',
        'dark-border': 'var(--color-border-dark, #4b5563)',
        // 扩展灰色系列
        gray: {
          750: '#2d3748', // 添加介于gray-700(#374151)和gray-800(#1f2937)之间的颜色
        },
        blue: {
          50: '#eef9ff',
          100: '#d9f1ff',
          200: '#bae5ff',
          300: '#8bd9ff',
          400: '#54c1ff',
          500: '#30a4ff',
          600: '#1a84ff',
          700: '#0e68e3',
          800: '#1156b8',
          900: '#144a91',
        },
      },
      boxShadow: {
        'inner': 'inset 0 2px 4px 0 rgba(0, 0, 0, 0.06)',
      },
      typography: {
        DEFAULT: {
          css: {
            maxWidth: '100%',
            a: {
              color: 'var(--color-primary-600)',
              '&:hover': {
                color: 'var(--color-primary-700)',
              },
            },
            code: {
              color: '#dd1144',
              backgroundColor: 'var(--color-code-bg)',
              padding: '0.2em 0.4em',
              borderRadius: '0.25em',
            },
            pre: {
              backgroundColor: 'var(--color-pre-bg)',
              color: 'var(--color-pre-text)',
            },
          },
        },
        dark: {
          css: {
            color: 'var(--color-text-dark)',
            a: {
              color: 'var(--color-primary-400)',
              '&:hover': {
                color: 'var(--color-primary-300)',
              },
            },
            h1: {
              color: 'var(--color-text-dark)',
            },
            h2: {
              color: 'var(--color-text-dark)',
            },
            h3: {
              color: 'var(--color-text-dark)',
            },
            h4: {
              color: 'var(--color-text-dark)',
            },
            strong: {
              color: 'var(--color-text-dark)',
            },
            blockquote: {
              color: 'var(--color-text-dark-secondary)',
            },
          },
        },
      },
    },
  },
  plugins: [
    typography,
  ],
} 