/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts}'],
  theme: {
    extend: {
      fontFamily: {
        display: ['"DM Sans"', 'sans-serif'],
        body: ['"Inter"', 'sans-serif'],
        mono: ['"JetBrains Mono"', 'monospace'],
      },
      colors: {
        glass: {
          white: 'rgba(255,255,255,0.08)',
          border: 'rgba(255,255,255,0.15)',
          dark: 'rgba(0,0,0,0.25)',
        },
        accent: {
          cyan:   '#00E5FF',
          violet: '#7C3AED',
          rose:   '#F43F5E',
          amber:  '#F59E0B',
          emerald:'#10B981',
        },
      },
      backdropBlur: {
        glass: '20px',
      },
      boxShadow: {
        glass: '0 8px 32px 0 rgba(0,0,0,0.37)',
        glow:  '0 0 24px rgba(0,229,255,0.3)',
        'glow-violet': '0 0 24px rgba(124,58,237,0.4)',
      },
      animation: {
        'fade-up':   'fadeUp 0.5s ease forwards',
        'pulse-glow': 'pulseGlow 2s ease-in-out infinite',
        'float':     'float 6s ease-in-out infinite',
      },
      keyframes: {
        fadeUp: {
          '0%':   { opacity: 0, transform: 'translateY(20px)' },
          '100%': { opacity: 1, transform: 'translateY(0)' },
        },
        pulseGlow: {
          '0%, 100%': { boxShadow: '0 0 20px rgba(0,229,255,0.2)' },
          '50%':      { boxShadow: '0 0 40px rgba(0,229,255,0.5)' },
        },
        float: {
          '0%, 100%': { transform: 'translateY(0px)' },
          '50%':      { transform: 'translateY(-8px)' },
        },
      },
    },
  },
  plugins: [],
}
