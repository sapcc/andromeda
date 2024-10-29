import * as esbuild from 'esbuild'
import { default as pkg } from './package.json' with {type: 'json'}
import inlineImportPlugin from 'esbuild-plugin-inline-import';
import chalk from "chalk";
import imageInline from "esbuild-plugin-inline-image";
import postcss from "postcss";
import autoprefixer from "autoprefixer";
import tailwindcss from "tailwindcss"

const isProduction = process.env.NODE_ENV === "production"
const appProps = process.env.APP_PROPS || "{}"
const args = process.argv.slice(2)
const watch = args.indexOf("--watch") >= 0
const serve = args.indexOf("--serve") >= 0

// build app
let ctx = await esbuild.context({
    bundle: true,
    minify: isProduction,
    target: ["es2020", "chrome64", "firefox67", "safari11.1", "edge79"],
    format: "esm",
    platform: "browser",
    loader: {".js": "jsx"},
    sourcemap: !isProduction,
    entryPoints: [pkg.source],
    outdir: "public/build",
    define: {
        "process.env.APP_PROPS": appProps,
    },
    plugins: [
        inlineImportPlugin({
            transform: async (contents, args) => {
                let path = args.path
                if (path.endsWith(".css")) {
                    let {css} = await postcss([tailwindcss, autoprefixer]).process(contents, {
                        from: path,
                        to: path,
                    })
                    return css
                }
                return {contents}
            },
        }),
        imageInline({limit: 10240, extensions: ["png", "jpg", "jpeg", "gif"]}), // 10Kb
    ],
})

// watch and serve
if (watch || serve) {
    if (watch) await ctx.watch()
    if (serve) {
        let {host, port} = await ctx.serve({
            host: process.env.HOST || "127.0.0.1",
            port: parseInt(process.env.PORT || "3000"),
            servedir: "public",
        })
        console.log("serve on", `${host}:${port}`)
    }
} else {
    ctx
        .rebuild()
        .then(() => console.log(chalk.green("⚡ Build complete! ⚡")))
        .catch(() => process.exit(1))
    await ctx.dispose()
}
