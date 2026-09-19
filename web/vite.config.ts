import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
  base: process.env.VITE_BASE_PATH ?? "/Plane-crazy-multitool/",
  plugins: [react()],
  test: { environment: "node", include: ["src/**/*.test.ts"] },
});
