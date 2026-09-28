# Composed UI sections

Widgets compose entities/features/shared code into substantial reusable screen
sections. They must not import pages or app.

- Keep owned Query state and user-visible async states together.
- Call shared/api and compose features instead of duplicating transport/forms.
- Keep components cohesive; do not create a generic dashboard framework.
- Tests should exercise observable user behavior, errors, and updates.
