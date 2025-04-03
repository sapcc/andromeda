#!/usr/bin/env bash

set -ex

if test $# -ne 3; then
    printf 'Usage: %s <PKG_URL> <TMP_DIR> <BIN_PATH>\n' "$(basename "$0")" >&2
    exit 1
fi

pkg_url="$1"
tmp_dir="$2"
bin_path="$3"

bin_basename="$(basename "$bin_path")"
pkg_extracted="${tmp_dir}/${bin_basename}-extracted"
pkg_path="${tmp_dir}/${bin_basename}.tar.gz"

curl -L -o "$pkg_path" "$pkg_url"
mkdir "$pkg_extracted"
tar -xzf "$pkg_path" -C "$pkg_extracted"
cp "${pkg_extracted}/${bin_basename}" "$bin_path"

rm -rf "$pkg_extracted" "$pkg_path"
