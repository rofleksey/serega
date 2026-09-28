# Card presentation

Own card presentation and status metadata.

- Canonical values are todo, doing, done; labels and colors are presentation.
- Keep ID, version, creator/updater, and timestamps distinct.
- Present user text as text; never inject card content as HTML.
- Receive action callbacks where possible so presentation does not own page
  navigation or duplicate mutation orchestration.
- Preserve readable long text and accessible names for controls.
- Do not interpret creator identity as a permission boundary.
