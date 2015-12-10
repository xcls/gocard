# TODO

- [ ] Case-insensitive check of email addresses for authentication
  - Change to `citext` column
- [ ] Uniqueness validation of email address

- [ ] If the quality response was lower than 3 then start repetitions for the
  item from the beginning without changing the E-Factor (i.e. use intervals
  I(1), I(2) etc. as if the item was memorized anew).
- [ ] After each repetition session of a given day repeat again all items that
  scored below four in the quality assessment. Continue the repetitions until
  all of these items score at least four.
