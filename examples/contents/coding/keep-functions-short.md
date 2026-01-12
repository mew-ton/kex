---

title: Keep functions short by top-down decomposition
description: >
  When implementing logic, keep functions short by applying top-down decomposition using comments-first development.
keywords:
  - readability
  - maintainability
  - function
  - length
  - complexity
  - size
  - decomposition
---

## Summary
Do not write long logic immediately. First, describe the process in comments/verbs, then convert those comments into function calls.

## Rationale
- Breaking down logic into named functions explains "what" is happening without showing "how".
- It makes code readable like a sentence.
- It clarifies responsibilities early.

## Guidance
1.  Write the steps of your function as comments (verbs).
2.  Convert comments into function calls (even if they don't exist yet).
3.  Implement the sub-functions recursively using the same method.

## Examples

### Bad
```text
function ConfirmBilling(cart):
    // ... complex logic ...
    // ... 50 lines later ...
    return result
```

### Good
```text
function ConfirmBilling(cart):
    validateCart(cart)
    usage = collectUsage(cart)
    amount = calculateAmount(cart, usage)
    saveBilling(cart, amount)
    notifyBilling(result)
```
