---

title: Use Strategy Pattern
status: adopted
keywords: [design-pattern, strategy, interface, decoupling]
---

# Use Strategy Pattern

When implementing behavior that varies based on configuration, platform, or context, favor the **Strategy Pattern** over conditional logic (if/else chains).

## Principles

1.  **Define an Interface**: Create an interface that defines the contract for the behavior.
2.  **Implement Strategies**: Create separate structs/classes that implement this interface for each specific behavior.
3.  **Inject Dependencies**: Pass the selected strategy (as the interface) to the consumer/context, typically via constructor injection.

## Benefits

-   **Decoupling**: The consumer doesn't need to know the implementation details.
-   **Testability**: Strategies can be mocked for unit testing.
-   **Extensibility**: New behaviors can be added without modifying the consumer code (Open/Closed Principle).
-   **Maintainability**: Eliminates complex conditional logic.

## Anti-Pattern

```go
func (s *Service) DoSomething(ctx context.Context) {
    if s.Config.Type == "local" {
        // Local logic...
    } else if s.Config.Type == "remote" {
        // Remote logic...
    }
}
```

## Recommended Pattern

```go
type Provider interface {
    Do(ctx context.Context) error
}

type LocalProvider struct {}
func (p *LocalProvider) Do(ctx context.Context) error { /* ... */ }

type RemoteProvider struct {}
func (p *RemoteProvider) Do(ctx context.Context) error { /* ... */ }

type Service struct {
    Provider Provider
}

func NewService(p Provider) *Service {
    return &Service{Provider: p}
}
```
