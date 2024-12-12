import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./src/pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/components/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        background: "var(--foreground)",
        foreground: "var(--background)",
        black: "#000000", // Ensure black is defined
      },
      textColor: {
        DEFAULT: "#000000", // Set default text color to black
      },
    },
  },
  plugins: [],
};

export default config;
