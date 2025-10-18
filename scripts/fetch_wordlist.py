#!/usr/bin/env python3
"""
Fetch and normalize an open-source English word list for games.

- Source: dwyl/english-words (words_alpha.txt)
  URL: https://raw.githubusercontent.com/dwyl/english-words/master/words_alpha.txt
- Output: normalized, lowercase, ASCII-only [a-z], deduplicated, sorted

Usage:
  python scripts/fetch_wordlist.py --output data/wordlists/english_words.txt \
      --min-length 2 --max-length 0 --allow-hyphens 0 --allow-apostrophes 0

Notes:
- Default filters exclude hyphens and apostrophes and words shorter than 2 chars.
- Set --max-length 0 to disable max-length filtering.
"""

from __future__ import annotations

import argparse
import io
import os
import re
import sys
import urllib.request
from typing import Iterable, Set

SOURCE_URL = (
    "https://raw.githubusercontent.com/dwyl/english-words/master/words_alpha.txt"
)


def download_words(url: str) -> Iterable[str]:
    request = urllib.request.Request(
        url,
        headers={
            "User-Agent": "wordlist-fetcher/1.0 (+https://github.com/dwyl/english-words)"
        },
    )
    with urllib.request.urlopen(request, timeout=60) as response:
        # words_alpha.txt is ASCII/UTF-8, one word per line
        for raw_line in io.TextIOWrapper(response, encoding="utf-8", errors="ignore"):
            yield raw_line.rstrip("\n\r")


def build_pattern(allow_hyphens: bool, allow_apostrophes: bool) -> re.Pattern[str]:
    chars = "a-z"
    if allow_hyphens:
        chars += "-"
    if allow_apostrophes:
        chars += "'"
    return re.compile(rf"^[{chars}]+$")


def normalize_words(
    words: Iterable[str],
    *,
    min_length: int,
    max_length: int | None,
    allow_hyphens: bool,
    allow_apostrophes: bool,
) -> Set[str]:
    pattern = build_pattern(allow_hyphens=allow_hyphens, allow_apostrophes=allow_apostrophes)
    normalized: Set[str] = set()

    for word in words:
        if not word:
            continue
        w = word.strip().lower()
        if not w:
            continue
        if not pattern.match(w):
            continue
        if len(w) < min_length:
            continue
        if max_length is not None and len(w) > max_length:
            continue
        normalized.add(w)

    return normalized


def write_output(words: Set[str], output_path: str) -> None:
    os.makedirs(os.path.dirname(output_path), exist_ok=True)
    with open(output_path, "w", encoding="utf-8") as f:
        for w in sorted(words):
            f.write(w)
            f.write("\n")


def parse_args(argv: list[str]) -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Fetch and normalize an English word list")
    default_output = os.path.join(
        os.path.dirname(os.path.dirname(__file__)), "data", "wordlists", "english_words.txt"
    )
    parser.add_argument(
        "--output",
        default=default_output,
        help="Path to write the normalized word list (default: %(default)s)",
    )
    parser.add_argument(
        "--min-length",
        type=int,
        default=2,
        help="Minimum word length to include (default: %(default)s)",
    )
    parser.add_argument(
        "--max-length",
        type=int,
        default=0,
        help="Maximum word length to include; 0 disables (default: %(default)s)",
    )
    parser.add_argument(
        "--allow-hyphens",
        type=int,
        choices=(0, 1),
        default=0,
        help="Allow hyphens '-' in words (default: %(default)s)",
    )
    parser.add_argument(
        "--allow-apostrophes",
        type=int,
        choices=(0, 1),
        default=0,
        help="Allow apostrophes '\'' in words (default: %(default)s)",
    )

    args = parser.parse_args(argv)

    max_length: int | None = None if args.max_length == 0 else args.max_length
    args.max_length = max_length

    return args


def main(argv: list[str]) -> int:
    args = parse_args(argv)

    try:
        source_iter = download_words(SOURCE_URL)
        words = normalize_words(
            source_iter,
            min_length=args.min_length,
            max_length=args.max_length,
            allow_hyphens=bool(args.allow_hyphens),
            allow_apostrophes=bool(args.allow_apostrophes),
        )
        write_output(words, args.output)
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        return 1

    print(
        f"Wrote {len(words):,} words to {args.output}"
        f" (min_len={args.min_length}, max_len={args.max_length or 'none'},"
        f" hyphens={bool(args.allow_hyphens)}, apostrophes={bool(args.allow_apostrophes)})"
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main(sys.argv[1:]))
