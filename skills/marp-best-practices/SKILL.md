---
name: marp-best-practices
description: Best practices for designing professional, bug-free presentations using Marp.
---

# Marp Design Best Practices

Use this skill to guide the creation of high-quality Marp presentations. Review these rules before generating any new deck.

## 1. Layout & Composition

### Split Slides (Text + Image)
The most professional layout for product presentations is the split view.
- **Usage**: Use `bg right` or `bg left` to place an image on one side and text on the other.
- **Standard Ratio**: `![bg right:40%](image.png)` (Image takes 40%, Text takes 60%).
- **Dense Text**: If you have bullet points > 4 lines or long paragraphs, give text more space by reducing image width: `![bg right:33%](image.png)`.

### Title Slides
Use the `lead` class for title slides to center everything gracefully.
```markdown
<!-- _class: lead -->
# My Presentation Title
![bg right:40%](hero-image.png)
```

## 2. Image Handling (Critical)

### Avoid Cropping Screenshots
By default, `bg` images use `cover` mode, which crops edges. For software screenshots (UI), you **MUST** use `contain` to show the full interface.

*   ❌ **Bad**: `![bg right](screenshot.png)` (Crops key UI elements)
*   ✅ **Good**: `![bg right contain](screenshot.png)` (Shows full UI)
*   ✅ **Precise**: `![bg right:40% contain](screenshot.png)` (Shows full UI, constrained width)

### Backgrounds
For decorative images (abstract, photos), standard `cover` (default) is fine.

## 3. Project Structure

### Local Assets
Always store images in a local `assets/` folder relative to the markdown file.
- **Structure**:
    ```text
    docs/marketing/presentations/
    ├── my-presentation.md
    └── assets/
        ├── screenshot1.png
        └── diagram.png
    ```
- **Reference**: `![bg right](assets/screenshot1.png)`

## 4. Troubleshooting Common Issues

### "Text Cutoff"
**Symptom**: Text runs off the bottom of the slide.
**Fix**:
1.  Reduce image width (give text 60-70% width).
2.  Split content into two slides.
3.  Use `<!-- fit -->` in the header (e.g. `# My Long Title <!-- fit -->`) to auto-scale text.

### "Missing Images"
**Symptom**: Images don't render in PPTX export.
**Fix**:
1.  Ensure paths are relative and correct.
2.  Use the flag `--allow-local-files` when running `marp`.

## 5. Standard Frontmatter
Start every file with this standard configuration:

```yaml
---
marp: true
theme: default
paginate: true
backgroundColor: #ffffff
---
```
