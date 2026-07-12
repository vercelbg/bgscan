# bgscan Documentation

This directory contains the source for the bgscan documentation site, built with [Hugo](https://gohugo.io/) and the [Hugo Book theme](https://github.com/alex-shpak/hugo-book).

## 📖 About This Documentation

The live documentation is published at:
**https://mohsenbg.github.io/bgscan/docs**

If you are looking for the bgscan project itself (source code, releases, installation), please see the parent directory or visit:
**https://github.com/MohsenBg/bgscan**

## 🛠️ Running the Documentation Locally

To preview or contribute to the documentation:

### Prerequisites
- [Hugo extended version](https://gohugo.io/getting-started/installing/) (v0.121.0 or newer recommended)
- Git

### Quick Start
```bash
# Clone the main repository (if you haven't already)
git clone https://github.com/MohsenBg/bgscan.git
cd bgscan/docs

# Start the Hugo development server
hugo server
```

The site will be available at http://localhost:1313/bgscan/docs/ (note the `/bgscan/docs/` path due to the repository's baseURL).

## 📝 Editing the Documentation

Content is written in Markdown and organized by language:

- **English**: `content/en/`
- **Persian (Farsi)**: `content/fa/`

Each language directory follows the same structure:
```
content/<lang>/
├── getting-started/
├── scanner/
├── settings/
├── developer/
└── ...
```

To add or modify a page:
1. Navigate to the appropriate language and section.
2. Edit the corresponding `.md` file (or create a new one).
3. Front matter (YAML at the top of each file) controls title, weight, and menu behavior.
4. Save and preview locally with `hugo server`.

### Images
Place images in `static/` and reference them in Markdown as:
```markdown
![Alt text](/bgscan-image-name.webp)
```
(The leading `/bgscan-` matches the site's baseURL configuration.)

## 🔧 Configuration

The site configuration lives in [`hugo.yaml`](./hugo.yaml). Key settings:
- `baseURL`: Set to `https://mohsenbg.github.io/bgscan/` for production.
- `languages`: Defines English (`en`) and Persian (`fa`) with respective content directories.
- Menus, theme, and output formats are also defined here.

## 🚦 Building for Production

The documentation is automatically built and deployed to GitHub Pages via the [.github/workflows/hugo.yaml](https://github.com/MohsenBg/bgscan/blob/main/.github/workflows/hugo.yaml) workflow on pushes to `main`.

To build locally:
```bash
hugo --minify
```
The output will be generated in the `public/` directory.

## 📄 License

The documentation is licensed under the same MIT license as the bgscan project. See the [root LICENSE](../LICENSE) file for details.