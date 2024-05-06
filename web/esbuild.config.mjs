import * as esbuild from 'esbuild'
import { default as pkg } from './package.json' with {type: 'json'}
import stylePlugin from 'esbuild-style-plugin'
import tailwindcss from "tailwindcss"
import autoprefixer from "autoprefixer";
import chalk from "chalk";
import imageInline from "esbuild-plugin-inline-image";
import { writeFile } from "fs";

const isProduction = process.env.NODE_ENV === "production"
const args = process.argv.slice(2)
const watch = args.indexOf("--watch") >= 0
const serve = args.indexOf("--serve") >= 0

// (optionally) load secrets and write appProps.js
import("./secretProps.json", { assert: { type: "json" } })
    .then((appProps) => {
        writeFile(
            `./public/appProps.js`,
            `export default ${JSON.stringify(appProps.default)}`,
            err => {
                if (err) {
                    console.error(err)
                    return
                }
                console.log("public/appProps.js refreshed")
            }
        )
    })
    .catch(() => console.log("secretProps.json not found, using default values"))

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
    plugins: [
        stylePlugin({postcss: {plugins: [tailwindcss, autoprefixer]}}),
        imageInline({limit: 10240, extensions: ["png", "jpg", "jpeg", "gif"]}), // 10Kb
    ],
})

// watch and serve
if (watch || serve) {
    if (watch) await ctx.watch()
    if (serve) {
        let {host, port} = await ctx.serve({
            host: "0.0.0.0",
            port: parseInt(process.env.PORT),
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
