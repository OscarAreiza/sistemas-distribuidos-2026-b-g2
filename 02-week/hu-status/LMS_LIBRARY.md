# Product Brief — Loan Management System (LMS)

| Field | Value |
|---|---|
| Project key | LMS-LIBRARY-V1 |
| Version | 1.0 (first version) |
| Last updated | 2026-08-16 |
| Authors | Oscar Areiza, Hermes Pascuas, Luis Meneses |
| Status | In discovery — scope subject to change |
| Initial use case | University library (loans of books/copies) |

## 1. Tech Stack

| Layer | Technology |
|---|---|
| Frontend | Vue.js |
| Backend | Go |
| Database | PostgreSQL |
| Runtime environment | Docker (containerized services) |
| Deployment target | Cloud-hosted where possible (virtualized infrastructure) |

The system is designed to run as a set of Dockerized services (frontend, backend, database) so the same setup works locally and, where possible, on a cloud/virtualized environment.

## 2. Context

The university library is an academic service used by students, faculty, and staff to search, borrow, and return library materials. Today, part of this process relies on manual records, spreadsheets, or disconnected tools, which makes it difficult to know in real time which materials are available and what the status of each loan is.

This project is being built by eighth-semester students as an applied academic solution. The goal is a functional system that organizes core library management, improves material lookup, and streamlines loans and returns.

## 3. Scope

### In scope
- Centralized registry of books, users, and loans.
- Availability lookup and a basic operations history.
- Core loan/return workflow described in [Section 5](#5-expected-process).

### Out of scope (v1)
- A full enterprise-grade library management platform.
- Advanced features not explicitly listed in [Section 4](#4-needs-and-problems) unless promoted from [Section 6](#6-open-questions).

## 4. Needs and Problems

- Register books with basic information: title, author, ISBN, category, year, and available quantity.
- Register and manage the users authorized to use the library service.
- Search books by title, author, ISBN, or category.
- Quickly determine whether a book is available, on loan, or unavailable.
- Register loans of books to students, faculty, or other authorized users.
- Register returns and automatically update material availability.
- Query active loans and a basic loan history.
- Identify overdue loans to support follow-up by library staff.
- Remove the dependency on paper records, spreadsheets, or scattered information.
- Keep the operation simple and appropriate for a project built by eighth-semester students.

**Core problem:** information about books, users, and loans can be scattered or updated manually, making it hard to know in real time what material is available, who has it on loan, and which loans are pending return.

## 5. Expected Process

1. Register the book with its bibliographic data and available quantity.
2. Register the user who will use the library service.
3. Search for the requested book by title, author, ISBN, or category.
4. Check whether copies are available.
5. Register the loan, including user, book, loan date, and expected return date.
6. Update book availability once the loan is registered.
7. Register the return of the material.
8. Update book availability again.
9. Query active loans, loan history, and overdue loans.

## 6. Open Questions

| # | Question | Status | Notes |
|---|---|---|---|
| 1 | What user types can the system register: students, faculty, staff, others? | Open | |
| 2 | How many books/copies can a user request at the same time? | Open | |
| 3 | What is the maximum allowed loan period? | Open | |
| 4 | Can loans be renewed before the due date? | Open | |
| 5 | What happens when a user returns a book after the due date? | Open | |
| 6 | Do we need to track fines or penalties for late returns? | Open | |
| 7 | Should books and users be editable or deletable once registered? | Open | |
| 8 | Is a role-based authentication system required? | Open | |
| 9 | Will there be a single library administrator or multiple admin users? | Open | |
| 10 | Are reports needed (most borrowed books, overdue loans, users with active loans)? | Open | |
| 11 | Do reports need to be exportable to CSV or PDF? | Open | |
| 12 | Are notifications/reminders needed for due dates? | Open | |
| 13 | Can a book be reserved when all copies are on loan? | Open | |
| 14 | Does the physical condition of books need to be tracked? | Open | |
| 15 | What level of security and data backup is expected for v1? | Open | |

## 7. Business Glossary

| Term | Definition |
|---|---|
| Administrator | User responsible for managing books, users, and library operations. |
| Availability | Number of copies currently available for loan. |
| Book | Bibliographic material registered in the system and available for lookup or loan. |
| Category | Subject classification used to organize and search books. |
| Copy | A specific physical instance of a book that can be loaned to a user. |
| ISBN | International identifier used to identify a specific edition of a book. |
| Loan | Record of the temporary delivery of a book to a user. |
| Loan date | Date on which the book is delivered to the user. |
| Loan history | Record of past loans and returns. |
| Active loan | A loan that has not yet been marked as returned. |
| Overdue loan | A loan whose expected return date has passed and is still active. |
| Return | The process by which a user gives a book back to the library. |
| Return date | Date on which the user must return the book, or the date it was actually returned. |
| Reservation | A request made by a user to obtain a book when no copies are currently available. |
| User | A person authorized to use the library services. |
