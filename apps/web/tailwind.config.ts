import type { Config } from 'tailwindcss'

const config: Config = {
  darkMode: 'class',
  content: [
    './index.html',
    './src/**/*.{ts,tsx}',
  ],
  theme: {
    extend: {
      colors: {
        background: '#0a0a0f',
        foreground: '#e8e8f0',
        border: '#2a2a3a',
        accent: {
          DEFAULT: '#6c5ce7',
          hover: '#7c6cf7',
        },
        muted: {
          DEFAULT: '#1a1a24',
          foreground: '#9090a8',
        },
        card: {
          DEFAULT: '#22222e',
          foreground: '#e8e8f0',
        },
        destructive: {
          DEFAULT: '#ff6b6b',
          foreground: '#e8e8f0',
        },
        success: '#00b894',
        warning: '#fdcb6e',
        info: '#74b9ff',
      },
      fontFamily: {
        sans: ['Onest', 'system-ui', 'sans-serif'],
        mono: ['JetBrains Mono', 'monospace'],
      },
    },
  },
  plugins: [],
}

export default config
