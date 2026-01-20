# Foresee

Foresee is a small prediction market web application built as a learning project after finishing the book *Let’s Go* by Alex Edwards.

The goal of the project was to consolidate Go web fundamentals by building a complete application end to end, focusing on applying everything that I have learned by reading the book rather than feature richness or financial realism.

---

## Project Idea

Users can create markets by asking a question with a Yes and No outcomes. Other users place bets using virtual coins on exactly one outcome. When a market expires, an authorized resolver selects the winning outcome and all coins from losing outcomes are redistributed proportionally to the winners.

Coins have no real-world value and exist only to explore backend logic and accounting behavior.

---

## Purpose

This project was built to:

* Apply the concepts learned in *Let’s Go* in a non-trivial project
* Practice building a Go web app with routing, templates, sessions, auth, and a database
* Work through real state transitions (open → expired → resolved)
* Implement safe, transactional updates for money-like values

It is a learning and portfolio project, not a production product.

---

## Key Characteristics

* Pool-based markets
* Explicit market resolution by a user
* Proportional payouts to winning bets
* Server-rendered HTML using Go templates
* Minimal dependencies, standard library first
---

## Challenges

* Designing payout logic that is simple but correct
* Preventing double resolution and double payouts
* Structuring code to keep responsibilities clear
* Handling edge cases like no bets on the winning outcome
---

## How to Run

The project is fully containerized and works out of the box.

```bash
docker compose up -d --build
```

Once the containers are up, the application will be available locally with no additional setup required.

---

## Pending Improvements

* Cleanup of the application to rely consistently on services
* Reduce direct database access outside repositories
* Add unit tests for business logic
* Add integration tests for critical flows
* Improve error handling and messages
* Add automated market expiration

---

## Summary

Foresee is a deliberately scoped Go project built to turn theory into practice after reading *Let’s Go*.