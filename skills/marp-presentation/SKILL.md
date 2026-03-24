---
name: marp-presentation
description: Generate professional slide decks (PPTX, PDF) from Markdown files using the Marp CLI.
---

# Marp Presentation Skill

Use this skill when the user wants to generate a presentation file (PowerPoint .pptx or PDF) from a Markdown document.

## Prerequisites
- The system must have **Marp CLI** installed and accessible via the `marp` command.
- The input file must be a valid Markdown file, preferably with Marp directives (e.g., `--- marp: true ---` in frontmatter).

## Usage

### 1. Basic Conversion

To convert a Markdown file to PowerPoint (PPTX):
```bash
marp --pptx <input_file.md> -o <output_file.pptx>
```

To convert a Markdown file to PDF:
```bash
marp --pdf <input_file.md> -o <output_file.pdf>
```

### 2. Enabling Marp in Markdown

If the input Markdown file does not have the Marp header, you should add it or specify `--engine` settings, but usually, it's best to ensure the file starts with:

```yaml
---
marp: true
theme: default
paginate: true
---
```

### 3. Image Handling

Ensure that all images referenced in the Markdown file are accessible.
- If images are local (e.g., `./assets/image.png`), run the `marp` command from the directory where the relative paths resolve correctly, or allow Marp to access local files.
- `marp` automatically handles embedding local images into the final PPTX/PDF.
- **Tip**: If you encounter security warnings about local files, add the `--allow-local-files` flag.
- **Tip**: If you encounter security warnings about local files, add the `--allow-local-files` flag.

### 4. Custom Themes

If the user requests a specific look, you can specify a theme:
- `default` (white/neutral)
- `gaia` (colorful)
- `uncover` (centered/minimal)

Example:
```yaml
---
marp: true
theme: gaia
class: lead
---
```

## Workflow for Agents

1.  **Verify Input**: Check if the target Markdown file exists.
2.  **Add Directives**: If the file is missing `marp: true`, suggest adding it or add it temporarily.
3.  **Run Command**: Execute the `marp` CLI command.
4.  **Verify Output**: Check if the output file (pptx/pdf) was created successfully.
5.  **Notify User**: Inform the user where the generated file is located.
