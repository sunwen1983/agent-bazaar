#!/usr/bin/env python3
"""Merge a second nsdeploy static bundle (guide) into the origin-C bundle (stays).

nsdeploy generates one bundle per network, but origin C hosts two networks
(stays + local-guide) behind one manifest. This script:

- copies <guide>.network.json into the stays bundle,
- appends the guide network entry to agent-manifest.json `networks` and to
  `linking.exports`,
- merges links.json `networks` maps,
- adds the guide network JSON to _headers with Access-Control-Allow-Origin: *.

Usage: merge-origin-c.py <origin-c-dir> <guide-bundle-dir>
"""
import json
import sys
from pathlib import Path


def load(path: Path) -> dict:
    with open(path, encoding="utf-8") as f:
        return json.load(f)


def save(path: Path, data: dict) -> None:
    with open(path, "w", encoding="utf-8") as f:
        json.dump(data, f, indent=2)
        f.write("\n")


def main() -> None:
    origin_c = Path(sys.argv[1])
    guide_dir = Path(sys.argv[2])

    guide_networks = sorted(guide_dir.glob("*.network.json"))
    if len(guide_networks) != 1:
        raise SystemExit(f"expected exactly one *.network.json in {guide_dir}, "
                         f"found {[p.name for p in guide_networks]}")
    guide_json = guide_networks[0]
    target = origin_c / guide_json.name
    target.write_bytes(guide_json.read_bytes())
    print(f"copied {guide_json.name} -> {target}")

    manifest_path = origin_c / "agent-manifest.json"
    manifest = load(manifest_path)
    guide_manifest = load(guide_dir / "agent-manifest.json")

    names = {n["name"] for n in manifest.get("networks", [])}
    for entry in guide_manifest.get("networks", []):
        if entry["name"] in names:
            raise SystemExit(f"network {entry['name']} already in manifest")
        manifest.setdefault("networks", []).append(entry)
        print(f"manifest: added network {entry['name']}")

    exports = manifest.setdefault("linking", {}).setdefault("exports", [])
    have = {e["network"] for e in exports}
    for entry in guide_manifest.get("linking", {}).get("exports", []):
        if entry["network"] not in have:
            exports.append(entry)
            print(f"manifest: added linking export {entry['network']}")

    links_path = origin_c / "links.json"
    links = load(links_path) if links_path.exists() else {}
    guide_links = load(guide_dir / "links.json") if (guide_dir / "links.json").exists() else {}
    merged = links.get("networks", {})
    merged.update(guide_links.get("networks", {}))
    links["networks"] = merged
    save(links_path, links)
    save(manifest_path, manifest)

    headers_path = origin_c / "_headers"
    headers = headers_path.read_text(encoding="utf-8") if headers_path.exists() else ""
    rule = f"/{guide_json.name}\n  Access-Control-Allow-Origin: *\n"
    if guide_json.name not in headers:
        headers = headers.rstrip("\n") + "\n" + rule
        headers_path.write_text(headers, encoding="utf-8")
        print(f"_headers: added CORS rule for /{guide_json.name}")

    print("origin C merged OK:", sorted(p.name for p in origin_c.glob("*.network.json")))


if __name__ == "__main__":
    main()
