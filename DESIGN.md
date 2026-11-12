---
version: alpha
name: Stripe Inspired
website: "https://stripe.com"
description: An inspired interpretation of Stripe's design language — a financial-infrastructure brand built on a deep navy ink, an electric indigo primary, and a recurring atmospheric gradient mesh that occupies the upper third of nearly every marketing page. The system pairs the proprietary Sohne family at thin (300) weights with negative letter-spacing for editorial-density display headlines, and uses tabular-figure body type where money and numerics matter. Buttons are tight-radius pills, cards live on near-white surfaces, and the dashboard track flips polarity to a familiar dark-app shell.

seo:
  title: "Stripe Design System: Indigo, Söhne and Calm Gradients"
  metaDescription: "Stripe's electric indigo #533afd, thin Söhne headlines and gradient-mesh heroes in one DESIGN.md, with components for React, Go and AI tools."
  tags:
    - "Fintech & Crypto"
    - "Banking & Payments"
  highlights:
    - "Single voltage CTA — indigo #533afd is the only filled button on marketing surfaces, one per band"
    - "Sohne thin weight (300) display tier — negative tracking from -1.4px at 56px down to -0.2px at 20px"
    - "Tabular figures via tnum — every money and numeric cell carries the brand's quiet financial-data signal"
    - "Gradient mesh hero — cream, sherbet orange, lavender, indigo, and ruby pink washed across the upper third"
    - "Dark-app dashboard track — deep navy #1c1e54 mockups composited above the white canvas"
  lastUpdated: "2026-05-12"
  author:
    name: "Dov Azencot"
    url: "https://x.com/dovazencot"
  opening: |
    Stripe's design language is the modern reference for fintech marketing. One brand voltage (indigo #533afd). One proprietary typeface (Sohne at weight 300). One atmospheric gradient mesh — cream, sherbet orange, lavender, indigo, and ruby pink — washed across the upper third of nearly every marketing page. Deep navy ink #0d253d does all the body work, never pure black. Pill-shaped CTAs with tight 8x16px padding sit one-per-band, never two filled buttons competing in the same hero.

    This page packages the full system into a single DESIGN.md file. Inside: 22 colors, 16 type tokens, 6 corner radii, 8 spacing values, and 21 components — buttons, dashboard mockup chrome, pricing card variants, the gradient mesh backdrop, and the tabular-figure body type that quietly signals the brand's financial DNA. Every hex is quoted; every token name is reproducible.

    Feed the file to Claude, Cursor, or GitHub Copilot. The AI writes React components and Tailwind classes that match Stripe's editorial air — single-indigo CTA hierarchy, Sohne thin display tier, negative letter-spacing, dark-app dashboard composites — instead of a generic fintech theme. Inter at weight 300 with `font-feature-settings: "ss01"` substitutes for Sohne if you need an open-source path.
  related:
    - href: "https://github.com/google-labs-code/design.md"
      title: "The DESIGN.md specification"
      description: "Google Labs' open spec for machine-readable design system files — the format this page is built on."
    - href: "/design"
      title: "Browse all design systems"
      description: "The full directory of DESIGN.md files on shadcn.io, with live mockups for each."
    - href: "/blocks/pricing"
      title: "Pricing blocks for shadcn/ui"
      description: "Production-ready React pricing sections — including dark-featured-tier patterns inspired by Stripe."
  questions:
    - id: "primary-color"
      title: "What is Stripe's primary brand color?"
      answer: "Stripe's primary is indigo #533afd — an electric blue-violet used as the single filled CTA color across every marketing surface. It appears as `button-primary-pill` (the only filled button per band), inline link emphasis on light surfaces, and as a gradient mid-stop inside the mesh backdrop. Pressed state shifts to #2e2b8c (`primary-press`), and a paler #b9b9f9 (`primary-bg-subdued-hover`) backs the soft indigo tag pills. Ruby #ea2261 and magenta #f96bee live inside the gradient mesh and as chart accents — never as button fills."
    - id: "typography"
      title: "What typography does Stripe use, and what should I use if Sohne isn't available?"
      answer: 'Stripe runs Sohne — a proprietary variable typeface from Klim Type Foundry — at weights 300 (thin) and 400 (regular). The display tier sits at sizes 26 to 56px in weight 300 with negative letter-spacing from -0.26px at 26px down to -1.4px at 56px. The `ss01` stylistic set is enabled globally on the body, and `tnum` is applied per-element on money or numeric cells. If Sohne is unavailable, Inter at weight 300 with `font-feature-settings: "ss01"` and matching negative tracking is the closest open-source substitute — Helvetica and system-ui defaults are too heavy.'
    - id: "gradient-mesh"
      title: "What is the gradient mesh and why is it non-negotiable?"
      answer: "The gradient mesh is Stripe's signature atmospheric backdrop — a wide horizontal band of pastel cream #f5e9d4, sherbet orange #9b6829, lavender, indigo #533afd, and ruby pink #ea2261 washed across the upper third of nearly every marketing page. It's implemented as an SVG or large background image rather than a flat CSS gradient, because the real mesh has organic blob shapes that don't render cleanly in CSS. The mesh IS the brand's depth system — literal box-shadows are reserved for product UI mockups and stay subtle. A bare-canvas marketing hero feels off-brand."
    - id: "tabular-figures"
      title: "Why does Stripe use tabular figures everywhere?"
      answer: 'Tabular figures (`font-feature-settings: "tnum"`) are the brand''s quiet signal that it''s a financial-infrastructure platform. Every cell rendering currency, transaction amounts, or numeric counts uses tnum plus a tightening letter-spacing of -0.42px at 14px. The `body-tabular` token (14px / 300 / line-height 1.4 / -0.42px / tnum) is the canonical money body type. Without tnum, columns of numbers misalign and the financial-data signature collapses — apply it per-element on numeric content, not globally.'
    - id: "dark-pricing-tier"
      title: "Why is Stripe's featured pricing tier dark navy instead of indigo?"
      answer: "Stripe inverts polarity on the featured pricing card — `card-pricing-featured` uses background #1c1e54 (`brand-dark-900`) with white text rather than an indigo fill. This is the brand's distinctive featured-tier choice: indigo stays scarce as the CTA-only voltage, while the dark-navy fill gives the featured tier the visual weight it needs without doubling up on indigo. The same dark-navy fill carries the dashboard chrome and `button-on-dark` CTAs on dark surfaces. Indigo #533afd would over-saturate the page if it filled both the CTA and the featured tier."
    - id: "known-gaps"
      title: "What's missing from this DESIGN.md spec?"
      answer: "A handful of areas, captured in the Known Gaps section: precise hover states across the marketing system (only press states are documented), the full dashboard-product semantic palette for errors and success states (which lives inside the dashboard app, not the marketing site), Stripe Atlas / Climate / Connect sub-brand accents, the loading skeleton tokens, and a small set of motion specs for the gradient mesh subtle drift animation. The spec captures roughly the entire public marketing surface plus the dashboard chrome — sub-product semantic systems live on separate token sets."

colors:
  primary: "#533afd"
  primary-deep: "#4434d4"
  primary-press: "#2e2b8c"
  primary-soft: "#665efd"
  primary-bg-subdued-hover: "#b9b9f9"
  brand-dark-900: "#1c1e54"
  ink: "#0d253d"
  ink-secondary: "#273951"
  ink-mute: "#64748d"
  ink-mute-2: "#61718a"
  on-primary: "#ffffff"
  canvas: "#ffffff"
  canvas-soft: "#f6f9fc"
  canvas-cream: "#f5e9d4"
  hairline: "#e3e8ee"
  hairline-input: "#a8c3de"
  ruby: "#ea2261"
  magenta: "#f96bee"
  lemon: "#9b6829"
  shadow-blue: "#003770"

typography:
  display-xxl:
    fontFamily: "sohne-var, 'SF Pro Display', system-ui, -apple-system, sans-serif"
    fontSize: 56px
    fontWeight: 300
    lineHeight: 1.03
    letterSpacing: -1.4px
    fontFeature: ss01
  display-xl:
    fontFamily: "sohne-var, 'SF Pro Display', system-ui, -apple-system, sans-serif"
    fontSize: 48px
    fontWeight: 300
    lineHeight: 1.15
    letterSpacing: -0.96px
    fontFeature: ss01
  display-lg:
    fontFamily: "sohne-var, 'SF Pro Display', system-ui, -apple-system, sans-serif"
    fontSize: 32px
    fontWeight: 300
    lineHeight: 1.1
    letterSpacing: -0.64px
    fontFeature: ss01
  display-md:
    fontFamily: "sohne-var, 'SF Pro Display', system-ui, -apple-system, sans-serif"
    fontSize: 26px
    fontWeight: 300
    lineHeight: 1.12
    letterSpacing: -0.26px
    fontFeature: ss01
  heading-lg:
    fontFamily: "sohne-var, 'SF Pro Display', system-ui, -apple-system, sans-serif"
    fontSize: 22px
    fontWeight: 300
    lineHeight: 1.1
    letterSpacing: -0.22px
    fontFeature: ss01
  heading-md:
    fontFamily: "sohne-var, 'SF Pro Display', system-ui, -apple-system, sans-serif"
    fontSize: 20px
    fontWeight: 300
    lineHeight: 1.4
    letterSpacing: -0.2px
    fontFeature: ss01
  heading-sm:
    fontFamily: "sohne-var, 'SF Pro Display', system-ui, -apple-system, sans-serif"
    fontSize: 18px
    fontWeight: 300
    lineHeight: 1.4
    letterSpacing: 0
    fontFeature: ss01
  body-lg:
    fontFamily: "sohne-var, 'SF Pro Display', system-ui, -apple-system, sans-serif"
    fontSize: 16px
    fontWeight: 300
    lineHeight: 1.4
    letterSpacing: 0
    fontFeature: ss01
  body-md:
    fontFamily: "sohne-var, 'SF Pro Display', system-ui, -apple-system, sans-serif"
    fontSize: 15px
    fontWeight: 300
    lineHeight: 1.4
    letterSpacing: 0
    fontFeature: ss01
  body-tabular:
    fontFamily: "sohne-var, 'SF Pro Display', system-ui, -apple-system, sans-serif"
    fontSize: 14px
    fontWeight: 300
    lineHeight: 1.4
    letterSpacing: -0.42px
    fontFeature: tnum
  button-md:
    fontFamily: "sohne-var, 'SF Pro Display', system-ui, -apple-system, sans-serif"
    fontSize: 16px
    fontWeight: 400
    lineHeight: 1.0
    letterSpacing: 0
    fontFeature: ss01
  button-sm:
    fontFamily: "sohne-var, 'SF Pro Display', system-ui, -apple-system, sans-serif"
    fontSize: 14px
    fontWeight: 400
    lineHeight: 1.0
    letterSpacing: 0
    fontFeature: ss01
  caption:
    fontFamily: "sohne-var, 'SF Pro Display', system-ui, -apple-system, sans-serif"
    fontSize: 13px
    fontWeight: 400
    lineHeight: 1.4
    letterSpacing: -0.39px
    fontFeature: tnum
  micro:
    fontFamily: "sohne-var, 'SF Pro Display', system-ui, -apple-system, sans-serif"
    fontSize: 11px
    fontWeight: 300
    lineHeight: 1.4
    letterSpacing: 0
    fontFeature: ss01
  micro-cap:
    fontFamily: "sohne-var, 'SF Pro Display', system-ui, -apple-system, sans-serif"
    fontSize: 10px
    fontWeight: 400
    lineHeight: 1.15
    letterSpacing: 0.1px
    fontFeature: ss01

rounded:
  xs: 4px
  sm: 6px
  md: 8px
  lg: 12px
  xl: 16px
  pill: 9999px

spacing:
  xxs: 2px
  xs: 4px
  sm: 8px
  md: 12px
  lg: 16px
  xl: 24px
  xxl: 32px
  huge: 64px

components:
  button-primary-pill:
    backgroundColor: "{colors.primary}"
    textColor: "{colors.on-primary}"
    typography: "{typography.button-md}"
    rounded: "{rounded.pill}"
    padding: 8px 16px
  button-primary-pill-pressed:
    backgroundColor: "{colors.primary-press}"
    textColor: "{colors.on-primary}"
    typography: "{typography.button-md}"
    rounded: "{rounded.pill}"
    padding: 8px 16px
  button-secondary:
    backgroundColor: "{colors.canvas}"
    textColor: "{colors.primary}"
    typography: "{typography.button-md}"
    rounded: "{rounded.pill}"
    padding: 8px 16px
  button-on-dark:
    backgroundColor: "{colors.brand-dark-900}"
    textColor: "{colors.on-primary}"
    typography: "{typography.button-md}"
    rounded: "{rounded.pill}"
    padding: 8px 16px
  text-input:
    backgroundColor: "{colors.canvas}"
    textColor: "{colors.ink}"
    typography: "{typography.body-md}"
    rounded: "{rounded.sm}"
    padding: 8px 12px
  text-input-focused:
    backgroundColor: "{colors.canvas}"
    textColor: "{colors.ink}"
    typography: "{typography.body-md}"
    rounded: "{rounded.sm}"
    padding: 8px 12px
  card-feature-light:
    backgroundColor: "{colors.canvas}"
    textColor: "{colors.ink}"
    typography: "{typography.body-md}"
    rounded: "{rounded.lg}"
    padding: 32px
  card-pricing:
    backgroundColor: "{colors.canvas}"
    textColor: "{colors.ink}"
    typography: "{typography.body-md}"
    rounded: "{rounded.lg}"
    padding: 32px
  card-pricing-featured:
    backgroundColor: "{colors.brand-dark-900}"
    textColor: "{colors.on-primary}"
    typography: "{typography.body-md}"
    rounded: "{rounded.lg}"
    padding: 32px
  card-cream-band:
    backgroundColor: "{colors.canvas-cream}"
    textColor: "{colors.ink}"
    typography: "{typography.body-md}"
    rounded: "{rounded.lg}"
    padding: 32px
  card-dashboard-mockup:
    backgroundColor: "{colors.canvas}"
    textColor: "{colors.ink}"
    typography: "{typography.body-tabular}"
    rounded: "{rounded.lg}"
    padding: 24px
  pill-tag-soft:
    backgroundColor: "{colors.primary-bg-subdued-hover}"
    textColor: "{colors.primary-deep}"
    typography: "{typography.micro-cap}"
    rounded: "{rounded.pill}"
    padding: 4px 8px
  nav-bar-on-mesh:
    backgroundColor: "{colors.canvas}"
    textColor: "{colors.ink}"
    typography: "{typography.body-md}"
    rounded: "{rounded.xs}"
    padding: 16px 24px
  link-on-light:
    backgroundColor: "{colors.canvas}"
    textColor: "{colors.primary}"
    typography: "{typography.body-md}"
    rounded: "{rounded.xs}"
    padding: 0px
  footer-light:
    backgroundColor: "{colors.canvas}"
    textColor: "{colors.ink-mute}"
    typography: "{typography.caption}"
    rounded: "{rounded.xs}"
    padding: 64px 24px
---

## Overview

Stripe's design language opens with the gradient mesh. A wide horizontal band of pastel cream #f5e9d4, sherbet orange #9b6829, lavender, electric indigo #533afd, and ruby pink #ea2261 occupies the upper third of nearly every marketing page — the brand's instantly-recognizable atmospheric backdrop. Type and product UI mockups float above it on `{colors.canvas}` (#ffffff), with the gradient acting as both decoration and visual anchor. The lower portion of the page returns to white, with feature explanations on `{colors.canvas-soft}` (#f6f9fc — a barely-tinted cool off-white) and dashboard product mockups composited as faux IDE/console panels in deep navy #1c1e54.

The color system has two primary roles. **Indigo** (`{colors.primary}` — #533afd) is the brand's signature CTA color, used sparingly: one filled pill per band. **Deep navy** (`{colors.ink}` — #0d253d) is the universal body text color and the fill of dashboard mockups, the featured pricing tier, and the dark-app surfaces on the dashboard track. Ruby (`{colors.ruby}` — #ea2261) and magenta (`{colors.magenta}` — #f96bee) appear inside the gradient mesh and as accent dots in product UI mockups; they are not used as button colors.

Typography is built around **Sohne** at weight 300 with negative letter-spacing — the brand's editorial-density display signature. Display sizes (32–56px) use -1.4px to -0.64px tracking; body sizes use 0; tabular caption sizes (where money and numerics matter) use the OpenType `tnum` feature plus a tightening -0.36 to -0.42px tracking. The `ss01` stylistic set is enabled across all 16 type tokens.

**Key Characteristics:**

- Gradient-mesh backdrop on every marketing hero — cream #f5e9d4, sherbet orange #9b6829, lavender, indigo #533afd, and ruby pink #ea2261 horizontally washed across the upper third of the page.
- Single-indigo CTA hierarchy: filled `{colors.primary}` (#533afd) pill is the only filled button on marketing surfaces, one per band.
- Sohne thin (weight 300) display tier with negative tracking from -1.4px at 56px down to -0.2px at 20px.
- Tabular-figure body type (`body-tabular` — 14px / 300 / `tnum`) for any cell containing money or numerics — the brand's quiet financial-data signal.
- Dark-app dashboard track: deep navy #1c1e54 product UI mockups sit composited above the white canvas, frequently with rendered code or dashboard tables inside.
- Pill-shaped buttons (`{rounded.pill}` 9999px) with tight 8x16px padding — short, decisive, transactional.
- Cream-band feature cards (`{colors.canvas-cream}` — #f5e9d4) introduce a warm interlude between blue/white sections without breaking the brand's chromatic logic.

## Colors

> **Source pages:** home (`/`), `/payments`, `/pricing`, `dashboard.stripe.com/register/payments`.

### Brand & Accent

- **Indigo** (`{colors.primary}` — #533afd): The brand's signature CTA color. Filled-pill button, link emphasis, gradient anchor.
- **Indigo Deep** (`{colors.primary-deep}` — #4434d4): A deeper indigo used in gradient mid-stops and as the press-state warmer alternative.
- **Indigo Press** (`{colors.primary-press}` — #2e2b8c): Pressed-state lift of the primary.
- **Indigo Soft** (`{colors.primary-soft}` — #665efd): A lighter indigo used in product-UI accents and chart highlights.
- **Indigo Subdued** (`{colors.primary-bg-subdued-hover}` — #b9b9f9): Pale indigo fill used as soft tag background.
- **Brand Dark 900** (`{colors.brand-dark-900}` — #1c1e54): The deep navy used on the featured pricing tier and dashboard chrome.
- **Ruby** (`{colors.ruby}` — #ea2261): Gradient accent and chart highlight; never a button.
- **Magenta** (`{colors.magenta}` — #f96bee): Brighter pink stop in gradient meshes.
- **Lemon** (`{colors.lemon}` — #9b6829): Warm sherbet stop in gradient backdrops.

### Surface

- **Canvas** (`{colors.canvas}` — #ffffff): Default page background.
- **Canvas Soft** (`{colors.canvas-soft}` — #f6f9fc): Cool-tinted off-white used on feature bands beneath the gradient hero.
- **Canvas Cream** (`{colors.canvas-cream}` — #f5e9d4): Warm cream used as a feature-band fill — the brand's chromatic interlude.
- **Hairline** (`{colors.hairline}` — #e3e8ee): 1px borders on cards and tables.
- **Hairline Input** (`{colors.hairline-input}` — #a8c3de): Slightly cooler hairline used on form inputs.

### Text

- **Ink** (`{colors.ink}` — #0d253d): Default body text color across the brand. Deep navy, never pure black.
- **Ink Secondary** (`{colors.ink-secondary}` — #273951): Secondary text on white.
- **Ink Mute** (`{colors.ink-mute}` — #64748d): Helper text, captions, table labels.
- **Ink Mute 2** (`{colors.ink-mute-2}` — #61718a): Near-equivalent to ink-mute used in nav.
- **On Primary** (`{colors.on-primary}` — #ffffff): Text on indigo / dark-navy surfaces.

### Semantic

The brand does not use a separate semantic color palette in the marketing system — error / success states live in dashboard-product UI specifically.

## Typography

### Font Family

The display and UI tier is **Sohne** (proprietary, licensed from Klim Type Foundry) at weights 300 (thin) and 400 (regular). The variable font (`sohne-var`) is loaded with `font-feature-settings: "ss01"` enabled globally — the stylistic set substitutes a single-story `a` and other character variants that are part of the brand's typographic signature.

When Sohne is unavailable, fall back to **SF Pro Display** at thin weights, then system-ui. For maximum brand fidelity, **Inter** (open-source) at weight 300 with `font-feature-settings: "ss01"` and `letter-spacing: -1.4px` on display sizes approximates the rhythm closely.

### Hierarchy

| Token                       | Size | Weight | Line Height | Letter Spacing | Use                                  |
| --------------------------- | ---- | ------ | ----------- | -------------- | ------------------------------------ |
| `{typography.display-xxl}`  | 56px | 300    | 1.03        | -1.4px         | Hero headline                        |
| `{typography.display-xl}`   | 48px | 300    | 1.15        | -0.96px        | Section opener                       |
| `{typography.display-lg}`   | 32px | 300    | 1.1         | -0.64px        | Card title / sub-section             |
| `{typography.display-md}`   | 26px | 300    | 1.12        | -0.26px        | Compact card title                   |
| `{typography.heading-lg}`   | 22px | 300    | 1.1         | -0.22px        | Pricing tier name                    |
| `{typography.heading-md}`   | 20px | 300    | 1.4         | -0.2px         | Section sub-heading                  |
| `{typography.heading-sm}`   | 18px | 300    | 1.4         | 0              | Mini-section label                   |
| `{typography.body-lg}`      | 16px | 300    | 1.4         | 0              | Marketing body lead                  |
| `{typography.body-md}`      | 15px | 300    | 1.4         | 0              | Default UI body                      |
| `{typography.body-tabular}` | 14px | 300    | 1.4         | -0.42px        | Money / numeric tables (uses `tnum`) |
| `{typography.button-md}`    | 16px | 400    | 1.0         | 0              | Pill button label                    |
| `{typography.button-sm}`    | 14px | 400    | 1.0         | 0              | Compact pill label                   |
| `{typography.caption}`      | 13px | 400    | 1.4         | -0.39px        | Helper, table labels                 |
| `{typography.micro}`        | 11px | 300    | 1.4         | 0              | Fine print                           |
| `{typography.micro-cap}`    | 10px | 400    | 1.15        | 0.1px          | All-caps eyebrow                     |

### Principles

- **Thin weight is the brand.** Display tiers always render at weight 300. Bumping to 400+ removes the brand's editorial air.
- **Negative tracking on display.** -1.4px at 56px, scaling proportionally down to -0.2px at 20px. The negative tracking is the brand's typographic signature.
- **Tabular figures for money.** Any cell rendering currency, transaction amounts, or numeric counts uses `font-feature-settings: "tnum"` plus a tightening tracking. The brand quietly signals its financial DNA through this micro-detail.
- **`ss01` globally.** Apply `font-feature-settings: "ss01"` to the body element so the stylistic-set substitution is on for every text role.

### Note on Font Substitutes

Sohne is proprietary. Use **Inter** (open-source via Google Fonts) at weight 300 with `letter-spacing: -1.4px` and `font-feature-settings: "ss01"` for display tiers — Inter is the closest open-source analogue. For body sizes, Inter at 300 weight with `font-feature-settings: "tnum"` (where applicable) is the canonical substitute. Avoid Helvetica or system-ui defaults — they're heavier than the brand needs.

## Marketing surface

The public landing page and the auth screens use **Instrument Serif** for
editorial headings while the product shell stays on Inter. The pairing is
deliberate: the serif carries the marketing voice, Inter carries the work.
Display headings use Instrument Serif at 400 with tight negative tracking;
product page titles stay Inter at 300.

The landing page is also the live theme lab. Its palette strip writes the same
`color-theme` cookie as Settings → Appearance, and the gradient mesh reads the
active theme tokens instead of a fixed palette.

Motion is one intentional moment, not a layer of effects:

- a staggered hero load (`.animate-rise`),
- a 180ms route fade through the protected app layout,
- a slow mesh drift on the hero,
- press feedback on buttons and a small swatch scale on theme cards.

Everything except the functional feedback is gated with `motion-safe`. The
reduced-motion path is static and complete.

## Layout

### Spacing System

- **Base unit**: 8px (with 2 / 4 / 12 sub-tokens for fine work).
- **Tokens**: `{spacing.xxs}` 2px · `{spacing.xs}` 4px · `{spacing.sm}` 8px · `{spacing.md}` 12px · `{spacing.lg}` 16px · `{spacing.xl}` 24px · `{spacing.xxl}` 32px · `{spacing.huge}` 64px.
- **Section padding**: 64–96px on marketing surfaces; 32–48px on dashboard / product surfaces.
- **Card internal padding**: 32px on feature cards; 24px on dashboard mockups.

### Grid & Container

- Marketing pages center in a ~1200px container with the gradient mesh extending edge-to-edge above.
- Pricing collapses 4-up → 2-up → 1-up at 1024 / 768 breakpoints.
- Dashboard product mockups use their own internal grids (12-col tables, 3-col card grids) rendered as static composites.

### Whitespace Philosophy

The gradient mesh occupies the upper third of the page; the white canvas below is generously padded. Section gaps tend toward 96px, with content tightening to 32px on dashboard / pricing pages where users compare and act.

### App Shell & Page Anatomy

The authenticated product surface does not follow the marketing container.
These rules apply to every page inside the app shell.

- **Container.** The sidebar shell fills the viewport beside the sidebar, and pages use that width with `p-4` / `md:p-6` / `lg:p-8` padding. Do not center a page in a narrow column or cap its width. Constrain the content grid, not the page.
- **Field grids.** Form panels are `@container`s and grids use `grid gap-4 @lg:grid-cols-2 @3xl:grid-cols-3`, so fields follow the panel width, not the viewport. A field that needs the full row uses `@lg:col-span-2 @3xl:col-span-3`.
- **Page anatomy.** Every page renders, top to bottom: `PageHeader` (title, optional description, right-aligned actions), an optional toolbar for filters or search, the content (cards, table, or grid), and an optional pagination footer.
- **One title.** `PageHeader` (`web/src/components/common/page-header.tsx`) renders the single `h1`. Do not add a second title block, and do not hand-roll the title/actions row.
- **Vertical rhythm.** Page sections use `space-y-6`. Groups inside a section use `space-y-4`. Toolbars use `gap-2` or `gap-3`.
- **Tables.** A primary table sits in a bare `overflow-x-auto rounded-md border` wrapper. A secondary table may sit inside a Card. Do not mix both treatments on one page.
- **Stat rows.** Use `grid gap-4 sm:grid-cols-2 lg:grid-cols-4`, or `sm:grid-cols-3` for three stats. Content grids use `gap-6` and collapse to one column below `lg`.
- **Breakpoints.** The sidebar collapses to an icon rail and then a mobile sheet (shadcn Sidebar). Page headers stack below `sm`. Tables scroll horizontally; they do not reflow into cards.

### Page Composition

Every page has a main zone. At `xl` and up, most pages also have a context
rail. The rail carries context or secondary actions; it never repeats the
main zone.

| Archetype | Examples                 | Main zone                        | Context rail                                  |
| --------- | ------------------------ | -------------------------------- | --------------------------------------------- |
| Workbench | Users, Roles             | KPI strip, filters, table        | Selected record, saved views, recent activity |
| Record    | User edit, Settings      | Identity, form, related sections | Summary, metadata, danger zone                |
| Task      | User create, Role create | Tabs or form                     | Guidance: section help, policies, field hints |
| Exception | Auth, Error              | Centered card                    | None                                          |

- Grid: `grid gap-6 xl:grid-cols-12`. Main zone `xl:col-span-8`, rail `xl:col-span-4`. Below `xl` the rail stacks under the main zone.
- The rail holds real content: a summary, guidance, history, or actions. Do not use it for decoration or filler stats.
- Every card earns its space. A card with one line of content is merged or enriched.
- Empty states carry one sentence and one action. No blank panels.
- Tables keep the full main-zone width; the rail never squeezes them below about 900px.

## Elevation & Depth

| Level | Treatment                                                                   | Use                                                                             |
| ----- | --------------------------------------------------------------------------- | ------------------------------------------------------------------------------- |
| 0     | Flat                                                                        | Default surface                                                                 |
| 1     | `box-shadow: rgba(0,55,112,0.08) 0 1px 3px`                                 | Card lift on white                                                              |
| 2     | `box-shadow: rgba(0,55,112,0.08) 0 8px 24px, rgba(0,55,112,0.04) 0 2px 6px` | Floating panels, dashboard mockup chrome                                        |
| 3     | Gradient mesh backdrop                                                      | The brand's primary depth medium — atmospheric color rather than literal shadow |

### Decorative Depth

The gradient mesh IS the depth system. Implemented as a layered SVG or large background image rather than CSS gradients (the actual mesh has organic blob shapes that aren't CSS-renderable). The mesh provides the brand's signature lift; literal shadows are reserved for product-UI mockups and stay subtle.

## Shapes

### Border Radius Scale

| Token            | Value  | Use                             |
| ---------------- | ------ | ------------------------------- |
| `{rounded.xs}`   | 4px    | Hairline tags, table chrome     |
| `{rounded.sm}`   | 6px    | Form inputs                     |
| `{rounded.md}`   | 8px    | Compact cards, alerts           |
| `{rounded.lg}`   | 12px   | Pricing cards, feature cards    |
| `{rounded.xl}`   | 16px   | Dashboard product mockup chrome |
| `{rounded.pill}` | 9999px | All buttons, tag pills          |

### Photography Geometry

The brand uses **product UI mockups** more than photography. Dashboard composites render as faux IDE/terminal/dashboard chrome inside `{rounded.lg}` 12px containers with a subtle `box-shadow`. Real photography appears in customer logo strips and the rare case-study card; treated as inset 4:3 with no shadow.

## Components

### Buttons

**`button-primary-pill`** — the dominant CTA system-wide.

- Background `{colors.primary}` (#533afd), text `{colors.on-primary}` (#ffffff), type `{typography.button-md}`, padding `{spacing.sm} {spacing.lg}` (8px 16px), rounded `{rounded.pill}` 9999px.
- Pressed state `button-primary-pill-pressed` shifts background to `{colors.primary-press}` (#2e2b8c).

**`button-secondary`** — outline-style alternative.

- Background `{colors.canvas}` (#ffffff), text `{colors.primary}` (#533afd), 1px solid `{colors.primary}` border, same pill geometry.

**`button-on-dark`** — used on dashboard / dark surfaces.

- Background `{colors.brand-dark-900}` (#1c1e54), text `{colors.on-primary}`, same pill geometry.

### Cards & Containers

**`card-feature-light`** — feature explanation card on white.

- Background `{colors.canvas}` (#ffffff), padding `{spacing.xxl}` (32px), rounded `{rounded.lg}` 12px, 1px `{colors.hairline}` (#e3e8ee) border, optional Level 1 shadow.

**`card-pricing`** — standard pricing tier.

- Background `{colors.canvas}` (#ffffff), padding `{spacing.xxl}` (32px), rounded `{rounded.lg}`, 1px `{colors.hairline}` border. Title `{typography.heading-lg}`, price `{typography.display-md}`, body `{typography.body-md}`, CTA pinned bottom as `button-primary-pill`.

**`card-pricing-featured`** — the inverted dark featured tier.

- Background `{colors.brand-dark-900}` (#1c1e54), text `{colors.on-primary}`, otherwise identical structure to `card-pricing`. The deep-navy fill is the brand's distinctive featured-tier choice.

**`card-cream-band`** — warm interlude card.

- Background `{colors.canvas-cream}` (#f5e9d4), text `{colors.ink}` (#0d253d), padding 32px, rounded `{rounded.lg}`. Used to break up the indigo / white rhythm with warmth.

**`card-dashboard-mockup`** — composited dashboard / product UI screenshot.

- Background `{colors.canvas}`, type `{typography.body-tabular}` (with `tnum`), padding `{spacing.xl}` 24px, rounded `{rounded.lg}` 12px, Level 2 shadow. Often contains nested mini-mockups: code preview + dashboard table + chart card.

### Inputs & Forms

**`text-input`** — standard form field.

- Background `{colors.canvas}` (#ffffff), text `{colors.ink}` (#0d253d), type `{typography.body-md}`, padding `{spacing.sm} {spacing.md}` (8px 12px), rounded `{rounded.sm}` 6px, 1px `{colors.hairline-input}` (#a8c3de) border.
- Focus state `text-input-focused`: border swaps to `{colors.primary}` (#533afd).

**Form page** — the layout for create and edit screens.

- The form fills the page width. Fields use `grid gap-4 sm:grid-cols-2 xl:grid-cols-3`; do not cap the page to make inputs narrower.
- The form is one `Card`. Do not wrap it in a second card; nested cards are not a sectioning device.
- Forms with two or more sections use tabs: a segmented `TabsList` in a bordered strip at the top of the card, one `TabsContent` per section with its own heading and description.
- Tab labels are short nouns ("Account", "Employment"). The full section title and its description stay inside the panel.
- One submit footer closes the card: Cancel (outline) and the primary action, separated by `border-t`.
- Fields sit in a `grid gap-4 @lg:grid-cols-2 @3xl:grid-cols-3` inside an `@container` panel; a field that needs the full width uses `@lg:col-span-2 @3xl:col-span-3`.
- Server validation errors mark the affected tab with a destructive dot. A failed submit switches to the first tab that has an error.
- Every field is `FieldLabel` + control + `InputError`. Required and optional hints belong in `FieldLabel`, not in the label text.

### Navigation

**`nav-bar-on-mesh`** — top nav floating over the gradient hero.

- Background `{colors.canvas}` (or transparent depending on scroll), text `{colors.ink}`, padding `{spacing.lg} {spacing.xl}` (16px 24px). Logo wordmark on the left, primary nav center, sign-in + filled `button-primary-pill` on the right.

### Pills, Tags, and Chips

**`pill-tag-soft`** — subdued indigo tag.

- Background `{colors.primary-bg-subdued-hover}` (#b9b9f9), text `{colors.primary-deep}` (#4434d4), type `{typography.micro-cap}`, padding 4px 8px, rounded `{rounded.pill}`.

### Signature Components

**Gradient Mesh Backdrop** — pastel cream #f5e9d4 → sherbet orange #9b6829 → lavender → indigo #533afd → ruby pink #ea2261 stops blurred horizontally across the upper third of the page. Implemented as SVG or a large background image — not a flat CSS gradient (the real mesh has organic blob shapes).

**Composited Dashboard Mockup** — multi-layer faux-product-UI compositions: an IDE panel on the left, a dashboard table center, a chart card on the right, all rendered at small scale inside `{rounded.lg}` containers with subtle Level 2 shadows. The composite is the brand's most-photographed feature.

**Tabular-Figure Money Type** — every number rendering money, count, or transaction value uses `font-feature-settings: "tnum"`. The brand's quiet signal that it's a financial-infrastructure platform.

**`link-on-light`** — inline links on light surfaces.

- Text `{colors.primary}` (#533afd) rendered in `{typography.body-md}`, no underline by default.

**`footer-light`** — site-wide footer.

- Background `{colors.canvas}` (#ffffff), text `{colors.ink-mute}` (#64748d), type `{typography.caption}`, padding `{spacing.huge} {spacing.xl}` (64px 24px). Holds 4–6 columns of link groups, social icons, and a small legal row.

## Do's and Don'ts

### Do

- Reserve `{colors.primary}` (#533afd) for filled CTAs and inline link emphasis — it should appear sparingly, one filled button per band.
- Apply the gradient mesh to every marketing hero; bare-canvas heroes feel off-brand.
- Render display tiers at weight 300 with negative letter-spacing — the thin tracking is the typographic signature.
- Use `font-feature-settings: "tnum"` on every money / numeric cell.
- Apply `font-feature-settings: "ss01"` globally on the body element.
- Pair every feature explanation with a composited product UI mockup; the brand's argument is "look at the actual product."

### Don't

- Don't bump display weight above 300 — at 400 the brand's editorial air collapses.
- Don't add new accent colors outside the documented gradient stops (cream / orange / lavender / indigo / ruby / magenta).
- Don't use the indigo `{colors.primary}` (#533afd) as a body-text color — it's a CTA and link color, not a type color at body size.
- Don't shrink button padding below 8px 16px — the tight pill is part of the brand's transactional feel.
- Don't render money cells without `tnum` — it breaks the quiet financial-data signature.
- Don't replace the pill shape with rounded-rectangles for buttons.

## Responsive Behavior

### Breakpoints

| Name    | Width       | Key Changes                                                        |
| ------- | ----------- | ------------------------------------------------------------------ |
| Wide    | ≥ 1440px    | Full gradient mesh edge-to-edge; dashboard composite at full scale |
| Desktop | 1024–1440px | Default content max-width; pricing 4-up                            |
| Tablet  | 768–1023px  | Pricing 2-up; dashboard composite simplifies to 2 panels           |
| Mobile  | < 768px     | Pricing 1-up; hamburger nav; display drops 56 → 36px               |

### Touch Targets

- Pill buttons hit ≥ 40×40px on mobile via padding scaling. On smaller screens, buttons size up to 44×44px to maintain WCAG AAA.
- Form fields stay at 40px minimum height.

### Collapsing Strategy

- Display tiers stair-step 56 → 48 → 32 → 26 → 22px through the breakpoints.
- Gradient mesh re-tiles on mobile to preserve the wash without disappearing.
- Dashboard composites simplify to single-panel mockups on mobile; the multi-layer composition only renders at desktop+.
- Pricing tiers stair-step 4-up → 2-up → 1-up.

### Image Behavior

Product UI composites use `srcset` with art-direction crops at major breakpoints. Mobile crops focus on the most actionable inner panel; desktop crops show the full multi-layer composition.

## Known Gaps

- **Hover state colors:** only press states (`button-primary-pill-pressed` → #2e2b8c) are documented. Marketing hover specs were not reliably extractable from the captured surfaces.
- **Dashboard product semantic palette:** error / success / warning states live inside the authenticated dashboard product (`dashboard.stripe.com`) and are not part of this marketing-system spec.
- **Sub-brand accents:** Stripe Atlas, Climate, Connect, Terminal, and Issuing each carry small accent palettes layered on top of the core indigo system — not captured here.
- **Loading skeletons:** the dashboard skeleton screens (gray-shimmer placeholders sized to match `body-tabular` rows) were not extracted in full.
- **Motion specs:** the gradient mesh has a subtle drift animation on the homepage; precise easing and duration tokens are not captured.
- **Code-block / syntax highlighting:** the inline code preview inside dashboard composites uses a custom Stripe palette (purple keywords, peach strings) that lives in its own token set, not the core marketing one.
