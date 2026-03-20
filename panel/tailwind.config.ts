import type { Config } from 'tailwindcss'

const config: Config = {
  content: [
    './src/pages/**/*.{js,ts,jsx,tsx,mdx}',
    './src/components/**/*.{js,ts,jsx,tsx,mdx}',
    './src/app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      colors: {
        // White-only palette
        'accent-coral': '#ffffff',
        'accent-blue': '#ffffff',
        'accent-purple': '#ffffff',
        'accent-emerald': '#ffffff',
        'accent-rose': '#ffffff',
        'accent-amber': '#ffffff',
        'accent-cyan': '#ffffff',
        'accent-indigo': '#ffffff',
        // Neutral semantic colors
        'success': '#ffffff',
        'warning': '#ffffff',
        'error': '#ffffff',
        'info': '#ffffff',
      },
      backgroundImage: {
        'gradient-primary': 'linear-gradient(135deg, #ffffff 0%, #ffffff 100%)',
        'gradient-success': 'linear-gradient(135deg, #ffffff 0%, #ffffff 100%)',
        'gradient-warm': 'linear-gradient(135deg, #ffffff 0%, #ffffff 100%)',
      },
      boxShadow: {
        'glow-blue': '0 0 20px rgba(255, 255, 255, 0.3)',
        'glow-purple': '0 0 20px rgba(255, 255, 255, 0.3)',
        'glow-emerald': '0 0 20px rgba(255, 255, 255, 0.3)',
      },
      animation: {
        'pulse-glow': 'pulse-glow 2s cubic-bezier(0.4, 0, 0.6, 1) infinite',
        'slide-up': 'slide-up 0.3s ease-out',
      },
      keyframes: {
        'pulse-glow': {
          '0%, 100%': { opacity: '1' },
          '50%': { opacity: '0.7' },
        },
        'slide-up': {
          'from': { opacity: '0', transform: 'translateY(10px)' },
          'to': { opacity: '1', transform: 'translateY(0)' },
        },
      },
    },
  },
  plugins: [],
}
export default config
