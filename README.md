# Katapult ACME DNS Solver

This is a simple Go library which exposes functionality to create ACME challenge DNS queries using Katapult's DNS service.

The API is fairly straight forward as shown below. The zone in question must be available to the API token which has been provided.

```go
// Create a new solver by providing the organization sub-domain and the
// API token.
logger := log.New(os.Stdout, "", log.LstdFlags)
solver := solver.NewSolver("api-token", logger)

// Set the record for the given zone with the given value.
err := solver.Set("example.com", "_acme-challenge.example.com", "abcdef")

// Cleanup the given records when appropriate
err = solver.Cleanup("example.com", "_acme-challenge.example.com")
```
