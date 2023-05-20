# Dining Philosophers

This is another solution to the Dining Philosopher's problem written in [Go](https://go.dev) this time. There seem to be some issues where a
philosopher will continue to eat and think back to back. Not exactly sure why that's happening, but it's not happening enough to be a serious concern.

## Visualizer

I'd like to build a TUI based visualizer of this solution's processing. Should be quite doable. Just would probably use an atomic value for tracking the
state of a philosopher (eating or thinking) and another for tracking their apetite.

## Other Solutions

- [Rust](https://github.com/JingusJohn/dining_philosophe_rs)
- [C](https://github.com/JingusJohn/dining-philosophers/tree/main)

So far, this has been my simplest solution. Golang's concurrency approach is so much simpler than what I had to deal with for C and Rust. I've wanted to
create visualizers for those as well, but this one should actually be quite doable. Rust's borrow checker makes visualizing the changes quite difficult
for me.

