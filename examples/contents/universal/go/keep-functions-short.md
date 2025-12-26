---
id: keep-functions-short
title: Keep functions short by top-down decomposition
description: >
  Write functions by detailing steps as comments first, then converting them to function calls.
keywords:
  - readability
  - maintainability
  - function-design
  - go
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
```go
func ConfirmBilling(cart Cart) (BillingResult, error) {
    // Complex validation logic inline...
    if len(cart.Items) == 0 {
        return BillingResult{}, errors.New("empty cart")
    }
    // ...
    // Complex calculation inline...
    total := 0
    for _, item := range cart.Items {
        total += item.Price
    }
    // ...
    return BillingResult{Total: total}, nil
}
```

### Good
```go
func ConfirmBilling(cart Cart) (BillingResult, error) {
    if err := validateCart(cart); err != nil {
        return BillingResult{}, err
    }

    usage, err := collectUsage(cart)
    if err != nil {
        return BillingResult{}, err
    }

    amount := calculateAmount(cart, usage)
    
    result, err := saveBilling(cart, amount)
    if err != nil {
        return BillingResult{}, err
    }

    if err := notifyBilling(result); err != nil {
        // Log warning but don't fail flow
        log.Warn("notification failed", "err", err)
    }

    return result, nil
}
```
