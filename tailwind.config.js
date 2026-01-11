module.exports = {
    content: [
        "./ui/html/**/*.html"
    ],
    theme: {
        extend: {
            colors: {
                bg: {
                    main: "#0b0f14",
                    surface: "#111827",
                    elevated: "#151b26"
                },
                border: {
                    subtle: "#1f2937"
                },
                text: {
                    primary: "#e5e7eb",
                    secondary: "#9ca3af",
                    muted: "#6b7280"
                },
                accent: {
                    DEFAULT: "#22d3ee",
                    hover: "#67e8f9"
                },
                success: "#22c55e",
                danger: "#ef4444"
            }
        }
    },
    plugins: []
}
