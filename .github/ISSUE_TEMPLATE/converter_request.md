---
name: Converter Request
about: Request a new format converter for SB
title: "[CONVERTER] Add {FORMAT} converter"
labels: converter, enhancement
assignees: ''
---

## Converter Details

**Input Format(s)**: (e.g., HEIC, HEIF)
**Output Format**: (e.g., JPG)
**Converter Name**: (e.g., `sb jpg`)

## Use Case

Why do you need this converter? What workflow does it enable?

## Format Details

- **Common File Extensions**: (e.g., .heic, .heif)
- **Typical Use Case**: (e.g., iPhone photos, screenshots)
- **Special Requirements**: (e.g., lossy/lossless, metadata preservation)

## Quality/Options

What options should this converter support?

- [ ] Quality setting (e.g., compression level)
- [ ] Color space conversion
- [ ] Metadata preservation
- [ ] Resolution/scaling
- [ ] Other: _______________

## FFmpeg Support

Does FFmpeg support this conversion?

- [ ] Yes, natively
- [ ] Yes, with specific codec/filter
- [ ] No (requires different tool)
- [ ] Unknown

**FFmpeg Command Example** (if known):
```bash
ffmpeg -i input.heic output.jpg
```

## Priority

How important is this converter to you?

- [ ] Critical (blocking my workflow)
- [ ] High (would significantly improve workflow)
- [ ] Medium (nice to have)
- [ ] Low (just a suggestion)

## Willingness to Contribute

- [ ] I'm willing to implement this converter
- [ ] I can provide test files
- [ ] I can help with documentation
- [ ] I can test the implementation

## Additional Context

Add any other context, examples, or references about the converter request here.
