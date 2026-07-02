# Incident Handoff — Decisions
## First words

Incident Handoff is a handoff documentation.

I understand that for such 1-person-project, of course we just need minimal version with simple CRUD, that's it. no need to worry much about data-race micro-tuning, cache, query-speed, and so on.

But, this is also a learning project. It's interesting to do it as a decent project, as if I were building a real one — where I have to make decisions and be aware of the tradeoffs, and what if in scale. So I document the decisions here, also to support my learning.

## Why using versioning in IncidentHandler and IncidentStore, funcs AddEntry and UpdateIncident? What's the trade-off in current implementation?

Well, let's have a look into the workflow

1. IncidentHandler calls IncidentStore to get something
2. IncidentHandler check some thing
3. IncidentHandler calls IncidentStore to change something
4. IncidentStore changes it and returns something back to IncidentHandler

in this case, 1 and 2 belong to Time-Of-Check (TOC). 3 and 4 belong to Time-Of-Use (TOU).

Originally, I think, to prevent race, I'll use mutex. It's very traditional.
But, for that, I need to lock mutex before 1. and unlock mutex after 4. That's a very long jouney. Other requests need to wait for you. And at scale, it's not safe. Race can happen if multiple instances exist. So I researched.

https://stackoverflow.com/questions/129329/optimistic-vs-pessimistic-locking

So, in this case, versioning will solve this case with optimistic locking. Incident-handoff case is a low contention case. In usual case, everything works, no interuption. If races happen, they will be detected and rejected by returning Conflict and the user owns the next step of what to do. This is not expensive because the race barely happens. And this one also works easily even with scaling (multiple instances).

On other side, Pessimistic locking is suitable with high-volumn systems or cases (flash-sale i guess), using not mutex (process-locking-level), but an external reliable service. 

## Why VerifyPassword(ARealBrcyptHasedPassword, "dummy-password")?
Because, An attacker can measure the response time:
- fast = no such user
- slow = user exists
=> This leaks which usernames are registered.

Note: In current implementation, an acocunt is registered directly -> attacker can try to register an account to detect this. But if registration requires an E-mail and send confirmation to that E-Mail -> no problem.